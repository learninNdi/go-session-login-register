package controllers

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/learninNdi/go-session-login-register/config"
	"github.com/learninNdi/go-session-login-register/entities"
	"github.com/learninNdi/go-session-login-register/libraries"
	"github.com/learninNdi/go-session-login-register/models"
	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

var userModel = models.NewUserModel()
var validation = libraries.NewValidation()

func Index(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, config.SESSION_ID)

	if len(session.Values) == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		if session.Values["loggedIn"] != true {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			data := map[string]interface{}{
				"fullname": session.Values["fullname"],
			}

			temp, _ := template.ParseFiles("views/index.html")

			temp.Execute(w, data)
		}
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp, _ := template.ParseFiles("views/login.html")

		temp.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		// proses login
		r.ParseForm()

		userInput := &UserInput{
			Username: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		}

		//validation
		errorMessages := validation.Struct(userInput)

		if errorMessages != nil {
			data := map[string]interface{}{
				"validation": errorMessages,
			}

			temp, _ := template.ParseFiles("views/login.html")

			temp.Execute(w, data)
		} else {
			var user entities.User
			userModel.Where(&user, "username", userInput.Username)

			var message error

			if user.Username == "" {
				message = errors.New("username salah")
			} else {
				errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))

				if errPassword != nil {
					message = errors.New("password salah")
				}
			}

			if message != nil {
				data := map[string]interface{}{
					"error": message,
				}

				temp, _ := template.ParseFiles("views/login.html")

				temp.Execute(w, data)
			} else {
				// set session
				session, _ := config.Store.Get(r, config.SESSION_ID)

				session.Values["loggedIn"] = true
				session.Values["username"] = user.Username
				session.Values["fullname"] = user.FullName
				session.Values["email"] = user.Email

				session.Save(r, w)

				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, config.SESSION_ID)

	// delete session
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp, _ := template.ParseFiles("views/register.html")

		temp.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		// get form
		r.ParseForm()

		user := entities.User{
			Username:  r.Form.Get("username"),
			FullName:  r.Form.Get("fullname"),
			Email:     r.Form.Get("email"),
			Password:  r.Form.Get("password"),
			Cpassword: r.Form.Get("cpassword"),
		}

		// // VALIDASI MANUAL
		// errorMessage := make(map[string]interface{})

		// if user.Username == "" {
		// 	errorMessage["username"] = "Username harus diisi"
		// }

		// if user.FullName == "" {
		// 	errorMessage["fullname"] = "Nama lengkap harus diisi"
		// }

		// if user.Email == "" {
		// 	errorMessage["email"] = "Email harus diisi"
		// }

		// if user.Password == "" {
		// 	errorMessage["password"] = "Password harus diisi"
		// }

		// if user.CPassword == "" {
		// 	errorMessage["cpassword"] = "Konfirmasi password harus diisi"
		// } else {
		// 	if user.CPassword != user.Password {
		// 		errorMessage["cpassword"] = "Konfirmasi password tidak cocok"
		// 	}
		// }

		// if len(errorMessage) > 0 {
		// 	data := map[string]interface{}{
		// 		"validation": errorMessage,
		// 	}

		// 	temp, _ := template.ParseFiles("views/register.html")

		// 	temp.Execute(w, data)
		// } else {
		// 	// hashing password
		// 	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		// 	user.Password = string(hashedPassword)

		// 	// insert to db
		// 	_, err := userModel.Create(user)

		// 	var message string
		// 	if err != nil {
		// 		message = "Proses registrasi gagal" + err.Error()
		// 	} else {
		// 		message = "Registrasi berhasil, silakan login"
		// 	}

		// 	data := map[string]interface{}{
		// 		"pesan": message,
		// 	}

		// 	temp, _ := template.ParseFiles("views/register.html")

		// 	temp.Execute(w, data)
		// }

		errorMessages := validation.Struct(user)

		if errorMessages != nil {
			data := map[string]interface{}{
				"validation": errorMessages,
				"user":       user,
			}

			temp, _ := template.ParseFiles("views/register.html")

			temp.Execute(w, data)
		} else {
			// hashing password
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			user.Password = string(hashedPassword)

			// insert to db
			_, err := userModel.Create(user)

			var message string

			if err != nil {
				message = "Proses registrasi gagal" + err.Error()
			} else {
				message = "Registrasi berhasil, silakan login"
			}

			data := map[string]interface{}{
				"pesan": message,
			}

			temp, _ := template.ParseFiles("views/register.html")

			temp.Execute(w, data)
		}
	}
}
