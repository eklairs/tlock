package tlockcore

import (
	"log"
	"os"
	"path"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v2"
)

var USERS_PATH = path.Join(xdg.DataHome, "tlock", "users.yaml")

// Represents a user
type UserSpec struct {
    Username string
    Vault string
}

// Users API for tlock
type TLockUsers struct {
    Users []UserSpec
}

// Loads the list of users
func LoadTLockUsers() TLockUsers {
    raw, err := os.ReadFile(USERS_PATH);

    if err != nil {
        return TLockUsers{}
    }

    users := []UserSpec{}

    if err = yaml.Unmarshal(raw, users); err != nil {
        return TLockUsers{}
    }

    return TLockUsers{ Users: users }
}

// Writes the current users value to the file
func (users TLockUsers) write() {
    data, _ := yaml.Marshal(users.Users)

    file, err := os.Create(path.Join(xdg.DataHome, "tlock", "users.yaml"))

    if err != nil {
        log.Fatalf("Failed to create users file")
    }

    file.Write(data)
}

// Adds a new user
func (users *TLockUsers) AddNewUser(username, password string) {
    users.Users = append(users.Users, UserSpec{ Username: username })

    users.write()
}
