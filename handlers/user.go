package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/TinySchoolHub/tiny-school-hub-api-backend/models"
)

var users = []models.User{}


// @Summary Get all users
// @Description Get all users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User
// @Router /users [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User
// @Router /users{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
    // Implement the logic to get a single user by ID
}

// @Summary Create a user
// @Description Create a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.CreateUserRequest true "User data"
// @Success 201 {object} models.User
// @Router /users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
    var req models.CreateUserRequest
    json.NewDecoder(r.Body).Decode(&req)

    user := models.User{
        ID:    uint(len(users) + 1), // ID is generated by the program
        Name:  req.Name,
        Email: req.Email,
    }

    users = append(users, user)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

// @Summary Update an user by ID
// @Description Update an user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User
// @Router /users{id} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
    // Implement the logic to update a user by ID
}

// @Summary Delete an user by ID
// @Description Delete an user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User
// @Router /users{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
    // Implement the logic to delete a user by ID
}