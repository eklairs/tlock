package tlockvault

import (
	"log"
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

type TLockVaultData struct {
    Folders []struct {
        Name string
        Uris []string
    }
}

type TLockVault struct {
    Data TLockVaultData
    Vault_path string
}

func Initialize(password string) string {
    id := uuid.New()
    dir := path.Join(xdg.DataHome, "tlock", "root", id.String())

    if err := os.MkdirAll(dir, os.ModePerm); err != nil {
        log.Fatalf("Failed to create user's root dir: %v", err)
    }

    return path.Join(dir, "vault.dat")
}

func Load(path, password string) (*TLockVault, error) {
    raw, err := os.ReadFile(path)

    if err != nil {
        return nil, err
    }

    data, _err := Decrypt(raw)

    if _err != nil {
        return nil, _err
    }

    return &TLockVault {
        Data: *data,
        Vault_path: path,
    }, nil
}

func Decrypt(data []byte) (*TLockVaultData, error) {
    out := TLockVaultData{}

    if err := yaml.Unmarshal(data, &out); err != nil {
        return nil, err
    }

    return &out, nil
}
