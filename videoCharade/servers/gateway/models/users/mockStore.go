package users

import (
	"errors"
	"fmt"
)

//MockStore is fake
type MockStore struct {
}

//NewMockStore generates a fake store
func NewMockStore() *MockStore {
	return &MockStore{}
}

//SetUpNewUser creates a new user and converts it to a user
func SetUpNewUser() *User {
	newUser := &NewUser{
		Email:        "john@doe.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "johndoe123",
		FirstName:    "john",
		LastName:     "doe",
	}
	user, err := newUser.ToUser()
	if err != nil {
		fmt.Printf("unexpected error: %v", err)
		return nil
	}
	return user
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (s *MockStore) Insert(u *User) (*User, error) {
	return u, nil
}

//GetByID returns the User with the given ID
func (s *MockStore) GetByID(id int64) (*User, error) {
	if id != int64(0) {
		return nil, errors.New("Fake Error")
	}
	return SetUpNewUser(), nil
}

//GetByEmail returns the User with the given email
func (s *MockStore) GetByEmail(email string) (*User, error) {
	user := SetUpNewUser()
	if user.Email != email {
		return nil, errors.New("not matching email")
	}
	return SetUpNewUser(), nil
}

//GetByUserName returns the User with the given Username
func (s *MockStore) GetByUserName(username string) (*User, error) {
	return SetUpNewUser(), nil
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (s *MockStore) Update(id int64, updates *Updates) (*User, error) {
	if len(updates.FirstName) < 1 && len(updates.LastName) < 1 {
		return nil, errors.New("fake errors")
	}
	return SetUpNewUser(), nil
}

//Delete deletes the user with the given ID
func (s *MockStore) Delete(id int64) error {
	return nil
}

//SignInInsert Inserts a user sign in log
func (s *MockStore) SignInInsert(log *SignInLog) (*SignInLog, error) {
	return nil, nil
}

//GetAllSignInByUserID returns a particulars users list of signins
func (s *MockStore) GetAllSignInByUserID(userID int64) ([]*SignInLog, error) {
	return nil, nil
}
