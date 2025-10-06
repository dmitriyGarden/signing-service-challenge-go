package devices_test

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/devices"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func fixedTime() time.Time {
	return time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
}

func TestService_CreateDevice_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	keyStore := mocks.NewMockKeyStore(ctrl)
	keyGen := mocks.NewMockKeyGenerator(ctrl)
	signerFactory := mocks.NewMockSignerFactory(ctrl)
	sigStore := mocks.NewMockSignatureStore(ctrl)

	service := devices.NewService(repo, keyStore, keyGen, signerFactory, sigStore)
	service.WithClock(fixedTime)

	id := uuid.New()
	input := devices.CreateDeviceInput{
		ID:        id,
		Algorithm: domain.AlgorithmRSA,
		Label:     "  demo terminal  ",
	}

	material := domain.KeyMaterial{Public: []byte("pub"), Private: []byte("priv")}

	keyGen.EXPECT().Generate(domain.AlgorithmRSA).Return(material, nil)

	repo.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(domain.Device{})).DoAndReturn(func(_ context.Context, device domain.Device) error {
		if device.ID != id {
			t.Fatalf("unexpected device id %s", device.ID)
		}
		if device.Label != "demo terminal" {
			t.Fatalf("expected trimmed label, got %q", device.Label)
		}
		if !device.CreatedAt.Equal(fixedTime()) {
			t.Fatalf("expected created at %v, got %v", fixedTime(), device.CreatedAt)
		}
		if !device.UpdatedAt.Equal(fixedTime()) {
			t.Fatalf("expected updated at %v, got %v", fixedTime(), device.UpdatedAt)
		}
		return nil
	})

	keyStore.EXPECT().Store(gomock.Any(), id, material).Return(nil)

	result, err := service.CreateDevice(context.Background(), input)
	if err != nil {
		t.Fatalf("CreateDevice returned error: %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if result.Device.ID != id {
		t.Fatalf("unexpected device id %s", result.Device.ID)
	}
	if result.Device.Label != "demo terminal" {
		t.Fatalf("expected trimmed label, got %q", result.Device.Label)
	}
	if !result.Device.CreatedAt.Equal(fixedTime()) {
		t.Fatalf("unexpected created timestamp %v", result.Device.CreatedAt)
	}
}

func TestService_CreateDevice_KeyStoreFailureRollsBack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	keyStore := mocks.NewMockKeyStore(ctrl)
	keyGen := mocks.NewMockKeyGenerator(ctrl)
	signerFactory := mocks.NewMockSignerFactory(ctrl)
	sigStore := mocks.NewMockSignatureStore(ctrl)

	service := devices.NewService(repo, keyStore, keyGen, signerFactory, sigStore)
	service.WithClock(fixedTime)

	id := uuid.New()
	input := devices.CreateDeviceInput{ID: id, Algorithm: domain.AlgorithmRSA, Label: "label"}
	material := domain.KeyMaterial{Public: []byte("pub"), Private: []byte("priv")}

	keyGen.EXPECT().Generate(domain.AlgorithmRSA).Return(material, nil)
	repo.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(domain.Device{})).Return(nil)
	expectedErr := errors.New("boom")
	keyStore.EXPECT().Store(gomock.Any(), id, material).Return(expectedErr)
	repo.EXPECT().Delete(gomock.Any(), id).Return(nil)

	result, err := service.CreateDevice(context.Background(), input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected wrapped error, got %v", err)
	}
	if result != nil {
		t.Fatal("expected nil result when store fails")
	}
}

func TestService_SignTransaction_FirstSignature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	keyStore := mocks.NewMockKeyStore(ctrl)
	keyGen := mocks.NewMockKeyGenerator(ctrl)
	signerFactory := mocks.NewMockSignerFactory(ctrl)
	sigStore := mocks.NewMockSignatureStore(ctrl)
	signer := mocks.NewMockSigner(ctrl)

	service := devices.NewService(repo, keyStore, keyGen, signerFactory, sigStore)
	service.WithClock(fixedTime)

	id := uuid.New()
	device := domain.Device{ID: id, Algorithm: domain.AlgorithmRSA, Label: "demo"}
	material := domain.KeyMaterial{Public: []byte("pub"), Private: []byte("priv")}

	repo.EXPECT().Get(gomock.Any(), id).Return(device, nil)
	keyStore.EXPECT().Load(gomock.Any(), id).Return(material, nil)
	signerFactory.EXPECT().SignerFor(device, material).Return(signer, nil)
	sigStore.EXPECT().Last(gomock.Any(), id).Return(devices.SignatureRecord{}, false, nil)

	payload := domain.BuildSecuredPayload(1, "data", id[:])
	signatureBytes := []byte("signed")
	signer.EXPECT().Sign([]byte(payload)).Return(signatureBytes, nil)

	sigStore.EXPECT().Append(gomock.Any(), id, gomock.AssignableToTypeOf(devices.SignatureRecord{})).DoAndReturn(
		func(_ context.Context, _ uuid.UUID, record devices.SignatureRecord) (devices.SignatureRecord, error) {
			if record.SignedData != payload {
				t.Fatalf("unexpected payload %q", record.SignedData)
			}
			expectedSig := base64.StdEncoding.EncodeToString(signatureBytes)
			if record.Signature != expectedSig {
				t.Fatalf("expected signature %q, got %q", expectedSig, record.Signature)
			}
			if !record.CreatedAt.Equal(fixedTime()) {
				t.Fatalf("expected CreatedAt %v, got %v", fixedTime(), record.CreatedAt)
			}
			record.Counter = 1
			return record, nil
		},
	)

	input := devices.SignTransactionInput{DeviceID: id, Data: "data"}
	result, err := service.SignTransaction(context.Background(), input)
	if err != nil {
		t.Fatalf("SignTransaction returned error: %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.CounterValue != 1 {
		t.Fatalf("expected counter 1, got %d", result.CounterValue)
	}
	expectedSig := base64.StdEncoding.EncodeToString(signatureBytes)
	if result.Signature != expectedSig {
		t.Fatalf("unexpected signature %q", result.Signature)
	}
	if result.SignedData != payload {
		t.Fatalf("unexpected signed data %q", result.SignedData)
	}
}

func TestService_SignTransaction_UsesPreviousSignature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	keyStore := mocks.NewMockKeyStore(ctrl)
	keyGen := mocks.NewMockKeyGenerator(ctrl)
	signerFactory := mocks.NewMockSignerFactory(ctrl)
	sigStore := mocks.NewMockSignatureStore(ctrl)
	signer := mocks.NewMockSigner(ctrl)

	service := devices.NewService(repo, keyStore, keyGen, signerFactory, sigStore)
	service.WithClock(fixedTime)

	id := uuid.New()
	device := domain.Device{ID: id, Algorithm: domain.AlgorithmRSA, Label: "demo"}
	material := domain.KeyMaterial{Public: []byte("pub"), Private: []byte("priv")}

	repo.EXPECT().Get(gomock.Any(), id).Return(device, nil)
	keyStore.EXPECT().Load(gomock.Any(), id).Return(material, nil)
	signerFactory.EXPECT().SignerFor(device, material).Return(signer, nil)

	prevBytes := []byte("previous")
	prevSignature := base64.StdEncoding.EncodeToString(prevBytes)
	sigStore.EXPECT().Last(gomock.Any(), id).Return(devices.SignatureRecord{Signature: prevSignature, Counter: 1}, true, nil)

	payload := domain.BuildSecuredPayload(2, "payload", prevBytes)
	signatureBytes := []byte("new-sig")
	signer.EXPECT().Sign([]byte(payload)).Return(signatureBytes, nil)

	sigStore.EXPECT().Append(gomock.Any(), id, gomock.AssignableToTypeOf(devices.SignatureRecord{})).DoAndReturn(
		func(_ context.Context, _ uuid.UUID, record devices.SignatureRecord) (devices.SignatureRecord, error) {
			if record.SignedData != payload {
				t.Fatalf("unexpected payload %q", record.SignedData)
			}
			record.Counter = 2
			return record, nil
		},
	)

	input := devices.SignTransactionInput{DeviceID: id, Data: "payload"}
	result, err := service.SignTransaction(context.Background(), input)
	if err != nil {
		t.Fatalf("SignTransaction returned error: %v", err)
	}
	if result.CounterValue != 2 {
		t.Fatalf("expected counter 2, got %d", result.CounterValue)
	}
	if result.SignedData != payload {
		t.Fatalf("unexpected signed data %q", result.SignedData)
	}
}

func TestService_SignTransaction_ValidatesData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	keyStore := mocks.NewMockKeyStore(ctrl)
	keyGen := mocks.NewMockKeyGenerator(ctrl)
	signerFactory := mocks.NewMockSignerFactory(ctrl)
	sigStore := mocks.NewMockSignatureStore(ctrl)

	service := devices.NewService(repo, keyStore, keyGen, signerFactory, sigStore)

	_, err := service.SignTransaction(context.Background(), devices.SignTransactionInput{DeviceID: uuid.New(), Data: "  \t  "})
	if err == nil {
		t.Fatal("expected validation error")
	}
	var vErr domain.ValidationError
	if !errors.As(err, &vErr) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestService_GetCounters_ForwardsToStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	keyStore := mocks.NewMockKeyStore(ctrl)
	keyGen := mocks.NewMockKeyGenerator(ctrl)
	signerFactory := mocks.NewMockSignerFactory(ctrl)
	sigStore := mocks.NewMockSignatureStore(ctrl)

	service := devices.NewService(repo, keyStore, keyGen, signerFactory, sigStore)

	ids := []uuid.UUID{uuid.New(), uuid.New()}
	expected := map[uuid.UUID]uint64{ids[0]: 3, ids[1]: 4}

	sigStore.EXPECT().GetCounters(gomock.Any(), ids).Return(expected, nil)

	counters, err := service.GetCounters(context.Background(), ids)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(counters) != len(expected) {
		t.Fatalf("expected %d counters, got %d", len(expected), len(counters))
	}
	for id, value := range expected {
		if counters[id] != value {
			t.Fatalf("counter mismatch for %s: want %d got %d", id, value, counters[id])
		}
	}
}

func TestService_DeleteDevice_RemovesAllState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	keyStore := mocks.NewMockKeyStore(ctrl)
	keyGen := mocks.NewMockKeyGenerator(ctrl)
	signerFactory := mocks.NewMockSignerFactory(ctrl)
	sigStore := mocks.NewMockSignatureStore(ctrl)

	service := devices.NewService(repo, keyStore, keyGen, signerFactory, sigStore)

	id := uuid.New()

	gomock.InOrder(
		repo.EXPECT().Delete(gomock.Any(), id).Return(nil),
		keyStore.EXPECT().Delete(gomock.Any(), id).Return(nil),
		sigStore.EXPECT().Delete(gomock.Any(), id).Return(nil),
	)

	if err := service.DeleteDevice(context.Background(), id); err != nil {
		t.Fatalf("DeleteDevice returned error: %v", err)
	}
}
