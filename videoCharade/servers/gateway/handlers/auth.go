package handlers

import (
	"encoding/json"
	"fmt"

	"final-project-crew/videoCharade/servers/gateway/models/users"
	"final-project-crew/videoCharade/servers/gateway/sessions"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//UsersHandler handles requests for the "users" resource. For now we will accept POST requests to create new user accounts, but in the next assignment, this will be extended to also handle GET requests for searching users.
func (ctx *Context) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "Not valid Content-Type", http.StatusUnsupportedMediaType)
			return
		}
		newUser := &users.NewUser{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&newUser); err != nil {
			http.Error(w, "Must be JSON", http.StatusBadRequest)
			return
		}

		_, err := ctx.UserStore.GetByEmail(newUser.Email)
		if err == nil {
			http.Error(w, "An account is already set up with this email", http.StatusBadRequest)
			return
		}

		_, unErr := ctx.UserStore.GetByUserName(newUser.UserName)
		if unErr == nil {
			http.Error(w, "An account is already set up with this user name", http.StatusBadRequest)
			return
		}

		user, err := newUser.ToUser()
		if err != nil {
			http.Error(w, "ToUser Request: "+err.Error(), http.StatusBadRequest)
			return
		}

		insertedUser, err := ctx.UserStore.Insert(user)
		if err != nil {
			http.Error(w, "Insert Error"+err.Error(), http.StatusInternalServerError)
			return
		}

		sesstionState := SessionState{time.Now(), insertedUser}
		_, err = sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, sesstionState, w)
		if err != nil {
			http.Error(w, "Session Error"+err.Error(), http.StatusInternalServerError)
			return
		}
		//add inserted user to trie
		ctx.UserTrie.Add(strings.ToLower(insertedUser.FirstName), insertedUser.ID)
		ctx.UserTrie.Add(strings.ToLower(insertedUser.LastName), insertedUser.ID)
		ctx.UserTrie.Add(strings.ToLower(insertedUser.UserName), insertedUser.ID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		data, err := json.Marshal(insertedUser)
		if err == nil {
			w.Write(data)
		} else {
			w.Write([]byte("[]"))
		}
	} else if r.Method == http.MethodGet {
		var state SessionState
		_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, &state)
		if err != nil {
			http.Error(w, "Status Unauthorized", http.StatusUnauthorized)
			return
		}
		query := r.FormValue("q")
		if len(query) < 1 {
			http.Error(w, "Query can't be empty. search something...", http.StatusBadRequest)
			return
		}
		idList := ctx.UserTrie.Find(query, 20)
		userList := []*users.User{}
		for _, id := range idList {
			user, error := ctx.UserStore.GetByID(id)
			if error != nil {
				http.Error(w, "Server Error", http.StatusInternalServerError)
				return
			}
			userList = append(userList, user)
		}
		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(userList)
		if err == nil {
			w.Write(data)
		} else {
			w.Write([]byte("[]"))
		}
	} else {
		http.Error(w, "MethodNotAllowed", http.StatusMethodNotAllowed)
		return
	}
}

//SpecificUserHandler handles requests for a specific user. The resource path will be /v1/users/{UserID}, where {UserID} will be the user's ID. As a convenience, clients can also use the literal string me to refer to the UserID of the currently-authenticated user. So /v1/users/me refers to the currently-authenticated user, and /v1/users/1234 refers to the user with the ID 1234 (which could be the same user).
func (ctx *Context) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	//get SessionState
	var state SessionState
	_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, &state)
	if err != nil {
		http.Error(w, "Status Unauthorized", http.StatusUnauthorized)
		return
	}
	//Get ID from URL SEGMENT
	segment := getLastURLSegment(r.URL.Path)
	var userID int64
	var userErr error
	if segment == "me" {
		userID = state.User.ID
	} else {
		userID, userErr = strconv.ParseInt(segment, 10, 64)
		if userErr != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}
	if r.Method == http.MethodGet {
		user, err := ctx.UserStore.GetByID(userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		data, err := json.Marshal(user)
		if err == nil {
			w.Write(data)
		} else {
			w.Write([]byte("[]"))
		}
	} else if r.Method == http.MethodPatch {
		if segment != "me" && userID != state.User.ID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "Not valid Content-Type", http.StatusUnsupportedMediaType)
			return
		}
		updates := &users.Updates{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(updates); err != nil {
			http.Error(w, "Must be JSON", http.StatusBadRequest)
			return
		}
		user, err := ctx.UserStore.Update(userID, updates)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		//remove user names
		ctx.UserTrie.Remove(state.User.FirstName, state.User.ID)
		ctx.UserTrie.Remove(state.User.LastName, state.User.ID)

		//add user updated names
		ctx.UserTrie.Add(strings.ToLower(user.FirstName), user.ID)
		ctx.UserTrie.Add(strings.ToLower(user.LastName), user.ID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		data, err := json.Marshal(user)
		if err == nil {

			w.Write(data)
		} else {
			w.Write([]byte("[]"))
		}
	} else {
		http.Error(w, "MethodNotAllowed", http.StatusMethodNotAllowed)
		return
	}
}

//SessionsHandler handles requests for the "sessions" resource, and allows clients to begin a new session using an existing user's credentials.
func (ctx *Context) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "Not valid Content-Type", http.StatusUnsupportedMediaType)
			return
		}
		credentials := &users.Credentials{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(credentials); err != nil {
			http.Error(w, "Must be JSON", http.StatusBadRequest)
			return
		}
		user, err := ctx.UserStore.GetByEmail(credentials.Email)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		err = user.Authenticate(credentials.Password)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		ipAddress := strings.Split(r.Header.Get("X-Forwarded-For"), ", ")[0]
		if len(ipAddress) < 1 {
			ipAddress = r.RemoteAddr
		}
		dateTime := time.Now()
		userSignIn := &users.SignInLog{
			UserID:    user.ID,
			DateTime:  dateTime,
			IPAddress: ipAddress,
		}
		ctx.UserStore.SignInInsert(userSignIn)
		sesstionState := SessionState{dateTime, user}
		_, err = sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, sesstionState, w)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		data, err := json.Marshal(user)
		if err == nil {
			w.Write(data)
		} else {
			w.Write([]byte("[]"))
		}
	} else {
		http.Error(w, "Server Error", http.StatusMethodNotAllowed)
		return
	}
}

//SpecificSessionHandler handles requests related to a specific authenticated session. For now, the only operation we will support is ending the current user's session. But this could be expanded to allow administrators to end sessions started by other users that have gone rogue. The resource path for the current user's session will be /v1/sessions/mine.
func (ctx *Context) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		segment := getLastURLSegment(r.URL.Path)
		if segment != "mine" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		var state SessionState
		_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, &state)
		if err != nil {
			http.Error(w, "Status Unauthorized", http.StatusUnauthorized)
			return
			fmt.Println(err)
		}

		_, err = sessions.EndSession(r, ctx.SigningKey, ctx.SessionStore)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("signed out"))
	} else {
		http.Error(w, "MethodNotAllowed", http.StatusMethodNotAllowed)
		return
	}
}

//get the last segment in URL
func getLastURLSegment(path string) string {
	urlSlice := strings.Split(path, "/")
	urlSeg := urlSlice[len(urlSlice)-1]
	return urlSeg
}

//Get ID from state
//userID := state.User.ID
//ctx.RemoveConnection(userID)
