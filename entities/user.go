package entities

type User struct {
	ID        int64
	Username  string `validate:"required,gte=3,isUnique=users-username"`
	FullName  string `validate:"required" label:"Nama lengkap"`
	Email     string `validate:"required,email,isUnique=users-email"`
	Password  string `validate:"required,gte=6"`
	Cpassword string `validate:"required,eqfield=Password" label:"Konfirmasi password"`
}
