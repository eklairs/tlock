package tlockcore

import (
	"log"
	"os"
	"path"

	"github.com/adrg/xdg"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"gopkg.in/yaml.v2"
)

var USERS_PATH = path.Join(xdg.DataHome, "tlock", "users.yaml")

// Represents a user
type UserSpec struct {
    Username string `yaml:"username"`
    Vault string `yaml:"vault"`
}

// Users API for tlock
type TLockUsers struct {
    Users []UserSpec `yaml:"users"`
}

// Loads the list of users
func LoadTLockUsers() TLockUsers {
    // Read users path
    raw, err := os.ReadFile(USERS_PATH);

    // Check for errors
    if err != nil {
        log.Printf("[users] Users map does not exist, returning empty users map")

        return TLockUsers{}
    }

    // Empty users where yaml will populate
    users := TLockUsers{}

    // Parse
    if err = yaml.Unmarshal(raw, &users); err != nil {
        log.Fatalf("[users] Failed to parse users map, syntax error possibly: %v", err)
    }

    return users
}

// Writes the current users value to the file
func (users TLockUsers) write() {
    data, _ := yaml.Marshal(users)

    file, err := os.Create(path.Join(xdg.DataHome, "tlock", "users.yaml"))

    if err != nil {
        log.Fatalf("Failed to create users file")
    }

    file.Write(data)
}

// Adds a new user
func (users *TLockUsers) AddNewUser(username, password string) tlockvault.TLockVault {
    vault := tlockvault.Initialize(password)
    users.Users = append(users.Users, UserSpec{ Username: username, Vault: vault.Vault_path })

    users.write()

    return vault
}
