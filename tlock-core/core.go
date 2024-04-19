package tlockcore

// Core for tlock
type TLockCore struct {
    // Users API
    Users TLockUsers
}

// Initializes a new core instance of tlock
func New() TLockCore {
    return TLockCore{
        Users: LoadTLockUsers(),
    }
}

