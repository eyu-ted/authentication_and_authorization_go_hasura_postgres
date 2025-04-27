package main

import (
	"blog/config"
	// "blog/controller"
	repository "blog/respository"
	controller "blog/controller"
	"blog/services"
	"log"
	// "os"
	// "path/filepath"
	// "net/http"

	"github.com/gin-gonic/gin"
	"fmt"
	// "github.com/google/uuid"
	// "time"

)
// func GenerateUniqueFileName(originalFilename string) string {
// 	extension := filepath.Ext(originalFilename) // .jpg, .png, etc
// 	newName := uuid.New().String() + "-" + time.Now().Format("20060102150405") + extension
// 	return newName
// }

func main() {
	cfg := config.NewConfig()
	endpoint := cfg.HasuraEndpoint

	adminSecret := cfg.HasuraAdminSecret
	JWT := cfg.JwtSecret

	fmt.Println("check point 1", "endpoint", endpoint,adminSecret, JWT)



	userRepo := repository.NewUserRepository(endpoint, adminSecret)

	userUsecase := services.NewAuthUsecase(userRepo, JWT)

	userController := controller.NewAuthController(userUsecase)

	router := gin.Default()

	public := router.Group("api/auth")
	{
		public.POST("/signup", userController.Signup)
		public.POST("/login", userController.Login)
		public.POST("/refresh", userController.Refresh)
	}
	router.Static("/uploads", "./uploads")

    // Single Upload
    router.POST("/upload", controller.UploadSingleHandler)

    // Multiple Uploads
    router.POST("/upload/multiple", controller.UploadMultipleHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("this is the error:", err)
	}

	// const uploadPath = "./uploads"
	
	
	// router.POST("/upload", func(c *gin.Context) {
	// 	file, err := c.FormFile("file")
	// 	if err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
	// 		return
	// 	}

	// 	Filename := GenerateUniqueFileName(file.Filename)

	// 	dst := filepath.Join(uploadPath, Filename)
	// 	if err := c.SaveUploadedFile(file, dst); err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, gin.H{"message": "Uploaded", "fielename": file.Filename})
	// })

	// router.GET("/image/:filename", func(c *gin.Context) {
	// 	filename := c.Param("filename")
	// 	fullPath := filepath.Join(uploadPath, filename)

	// 	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
	// 		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
	// 		return
	// 	}

	// 	c.File(fullPath)
	// })

	
	// end := "https://real-troll-94.hasura.app/v1/graphql"
	// email := "eyuu"
	// client := graphql.NewClient(end, http.DefaultClient).
	// 	WithRequestModifier(func(r *http.Request) {
	// 		r.Header.Set(
	// 			"x-hasura-admin-secret", "5qKQBAx0YO4v6EA32y0mtm78F0fU4vwMF7YZLY4qMri90oymKlzVLa1eBTbfPk82")

	// 		// r.Header.Set("Authorization", "Bearer "+token)
	// 		// r.Header.Set("x-hasura-role", "user")
	// 	})

	// var query struct {
	// 	Users []struct {
	// 		Name  string `graphql:"name"`
	// 		Email string `graphql:"email"`
	// 		ID    string `graphql:"id"`
	// 	} `graphql:"users(where: {email: {_eq: $email}})"`
	// }

	// variables := map[string]interface{}{
	// 	"email": email,
	// }
	// fmt.Println("check point 7", "email", email)
	// err := client.Query(context.Background(), &query, variables)
	// if err != nil {
	// 	log.Fatalf("Error fetching user by email:", err)

	// }
	// fmt.Println("check point 8", "user", query.Users)
	// if len(query.Users) != 0 {
	// 	x := &domain.User{
	// 		ID:    query.Users[0].ID,
	// 		Name:  query.Users[0].Name,
	// 		Email: query.Users[0].Email,
	// 	}
	// 	fmt.Println("check point 9", "user", x)
	// }

	// fmt.Printf("Client: %v", client)

	// authService := authService.NewAuthService(client)
	// var q struct {
	// 	Users []struct {
	// 		Email string `graphql:"email"`
	// 		ID    string `graphql:"id"`
	// 		Name  string `graphql:"name"`
	// 	} `graphql:"users"`
	// }

	// err := client.Query(context.Background(), &q, nil)
	// if err != nil {
	// 	fmt.Printf("Error: %v", err)
	// 	return
	// }

	// fmt.Printf("Query: %v", q)

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	// 	var request struct{
	// 		Input struct {
	// 			Email string `graphql:"email"`
	// 			Password string `graphql:"password"`
	// 		} `graphql:"input"`
	// 		}

	// 	err := json.NewDecoder(r.Body).Decode(&request)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return
	// 	}

	// 	playload,err := authService.Signup(r.Context(), request.Input.Email, request.Input.Password)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// 	json.NewEncoder(w).Encode(map[string]interface{}{
	// 		"token": playload["user_id"],
	// 		"user": map[string]interface{}{
	// 			"id": playload["user_id"],
	// 			"email": request.Input.Email,
	// 		},
	// 	})

	// })
}
