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

	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {

		return "", "", err
	}
	fmt.Println("check point 2", "user")
	if user != nil {
		return "user already exist", "", nil
	}
	hashpassword,err := utils.HashPassword(password)
	use := &domain.User{
		Email:    email,
		Name:     name,
		Password: hashpassword,
	}
	fmt.Println("check point 3", "user", use)
	err = u.userRepo.CreateUser(use)
	if err != nil {
		return "", "", err
	}
	fmt.Println("check point 4", "user", use)
	accessToken, err := utils.GenerateAccessToken(use.ID, use.Email, u.jwtSecret)
	if err != nil {
		return "", "", err
	}
	fmt.Println("check point 5", "user", use)
	refreshToken, err := utils.GenerateRefreshToken(use.ID, email, u.jwtSecret)
	if err != nil {
		return "", "", err
	}
	fmt.Println("check point 6", "user", use)
	return accessToken, refreshToken, nil

}
func (u *authUsecase) Login(email, password string) (string, string, error) {

	user, err := u.userRepo.GetUserByEmail(email)
	print(user)
	if err != nil {
		return "", "", err
	}
	if user.ID == "" {
		return "", "", nil
	}
	fmt.Println("check new 1",password,user.Password)

	if !utils.CheckPasswordHash(password, user.Password) {
		return "checkthe password", "", nil
	}
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, u.jwtSecret)
	if err != nil {
		return "check the accesstoken", "", err
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID, email, u.jwtSecret)
	if err != nil {
		return "check the refresh token", "", err
	}
	return accessToken, refreshToken, nil

	// return "", "", nil
}
func (u *authUsecase) RefreshToken(refreshToken string) (string, error) {
	claims, err := utils.ValidateToken(refreshToken, u.jwtSecret)
	if err != nil {
		return "", err
	}
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
	return "", nil
}
