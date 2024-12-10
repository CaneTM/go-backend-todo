package reps

type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Hash     string `json:"hash,omitempty"`
}
