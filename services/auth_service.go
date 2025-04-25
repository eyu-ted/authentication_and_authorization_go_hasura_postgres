package services

import (
	domain "blog/models"
	repository "blog/respository"
	"blog/utils"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type AuthUsecase interface {
	Signup(name, email, password string) (string, string, error)
	Login(email, password string) (string, string, error)
	RefreshToken(refreshToken string) (string, error)
}
type authUsecase struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthUsecase(userRepo repository.UserRepository, jwtSecret string) AuthUsecase {
	return &authUsecase{userRepo, jwtSecret}
}
func (u *authUsecase) Signup(name, email, password string) (string, string, error) {
	fmt.Println("check point 1", name, email, password)
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {

		return "", "", err
	}
	if user != nil {
		return "user already exist", "", nil
	}
	hashpassword,err := utils.HashPassword(password)
	use := &domain.User{
		Email:    email,
		Name:     name,
		Password: hashpassword,
	}

	err = u.userRepo.CreateUser(use)
	if err != nil {
		return "", "", err
	}

	accessToken, err := utils.GenerateAccessToken(use.ID, use.Email, u.jwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateRefreshToken(use.ID, email, u.jwtSecret)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil

}
func (u *authUsecase) Login(email, password string) (string, string, error) {

	user, err := u.userRepo.GetUserByEmail(email)
	print()
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", nil
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", "", nil
	}
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, u.jwtSecret)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID, email, u.jwtSecret)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil


}
func (u *authUsecase) RefreshToken(refreshToken string) (string, error) {
	claims, err := utils.ValidateToken(refreshToken, u.jwtSecret)
	if err != nil {
		return "", err
	}
	fmt.Println("check point refresh 1", "claims", claims)
	email, ok := claims.Claims.(jwt.MapClaims)["email"].(string)
	if !ok {
		return "", errors.New("invalid token claims")
	}
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if user.ID == "" {
		return "", nil
	}
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, u.jwtSecret)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
