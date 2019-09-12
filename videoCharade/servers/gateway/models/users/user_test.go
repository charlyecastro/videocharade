package users

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestValidate(t *testing.T) {
	//function to ensure it catches all possible validation errors, and returns no error when the new user is valid.
	cases := []struct {
		name        string
		hint        string
		newUser     NewUser
		expectError bool
	}{
		{
			"Invalid New User Email",
			"Remember to return an error if Email is invalid",
			NewUser{"charlye.com", "password", "password", "charles", "charlye", "castro"},
			true,
		},
		{
			"Invalid New User Password",
			"Remember to return a error if passwrod is less than 6 characters",
			NewUser{"charlye@gmail.com", "pass", "pass", "charles", "charlye", "castro"},
			true,
		},
		{
			"Invalid New User Password Confirmation",
			"Remember to return a error if password and passwordConf don't match",
			NewUser{"charlye@gmail.com", "password", "pass", "charles", "charlye", "castro"},
			true,
		},
		{
			"Invalid New User UserName (space between)",
			"Remember to return a error if UserName contains any spaces",
			NewUser{"charlye@gmail.com", "password", "password", "char les", "charlye", "castro"},
			true,
		},
		{
			"Invalid New User UserName (Empty)",
			"Remember to return a error if UserName is empty",
			NewUser{"charlye@gmail.com", "password", "password", "", "charlye", "castro"},
			true,
		},
		{
			"Valid Everything",
			"Remember to return nil if there are no errors",
			NewUser{"charlye@gmail.com", "password", "password", "charles", "charlye", "castro"},
			false,
		},
	}
	for _, c := range cases {
		err := c.newUser.Validate()
		if err != nil && !c.expectError {
			t.Errorf("case %s: unexpected error in authenticate: %v\nHINT: %s", c.name, err, c.hint)
		}

		//if user.PassHash
	}
}

func TestToUser(t *testing.T) {
	//function to ensure it calculates the PhotoURL field correctly, even when the email address has upper case letters or spaces (Links to an external site.)Links to an external site., and sets the PassHash field to the password hash. Since bcrypt hashes are salted with a random value, you can't anticipate what the hash should be, but you can verify the generated hash by comparing it to the original password using the bcrypt package functions.
	const expectedGravatar = "https://www.gravatar.com/avatar/277b890ff823f14d2079d4fcbec3cc2f"

	cases := []struct {
		name       string
		hint       string
		newUser    NewUser
		expectGrav string
	}{
		{
			"Normal Email",
			"Remember to use gravatarBasePhotoURL",
			NewUser{"charlyecastro@gmail.com", "password", "password", "charles", "charlye", "castro"},
			expectedGravatar,
		},
		{
			"Capitalized Email",
			"Remember to to convert email to lower case",
			NewUser{"CHARLYEcastro@gmail.com", "password", "password", "charles", "charlye", "castro"},
			expectedGravatar,
		},
		{
			"Spaces within Email",
			"Remember to spaces in email",
			NewUser{"charlye castro@gmail.com", "password", "password", "charles", "charlye", "castro"},
			expectedGravatar,
		},
		{
			"Verify Password",
			"Remember to spaces in email",
			NewUser{"charlye castro@gmail.com", "password", "password", "charles", "charlye", "castro"},
			expectedGravatar,
		},
	}
	for _, c := range cases {
		user, err := c.newUser.ToUser()
		if err != nil {
			t.Errorf("case %s: unexpected error in ToUser: %v\nHINT: %s", c.name, err, c.hint)
		}

		if user.PhotoURL != c.expectGrav {
			t.Errorf("case %s: unexpected error in Getting Grav: %v\nHINT: %s", c.name, err, c.hint)
		}
		if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(c.newUser.Password)); err != nil {
			t.Errorf("case %s: unexpected error in Comparing: %v\nHINT: %s", c.name, err, c.hint)
		}
	}
}

func TestFullName(t *testing.T) {
	//function to verify that it returns the correct results given the various possible inputs (no FirstName, no LastName, neither field set, both fields set).

	userNormal := &User{
		ID:        1,
		FirstName: "charlye",
		LastName:  "castro",
		Email:     "charlyecastro@gmail.com",
		UserName:  "charles",
		PassHash:  []byte("password123"),
		PhotoURL:  "https://www.gravatar.com/avatar/sakjfddslkjfei",
	}

	userEmptyFName := &User{
		ID:        1,
		FirstName: "",
		LastName:  "castro",
		Email:     "charlyecastro@gmail.com",
		UserName:  "charles",
		PassHash:  []byte("password123"),
		PhotoURL:  "https://www.gravatar.com/avatar/sakjfddslkjfei",
	}

	userEmptyLName := &User{
		ID:        1,
		FirstName: "charlye",
		LastName:  "",
		Email:     "charlyecastro@gmail.com",
		UserName:  "charles",
		PassHash:  []byte("password123"),
		PhotoURL:  "https://www.gravatar.com/avatar/sakjfddslkjfei",
	}

	userEmpty := &User{
		ID:        1,
		FirstName: "",
		LastName:  "",
		Email:     "charlyecastro@gmail.com",
		UserName:  "charles",
		PassHash:  []byte("password123"),
		PhotoURL:  "https://www.gravatar.com/avatar/sakjfddslkjfei",
	}

	cases := []struct {
		name       string
		hint       string
		user       *User
		expectName string
	}{
		{
			"First Name and Last name",
			"Remember to not return a space when either first or last name is empty",
			userNormal,
			"charlye castro",
		},
		{
			"Empty First Name",
			"Remember to not return a space when either first or last name is empty",
			userEmptyFName,
			"castro",
		},
		{
			"Empty Last Name",
			"Remember to not return a space when either first or last name is empty",
			userEmptyLName,
			"charlye",
		},
		{
			"Empty First Name & Last Name",
			"Remember to not return an empty string if both names are empty",
			userEmpty,
			"",
		},
	}
	for _, c := range cases {
		name := c.user.FullName()
		if name != c.expectName {
			t.Errorf("case %s: unexpected error in authenticate: %v\nHINT: %s", c.name, c.hint, c.expectName)
		}
	}
}

func TestAuthenticate(t *testing.T) {
	// function to verify that authentication happens correctly for the various possible inputs (incorrect password, correct password, empty string password).

	var newUserCharlye = NewUser{"charlyecastro@gmail.com", "password", "password", "charles", "charlye", "castro"}
	var userNormal, _ = newUserCharlye.ToUser()

	cases := []struct {
		name          string
		hint          string
		passAttempt   string
		expectedMatch bool
	}{
		{
			"correct password",
			"Remember to use bcrypt.CompareHashAndPassword",
			"password",
			false,
		},
		{
			"incorrect password",
			"Remember to not return a space when either first or last name is empty",
			"wrongPassword",
			true,
		},
		{
			"Empty password",
			"Remember to not return a space when either first or last name is empty",
			"",
			true,
		},
		{
			"Mispelled Password",
			"Remember to not return an empty string if both names are empty",
			"PassWord",
			true,
		},
	}
	for _, c := range cases {
		err := userNormal.Authenticate(c.passAttempt)
		if err != nil && !c.expectedMatch {
			t.Errorf("case %s: unexpected error in authenticate: %v\nHINT: %s", c.name, err, c.hint)
		}
	}

}

func TestApplyUpdates(t *testing.T) {
	//function to ensure the user's fields are updated properly given an Updates struct.
	var newUserNormal = NewUser{"charlyecastro@gmail.com", "password", "password", "charles", "charlye", "castro"}

	var userNormal, _ = newUserNormal.ToUser()

	var updateNormal = Updates{"Heri", "Sarmiento"}
	var updateEmptyFName = Updates{"charlye", "Sarmiento"}
	var updateEmptyLName = Updates{"Heri", "castro"}
	var updateEmpty = Updates{"", ""}

	cases := []struct {
		name        string
		hint        string
		updates     Updates
		expectError bool
	}{
		{
			"First Name and Last Name",
			"Remember to not return a space when either first or last name is empty",
			updateNormal,
			false,
		},
		{
			"Empty First Name",
			"Remember to not return a space when either first or last name is empty",
			updateEmptyFName,
			false,
		},
		{
			"Empty Last Name",
			"Remember to not return a space when either first or last name is empty",
			updateEmptyLName,
			false,
		},
		{
			"Empty First Name & Last Name",
			"Remember to not return an empty string if both names are empty",
			updateEmpty,
			true,
		},
	}

	for _, c := range cases {
		err := userNormal.ApplyUpdates(&c.updates)
		if err != nil && !c.expectError {
			t.Errorf("case %s: unexpected error generating new SessionID: %v", c.name, err)
		}
	}
}
