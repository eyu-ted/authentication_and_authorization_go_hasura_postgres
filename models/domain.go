package domain

type User struct {
	ID       string    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	LastSeen        string `json:"last_seen"`
	IsEmailVerified bool   `json:"is_email_verified"`
	RefreshToken		string `json:"refresh_token"`
}
