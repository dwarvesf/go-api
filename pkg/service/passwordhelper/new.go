package passwordhelper

// Helper password helper
type Helper interface {
	GenerateSalt() string
	Hash(password, salt string) (string, error)
	Compare(password, hashedPassword, salt string) bool
}

// NewArgon2 init argon2 helper
func NewArgon2() Helper {
	return newArgon2Default()
}
