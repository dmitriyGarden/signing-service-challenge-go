package inmemory

import (
	"context"
	"testing"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
)

func TestDeviceRepositoryCRUD(t *testing.T) {
	repo := NewDeviceRepository()
	uid := uuid.New()
	device := domain.Device{ID: uid, Algorithm: domain.AlgorithmRSA, Label: "POS", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	if err := repo.Create(context.Background(), device); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if err := repo.Create(context.Background(), device); err != domain.ErrDeviceExists {
		t.Fatalf("expected duplicate error, got %v", err)
	}

	stored, err := repo.Get(context.Background(), uid)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if stored.ID != device.ID || stored.Label != device.Label {
		t.Fatalf("unexpected stored device: %#v", stored)
	}

	stored.Label = "POS-2"
	if err := repo.Update(context.Background(), stored); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	list, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(list) != 1 || list[0].Label != "POS-2" {
		t.Fatalf("unexpected list result: %#v", list)
	}

	if err := repo.Delete(context.Background(), uid); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	if err := repo.Delete(context.Background(), uid); err != nil {
		t.Fatalf("unexpected error deleting missing device")
	}
}
