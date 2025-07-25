package domain

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	ID           string
	Username     string
	Password     string
	PasswordHash string
	Role         string
}
