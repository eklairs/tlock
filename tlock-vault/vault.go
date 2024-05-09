package tlockvault

// Vault securely stores all the tokens inside of the file for tlock
type Vault struct {
	// All the folders and their data
	Folders []Folder

	// Path to the file
	path string

	// Password
	password string

	// Channel to send the data to be written
	dataChan chan []Folder
}

// Sends the data to be written to the channel
func (vault Vault) write() {
	// Clear any existing data
	select {
	case <-vault.dataChan:
	default:
	}

	// Send the new data to write
	vault.dataChan <- vault.Folders
}

// Updates the password for the vault
func (vault *Vault) ChangePassword(password string) {
	// Set the master password
	vault.password = password

	// Rewrite
	vault.write()
}
