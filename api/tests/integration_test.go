package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"sync"
	"testing"
	"time"

	v0 "github.com/fiskaly/coding-challenges/signing-service-challenge/api/v0"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	appdevices "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/devices"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/persistence/inmemory"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/crypto"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func newTestHandler() http.Handler {
	repo := inmemory.NewDeviceRepository()
	keyStore := inmemory.NewKeyStore()
	keyGenerator := crypto.NewDefaultKeyGenerator()
	signerFactory := crypto.NewSignerFactory()
	signatureStore := inmemory.NewSignatureStore()

	core := appdevices.NewService(repo, keyStore, keyGenerator, signerFactory, signatureStore)
	core.WithClock(func() time.Time { return time.Unix(0, 0).UTC() })
	handler := v0.NewHandler(core)

	router := chi.NewRouter()
	router.Route("/api/v0", handler.Register)

	return router
}

type testClient struct {
	handler http.Handler
}

type httpResult struct {
	status int
	body   []byte
}

func (c *testClient) request(t *testing.T, method, path string, payload interface{}) httpResult {
	t.Helper()
	var body bytes.Buffer
	if payload != nil {
		if err := json.NewEncoder(&body).Encode(payload); err != nil {
			t.Fatalf("encode payload: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &body)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	c.handler.ServeHTTP(rec, req)
	return httpResult{status: rec.Code, body: rec.Body.Bytes()}
}

func decodeData[T any](t *testing.T, result httpResult, target *T) {
	t.Helper()
	if result.status < 200 || result.status >= 300 {
		t.Fatalf("unexpected status %d", result.status)
	}
	var wrapper struct {
		Data T `json:"data"`
	}
	if err := json.Unmarshal(result.body, &wrapper); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	*target = wrapper.Data
}

func TestDeviceLifecycleIntegration(t *testing.T) {
	client := testClient{handler: newTestHandler()}
	basePath := "/api/v0"

	deviceID := uuid.New()

	createResp := client.request(t, http.MethodPost, basePath+"/devices/", map[string]any{
		"id":        deviceID.String(),
		"algorithm": string(domain.AlgorithmRSA),
		"label":     "Register",
	})
	var created struct {
		ID        string `json:"id"`
		Algorithm string `json:"algorithm"`
		Label     string `json:"label"`
		Counter   uint64 `json:"counter"`
	}
	decodeData(t, createResp, &created)
	if created.ID != deviceID.String() {
		t.Fatalf("expected id %s, got %s", deviceID, created.ID)
	}

	listResp := client.request(t, http.MethodGet, basePath+"/devices/", nil)
	var devicesList []struct {
		ID        string `json:"id"`
		Algorithm string `json:"algorithm"`
		Label     string `json:"label"`
		Counter   uint64 `json:"counter"`
	}
	decodeData(t, listResp, &devicesList)
	if len(devicesList) != 1 || devicesList[0].Counter != 0 {
		t.Fatalf("unexpected device list: %#v", devicesList)
	}

	signResp := client.request(t, http.MethodPost, basePath+"/devices/"+deviceID.String()+"/sign", map[string]any{"data": "sale"})
	var signResult struct {
		Signature  string `json:"signature"`
		SignedData string `json:"signed_data"`
	}
	decodeData(t, signResp, &signResult)
	if signResult.Signature == "" {
		t.Fatal("expected signature to be set")
	}

	signaturesResp := client.request(t, http.MethodGet, basePath+"/devices/"+deviceID.String()+"/signatures", nil)
	var signatures []struct {
		Counter    uint64    `json:"counter"`
		Signature  string    `json:"signature"`
		SignedData string    `json:"signed_data"`
		CreatedAt  time.Time `json:"created_at"`
	}
	decodeData(t, signaturesResp, &signatures)
	if len(signatures) != 1 || signatures[0].Counter != 1 {
		t.Fatalf("unexpected signatures: %#v", signatures)
	}

	singleResp := client.request(t, http.MethodGet, basePath+"/devices/"+deviceID.String()+"/signatures/1", nil)
	var single struct {
		Counter    uint64    `json:"counter"`
		Signature  string    `json:"signature"`
		SignedData string    `json:"signed_data"`
		CreatedAt  time.Time `json:"created_at"`
	}
	decodeData(t, singleResp, &single)
	if single.Counter != 1 {
		t.Fatalf("expected counter 1, got %d", single.Counter)
	}
}

func TestConcurrentSigningIntegration(t *testing.T) {
	client := testClient{handler: newTestHandler()}
	basePath := "/api/v0"

	deviceID := uuid.New()
	createResp := client.request(t, http.MethodPost, basePath+"/devices/", map[string]any{
		"id":        deviceID.String(),
		"algorithm": string(domain.AlgorithmECDSA),
		"label":     "Concurrent",
	})
	var created struct {
		ID string `json:"id"`
	}
	decodeData(t, createResp, &created)

	const workers = 10
	const perWorker = 5
	var wg sync.WaitGroup
	wg.Add(workers)
	errCh := make(chan error, workers*perWorker)

	for i := 0; i < workers; i++ {
		go func(worker int) {
			defer wg.Done()
			for j := 0; j < perWorker; j++ {
				payload := map[string]any{"data": fmt.Sprintf("payload-%d-%d", worker, j)}
				resp := client.request(t, http.MethodPost, basePath+"/devices/"+deviceID.String()+"/sign", payload)
				if resp.status != http.StatusOK {
					errCh <- fmt.Errorf("unexpected status %d", resp.status)
					return
				}
			}
		}(i)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
		close(errCh)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("concurrent signing did not finish in time")
	}

	for err := range errCh {
		if err != nil {
			t.Fatalf("concurrent sign failed: %v", err)
		}
	}

	signaturesResp := client.request(t, http.MethodGet, basePath+"/devices/"+deviceID.String()+"/signatures", nil)
	var signatures []struct {
		Counter uint64 `json:"counter"`
	}
	decodeData(t, signaturesResp, &signatures)

	expected := workers * perWorker
	if len(signatures) != expected {
		t.Fatalf("expected %d signatures, got %d", expected, len(signatures))
	}

	counters := make([]uint64, 0, len(signatures))
	for _, record := range signatures {
		counters = append(counters, record.Counter)
	}
	sort.Slice(counters, func(i, j int) bool { return counters[i] < counters[j] })
	for idx, value := range counters {
		want := uint64(idx + 1)
		if value != want {
			t.Fatalf("expected counter %d at position %d, got %d", want, idx, value)
		}
	}
}
