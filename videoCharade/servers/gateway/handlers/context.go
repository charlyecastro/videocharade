package handlers

import (
	"final-project-crew/videoCharade/servers/gateway/indexes"
	"final-project-crew/videoCharade/servers/gateway/models/users"
	"final-project-crew/videoCharade/servers/gateway/sessions"
)

//Context is a struct that stores, the key used to sign and validate SessionIDs,
//the sessions.Store to use when getting or saving session state and
//the users.Store to use when finding or saving user profiles
type Context struct {
	SigningKey   string
	SessionStore sessions.Store
	UserStore    users.Store
	UserTrie     *indexes.Trie
	Notifier     *Notifier
}

//NewContext constructs a new Context ensuring that the dependencies are valid values
// func NewContext(signingKey string, sessionStore *sessions.RedisStore, userStore *users.MySQLStore) *Context {
func NewContext(signingKey string, sessionStore sessions.Store, userStore users.Store, userTrie *indexes.Trie, notifier *Notifier) *Context {
	if len(signingKey) == 0 || sessionStore == nil || userStore == nil {
		panic("can not pass in nil")
	}
	return &Context{signingKey, sessionStore, userStore, userTrie, notifier}
}
