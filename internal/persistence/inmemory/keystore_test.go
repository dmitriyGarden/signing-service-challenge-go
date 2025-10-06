package inmemory

import (
	"context"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
)

func TestKeyStoreLifecycle(t *testing.T) {
	store := NewKeyStore()
	uid := uuid.New()
	material := domain.KeyMaterial{Public: []byte("pub"), Private: []byte("priv")}

	if err := store.Store(context.Background(), uid, material); err != nil {
		t.Fatalf("store failed: %v", err)
	}

	loaded, err := store.Load(context.Background(), uid)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	loaded.Public[0] = 'P'
	if string(material.Public) != "pub" {
		t.Fatalf("store mutated original material")
	}

	again, err := store.Load(context.Background(), uid)
	if err != nil {
		t.Fatalf("second load failed: %v", err)
	}
	if string(again.Public) != "pub" {
		t.Fatalf("stored material was mutated: %s", string(again.Public))
	}

	if err := store.Delete(context.Background(), uid); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	if _, err := store.Load(context.Background(), uid); err != domain.ErrKeyMaterialMissing {
		t.Fatalf("expected missing key error, got %v", err)
	}
}
