package handlers

import (
	"bytes"
	"encoding/json"
	"final-project-crew/videoCharade/servers/gateway/models/users"
	"final-project-crew/videoCharade/servers/gateway/sessions"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setUpContext() *Context {
	return NewContext("test",
		sessions.NewMemStore(time.Hour, time.Minute),
		users.NewMockStore())

}

//TestUsersHandler DONE!!
func TestUsersHandler(t *testing.T) {
	cases := []struct {
		name               string
		method             string
		contentType        string
		user               *users.NewUser
		expectedStatusCode int
		isWrong            bool
	}{
		{
			"Invalid Method",
			http.MethodGet,
			"",
			&users.NewUser{
				Email:        "test@test.com",
				Password:     "test1234",
				PasswordConf: "test1234",
				UserName:     "test",
				FirstName:    "test",
				LastName:     "test",
			},
			http.StatusMethodNotAllowed,
			false,
		},
		{
			"Invalid Header",
			http.MethodPost,
			"text/plain",
			&users.NewUser{
				Email:        "test@test.com",
				Password:     "test1234",
				PasswordConf: "test1234",
				UserName:     "test",
				FirstName:    "test",
				LastName:     "test",
			},
			http.StatusUnsupportedMediaType,
			false,
		},
		{
			"Valid Method and Header",
			http.MethodPost,
			"application/json",
			&users.NewUser{
				Email:        "test@test.com",
				Password:     "test1234",
				PasswordConf: "test1234",
				UserName:     "test",
				FirstName:    "test",
				LastName:     "test",
			},
			http.StatusCreated,
			false,
		},
		{
			"Invalid new user",
			http.MethodPost,
			"application/json",
			&users.NewUser{
				Email:        "test@test.com",
				Password:     "test1234",
				PasswordConf: "4321",
				UserName:     "test",
				FirstName:    "test",
				LastName:     "test",
			},
			http.StatusBadRequest,
			false,
		},
		{
			"wrong signing key",
			http.MethodPost,
			"application/json",
			&users.NewUser{
				Email:        "test@test.com",
				Password:     "test1234",
				PasswordConf: "test1234",
				UserName:     "test",
				FirstName:    "test",
				LastName:     "test",
			},
			http.StatusInternalServerError,
			true,
		},
		{
			"Email already in use",
			http.MethodPost,
			"application/json",
			&users.NewUser{
				Email:        "john@doe.com",
				Password:     "test1234",
				PasswordConf: "test1234",
				UserName:     "test",
				FirstName:    "john",
				LastName:     "doe",
			},
			http.StatusBadRequest,
			false,
		},
	}

	for _, c := range cases {
		data, _ := json.Marshal(c.user)
		req, _ := http.NewRequest(c.method, "/v1/user", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", c.contentType)
		respRec := httptest.NewRecorder()
		ctx := setUpContext()
		if c.isWrong {
			ctx.SigningKey = ""
		}
		ctx.UsersHandler(respRec, req)
		resp := respRec.Result()
		if resp.StatusCode != c.expectedStatusCode {
			t.Errorf("case %s: incorrect status code: expected %d but got %d",
				c.name, c.expectedStatusCode, resp.StatusCode)
		}
	}
}

//TestSpeceficUserHandler DONE!!
func TestSpecificUserHandler(t *testing.T) {
	ctx := setUpContext()
	newSessionState := &SessionState{
		time.Now(),
		users.SetUpNewUser(),
	}
	sid, err := sessions.NewSessionID(ctx.SigningKey)
	if err != nil {
		t.Fatal("unexpected error: ", err)
		return
	}
	err = ctx.SessionStore.Save(sid, newSessionState)
	if err != nil {
		t.Fatal("unexpected error: ", err)
		return
	}
	cases := []struct {
		name               string
		method             string
		url                string
		header             string
		authorization      string
		update             *users.Updates
		expectedStatusCode int
		isWrong            bool
	}{
		{
			"Invalid Session ID",
			"",
			"/v1/users/me",
			"application/json",
			"something",
			&users.Updates{
				FirstName: "Johnny",
				LastName:  "Dough",
			},
			http.StatusUnauthorized,
			false,
		},
		{
			"Invalid url segment",
			"",
			"/v1/users/john",
			"application/json",
			"Bearer " + sid.String(),
			&users.Updates{
				FirstName: "Johnny",
				LastName:  "Dough",
			},
			http.StatusBadRequest,
			false,
		},
		{
			"Invalid URL Request",
			"",
			"/v1/users/me/1234",
			"application/json",
			"something",
			&users.Updates{
				FirstName: "Johnny",
				LastName:  "Dough",
			},
			http.StatusUnauthorized,
			false,
		},
		{
			"Used GET Method and Requested well Using me",
			"GET",
			"/v1/users/me",
			"application/json",
			"Bearer " + sid.String(),
			&users.Updates{
				FirstName: "Johnny",
				LastName:  "Dough",
			},
			http.StatusOK,
			false,
		},
		{
			"Used PATCH Method and Valid 'me' url segment",
			"PATCH",
			"/v1/users/me",
			"application/json",
			"Bearer " + sid.String(),
			&users.Updates{
				FirstName: "Johnny",
				LastName:  "Dough",
			},
			http.StatusOK,
			false,
		},
		{
			"invalid GET user id url segment",
			http.MethodGet,
			"/v1/users/23",
			"application/json",
			"Bearer " + sid.String(),
			&users.Updates{
				FirstName: "Johnny",
				LastName:  "Dough",
			},
			http.StatusNotFound,
			false,
		},
		{
			"invalid PATCH user id url segment",
			http.MethodPatch,
			"/v1/users/23",
			"application/json",
			"Bearer " + sid.String(),
			&users.Updates{
				FirstName: "Johnny",
				LastName:  "Dough",
			},
			http.StatusForbidden,
			false,
		},
		{
			"GET Method and Valid Requested Using numberID",
			http.MethodGet,
			"/v1/users/0",
			"application/json",
			"Bearer " + sid.String(),
			&users.Updates{
				FirstName: "Johnny",
				LastName:  "Dough",
			},
			http.StatusOK,
			false,
		},
		{
			"PATCH Method and Valid Requested Using numberID",
			http.MethodPatch,
			"/v1/users/0",
			"application/json",
			"Bearer " + sid.String(),
			&users.Updates{
				FirstName: "Johnny",
				LastName:  "Dough",
			},
			http.StatusOK,
			false,
		},
		{
			"Invalid Header",
			http.MethodPatch,
			"/v1/users/0",
			"plain/text",
			"Bearer " + sid.String(),
			&users.Updates{
				FirstName: "Johnny",
				LastName:  "Dough",
			},
			http.StatusUnsupportedMediaType,
			false,
		},
		{
			"Invalid Method",
			http.MethodPost,
			"/v1/users/0",
			"plain/text",
			"Bearer " + sid.String(),
			&users.Updates{
				FirstName: "Johnny",
				LastName:  "Dough",
			},
			http.StatusMethodNotAllowed,
			false,
		},
		{
			"Invalid Update",
			http.MethodPatch,
			"/v1/users/0",
			"application/json",
			"Bearer " + sid.String(),
			&users.Updates{
				FirstName: "",
				LastName:  "",
			},
			http.StatusBadRequest,
			false,
		},
	}
	for _, c := range cases {
		data, _ := json.Marshal(c.update)
		req, _ := http.NewRequest(c.method, c.url, bytes.NewBuffer(data))
		req.Header.Set("Content-Type", c.header)
		respRec := httptest.NewRecorder()
		req.Header.Set("Authorization", c.authorization)
		if c.isWrong {
			ctx.SigningKey = ""
		}
		ctx.SpecificUserHandler(respRec, req)
		resp := respRec.Result()
		if resp.StatusCode != c.expectedStatusCode {
			t.Errorf("case %s: incorrect status code: expected %d but got %d",
				c.name, c.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestSessionsHandler(t *testing.T) {

	cases := []struct {
		name               string
		method             string
		header             string
		credential         *users.Credentials
		expectedStatusCode int
		isWrong            bool
	}{
		{
			"Invalid Method",
			http.MethodGet,
			"application/json",
			&users.Credentials{
				Email:    "test@test.com",
				Password: "password",
			},
			http.StatusMethodNotAllowed,
			false,
		},
		{
			"Invalid Content-Type",
			http.MethodPost,
			"text/plain",
			&users.Credentials{
				Email:    "test@test.com",
				Password: "password",
			},
			http.StatusUnsupportedMediaType,
			false,
		},
		{
			"Invalid Email",
			http.MethodPost,
			"application/json",
			&users.Credentials{
				Email:    "wrong@email.com",
				Password: "password",
			},
			http.StatusUnauthorized,
			false,
		},
		{
			"Valid Request",
			http.MethodPost,
			"application/json",
			&users.Credentials{
				Email:    "john@doe.com",
				Password: "password",
			},
			http.StatusCreated,
			false,
		},
		{
			"Valid Request",
			http.MethodPost,
			"application/json",
			&users.Credentials{
				Email:    "john@doe.com",
				Password: "password",
			},
			http.StatusInternalServerError,
			true,
		},
		{
			"Wrong Password",
			http.MethodPost,
			"application/json",
			&users.Credentials{
				Email:    "john@doe.com",
				Password: "wrongpassword",
			},
			http.StatusUnauthorized,
			false,
		},
	}
	for _, c := range cases {
		data, _ := json.Marshal(c.credential)
		req, _ := http.NewRequest(c.method, "/v1/user", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", c.header)
		respRec := httptest.NewRecorder()
		ctx := setUpContext()
		if c.isWrong {
			ctx.SigningKey = ""
		}
		ctx.SessionsHandler(respRec, req)
		resp := respRec.Result()
		if resp.StatusCode != c.expectedStatusCode {
			t.Errorf("case %s: incorrect status code: expected %d but got %d",
				c.name, c.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestSpecificSessionHandler(t *testing.T) {
	ctx := setUpContext()
	newSessionState := &SessionState{
		time.Now(),
		users.SetUpNewUser(),
	}
	sid, err := sessions.NewSessionID(ctx.SigningKey)
	if err != nil {
		t.Fatal("unexpected error: ", err)
		return
	}
	err = ctx.SessionStore.Save(sid, newSessionState)
	if err != nil {
		t.Fatal("unexpected error: ", err)
		return
	}
	cases := []struct {
		name               string
		method             string
		url                string
		expectedStatusCode int
		isWrong            bool
	}{
		{
			"Invalid Method",
			http.MethodGet,
			"/v1/user/mine",
			http.StatusMethodNotAllowed,
			false,
		},
		{
			"Invalid URL",
			http.MethodDelete,
			"/v1/user/me",
			http.StatusForbidden,
			false,
		},
		{
			"Valid Method and URL",
			http.MethodDelete,
			"/v1/user/mine",
			http.StatusOK,
			false,
		},
		{
			"invalid signing key",
			http.MethodDelete,
			"/v1/user/mine",
			http.StatusInternalServerError,
			true,
		},
	}
	for _, c := range cases {
		req, _ := http.NewRequest(c.method, c.url, nil)
		req.Header.Set("Authorization", "Bearer "+sid.String())
		respRec := httptest.NewRecorder()
		if c.isWrong {
			ctx.SigningKey = ""
		}
		ctx.SpecificSessionHandler(respRec, req)
		resp := respRec.Result()
		if resp.StatusCode != c.expectedStatusCode {
			t.Errorf("case %s: incorrect status code: expected %d but got %d",
				c.name, c.expectedStatusCode, resp.StatusCode)
		}
	}
}
