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

// Returns a list of all the available users
func (core TLockCore) GetUsers() []string {
    return []string {}
}

