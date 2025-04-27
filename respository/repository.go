package repository

// import "github.com/googleapis/enterprise-certificate-proxy/client"

import (
	"context"
	// "fmt"
	"log"
	"net/http"

	domain "blog/models"

	"github.com/hasura/go-graphql-client"
)

// CreateUser creates a new user with the given email and password.
func (r *userRepository) CreateUser(user *domain.User) error {

	var m struct {
		InsertUsersOne struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		} `graphql:"insert_users_one(object:{ name: $name, email: $email, password: $password})"`
	}

	variables := map[string]interface{}{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	}

	err := r.client.Mutate(context.Background(), &m, variables)

	if err != nil {
		log.Println("Error creating user:", err)
		return err
	}

	return nil
}

// GetUserByEmail retrieves a user by their email.
func (r *userRepository) GetUserByEmail(email string) (*domain.User, error) {

	var query struct {
		Users []struct {
			Name     string `graphql:"name"`
			Email    string `graphql:"email"`
			ID       string `graphql:"id"`
			Password string `graphql:password`
		} `graphql:"users(where: {email: {_eq: $email}})"`
	}

	variables := map[string]interface{}{
		"email": email,
	}
	
	err := r.client.Query(context.Background(), &query, variables)
	if err != nil {
		log.Println("Error fetching user by email:", err)
		return nil, err
	}

	if len(query.Users) != 0 {
		return &domain.User{
			ID:       query.Users[0].ID,
			Name:     query.Users[0].Name,
			Email:    query.Users[0].Email,
			Password: query.Users[0].Password,
		}, nil
	}


	return nil, nil
}

type userRepository struct {
	client      *graphql.Client
	adminSecret string
}
type UserRepository interface {
	CreateUser(user *domain.User) error
	GetUserByEmail(email string) (*domain.User, error)
}

func NewUserRepository(HasuraEndpoint, adminSecret string) UserRepository {

	return &userRepository{
		client: graphql.NewClient(HasuraEndpoint, http.DefaultClient).
			WithRequestModifier(func(r *http.Request) {
				r.Header.Set("x-hasura-admin-secret", adminSecret)
			}),
		adminSecret: adminSecret,
	}
}
