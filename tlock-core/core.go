package tlockcore

import "log"

// Core for tlock
type TLockCore struct {
    // Users API
    Users TLockUsers
}

// Initializes a new core instance of tlock
func New() TLockCore {
    log.Printf("[core] Initializing a new instance of tlock core")

    return TLockCore{
        Users: LoadTLockUsers(),
    }
}

