package tlockvault

import (
	"log"
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

func remove(slice []string, s int) []string {
    return append(slice[:s], slice[s+1:]...)
}

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

func Initialize(password string) TLockVault {
    id := uuid.New()
    dir := path.Join(xdg.DataHome, "tlock", "root", id.String())

    if err := os.MkdirAll(dir, os.ModePerm); err != nil {
        log.Fatalf("Failed to create user's root dir: %v", err)
    }

    vault := TLockVault {
        Data: TLockVaultData{},
        Vault_path: path.Join(dir, "vault.dat"),
    }

    return vault
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

func (vault TLockVault) write() {
    f, _ := os.Create(vault.Vault_path)

    b, _ := yaml.Marshal(vault.Data)

    if _, err := f.Write(b); err != nil {

    }
}

func (vault *TLockVault) AddFolder(name string) {
    vault.Data.Folders = append(vault.Data.Folders, struct{Name string; Uris []string}{ Name: name, Uris: []string {} })

    vault.write()
}

func (vault *TLockVault) UpdateURI(folder, uri int, newURI string) {
    vault.Data.Folders[folder].Uris[uri] = newURI

    vault.write()
}

func (vault *TLockVault) DeleteURI(folder, uri int) {
    vault.Data.Folders[folder].Uris = remove(vault.Data.Folders[folder].Uris, uri)

    vault.write()
}

