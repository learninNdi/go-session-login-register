package main

import (
	"fmt"
	"net/http"

	authController "github.com/learninNdi/go-session-login-register/controllers"
)

func main() {
	http.HandleFunc("/", authController.Index)
	http.HandleFunc("/login", authController.Login)
	http.HandleFunc("/logout", authController.Logout)
	http.HandleFunc("/register", authController.Register)

	fmt.Println("Server running on: http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
