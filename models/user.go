package models

type User struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
