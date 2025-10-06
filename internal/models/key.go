package models

// Key represents a database record for a key
type Key struct {
	ID   int64
	HKey string
	MKey string
	LKey string
}

// KeyResponse represents the response structure for key operations
type KeyResponse struct {
	SortedKeys string
	AllKeys    string
}
