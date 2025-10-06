package devices

import "time"

// SignatureRecord represents a stored signature for a device.
type SignatureRecord struct {
	Counter    uint64
	Signature  string
	SignedData string
	CreatedAt  time.Time
}

// Clone returns a copy to avoid leaking pointers.
func (r SignatureRecord) Clone() SignatureRecord {
	return r
}
