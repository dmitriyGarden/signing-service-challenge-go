package inmemory

import (
	"context"
	"testing"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/devices"
	"github.com/google/uuid"
)

func TestSignatureStoreSequence(t *testing.T) {
	store := NewSignatureStore()
	now := time.Now().UTC()
	uid := uuid.New()

	first, err := store.Append(context.Background(), uid, devices.SignatureRecord{Signature: "sig1", SignedData: "0_payload", CreatedAt: now})
	if err != nil {
		t.Fatalf("append failed: %v", err)
	}
	if first.Counter != 1 {
		t.Fatalf("expected counter 1, got %d", first.Counter)
	}

	second, err := store.Append(context.Background(), uid, devices.SignatureRecord{Signature: "sig2", SignedData: "1_payload", CreatedAt: now.Add(time.Second)})
	if err != nil {
		t.Fatalf("append failed: %v", err)
	}
	if second.Counter != 2 {
		t.Fatalf("expected counter 2, got %d", second.Counter)
	}

	list, err := store.List(context.Background(), uid)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(list) != 2 || list[1].Signature != "sig2" {
		t.Fatalf("unexpected list result: %#v", list)
	}

	record, err := store.Get(context.Background(), uid, 2)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if record.Signature != "sig2" {
		t.Fatalf("unexpected signature: %s", record.Signature)
	}

	if _, err := store.Get(context.Background(), uid, 3); err == nil {
		t.Fatalf("expected not found for missing counter")
	}

	last, ok, err := store.Last(context.Background(), uid)
	if err != nil {
		t.Fatalf("last failed: %v", err)
	}
	if !ok || last.Signature != "sig2" {
		t.Fatalf("unexpected last record: %#v", last)
	}
	unknown := uuid.New()
	_, ok, err = store.Last(context.Background(), unknown)
	if err != nil || ok {
		t.Fatalf("expected no record for unknown device")
	}
}

func TestSignatureStoreClone(t *testing.T) {
	store := NewSignatureStore()
	now := time.Now()
	uid := uuid.New()
	record, err := store.Append(context.Background(), uid, devices.SignatureRecord{Signature: "sig", SignedData: "data", CreatedAt: now})
	if err != nil {
		t.Fatalf("append failed: %v", err)
	}

	record.Signature = "mutated"

	stored, err := store.Get(context.Background(), uid, 1)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if stored.Signature != "sig" {
		t.Fatalf("expected original signature, got %s", stored.Signature)
	}
}

func TestSignatureStoreGetInvalidCounter(t *testing.T) {
	store := NewSignatureStore()
	uid := uuid.New()
	if _, err := store.Get(context.Background(), uid, 0); err == nil {
		t.Fatalf("expected error for counter 0")
	}
	_, err := store.Get(context.Background(), uid, 1)
	if _, ok := err.(domain.NotFoundError); !ok {
		t.Fatalf("expected NotFoundError, got %v", err)
	}
}
