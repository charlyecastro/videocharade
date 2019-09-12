package users

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//gravatarBasePhotoURL is the base URL for Gravatar image requests.
//See https://id.gravatar.com/site/implement/images/ for details
const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"-"` //never JSON encoded/decoded
	PassHash  []byte `json:"-"` //never JSON encoded/decoded
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhotoURL  string `json:"photoURL"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	nu.Email = strings.ToLower(nu.Email)
	nu.Email = strings.Replace(nu.Email, " ", "", -1)
	_, err := mail.ParseAddress(nu.Email)
	if err != nil {
		return fmt.Errorf("Email Address is not valid")
	}
	if len(nu.Password) < 6 {
		return fmt.Errorf("Password must be at least 6 characters")
	}
	if nu.Password != nu.PasswordConf {
		return fmt.Errorf("Password does not match Password Confirmation")
	}
	if strings.Contains(nu.UserName, " ") || len(nu.UserName) < 1 {
		return fmt.Errorf("Username can not contain spaces and cannot be empty")
	}
	return nil
}

//ToUser converts the NewUser to a User, setting the
//PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	err := nu.Validate()
	if err != nil {
		return nil, err
	}

	u := &User{}

	u.ID = 0
	u.Email = nu.Email
	u.SetPassword(nu.Password)
	u.UserName = nu.UserName
	u.FirstName = nu.FirstName
	u.LastName = nu.LastName
	u.PhotoURL = getGravatar(nu.Email)

	return u, nil
}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
//If either first or last name is an empty string, no
//space is put between the names. If both are missing,
//this returns an empty string
func (u *User) FullName() string {
	space := " "
	if len(u.FirstName) < 1 || len(u.LastName) < 1 {
		space = ""
	}
	return u.FirstName + space + u.LastName
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}
	u.PassHash = hash
	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(password))
	if err != nil {
		return err
	}
	return nil
}

//ApplyUpdates takes a updates struct and applies those to the user. if both names are empty return an error othewise
//change the names that are not emty
func (u *User) ApplyUpdates(updates *Updates) error {
	if len(updates.FirstName) < 1 && len(updates.LastName) < 1 {
		return errors.New("Cannot have both fields be empty")
	}
	if len(updates.FirstName) > 0 {
		u.FirstName = updates.FirstName
	}
	if len(updates.LastName) > 0 {
		u.LastName = updates.LastName
	}
	return nil
}

func getGravatar(email string) string {
	hasher := md5.New()
	hasher.Write([]byte(email))
	emailHash := hex.EncodeToString(hasher.Sum(nil))
	return gravatarBasePhotoURL + emailHash
}
