package handlers

import (
	"final-project-crew/videoCharade/servers/gateway/models/users"
	"time"
)

//TODO: define a session state struct for this web server
//see the assignment description for the fields you should include
//remember that other packages can only see exported fields!

//SessionState is a struct that tracks the Time at which this session began and
//The authenticated User who started the session
type SessionState struct {
	SessionBegin time.Time   `json:"sessionBegin,omitempty"`
	User         *users.User `json:"user,omitempty"`
}
