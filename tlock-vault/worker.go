package tlockvault

import (
	"time"

	"github.com/eklairs/tlock/tlock-internal/utils"
	"github.com/kelindar/binary"
)

// Writing to file implementation
func (vault *Vault) startFileWriterWorker(recv chan []Folder) {
	for {
		if data, ok := <-recv; ok {
			// Serialize
			serialized, _ := binary.Marshal(data)

			// Encrypt
			encrypted, _ := Encrypt(vault.password, serialized)

			// Create parent dir
			if file, err := utils.EnsureExists(vault.path); err == nil {
				file.Write(encrypted)
			}
		}

		// Sleep for 1 second
		time.Sleep(time.Second * 1)
	}
}
