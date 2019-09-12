package sessions

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	sid, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, err
	}

	store.Save(sid, sessionState)
	header := w.Header()
	bearer := schemeBearer + sid.String()
	header.Set(headerAuthorization, bearer)
	return sid, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	var sid string
	sid = r.Header.Get(headerAuthorization)
	if len(sid) < 1 {
		sid = r.FormValue(paramAuthorization)
		if len(sid) < 1 {
			return InvalidSessionID, ErrInvalidID
		}
	}
	if !strings.HasPrefix(sid, schemeBearer) {
		return InvalidSessionID, ErrInvalidScheme
	}
	sidArray := strings.SplitAfter(sid, schemeBearer)
	sid = sidArray[1]
	valSid, err := ValidateID(sid, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}

	return valSid, nil
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	valSid, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	storeErr := store.Get(valSid, sessionState)

	if storeErr != nil {
		return InvalidSessionID, storeErr
	}
	return valSid, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	valSid, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	err = store.Delete(valSid)
	if err != nil {
		fmt.Print("couldnt deled sid becasue of ERR: ", err)
	}
	return valSid, nil
}
