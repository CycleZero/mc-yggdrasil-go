package utils

import (
	"crypto/md5"
	"github.com/google/uuid"
)

// NameUUIDFromBytes replicates the Java UUID.nameUUIDFromBytes functionality
// It generates a version 3 (MD5-based) UUID from the provided byte array
func NameUUIDFromBytes(data []byte) (uuid.UUID, error) {
	// Use MD5 hash as per Java's UUID.nameUUIDFromBytes implementation

	// Create and return UUID from the hashed bytes
	result, err := uuid.FromBytes(hash[:])
	if err != nil {
		return uuid.Nil, err
	}
	return result, nil
}
