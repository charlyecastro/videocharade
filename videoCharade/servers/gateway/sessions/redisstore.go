package sessions

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

//RedisStore represents a session.Store backed by redis.
type RedisStore struct {
	Client          *redis.Client
	SessionDuration time.Duration
}

//NewRedisStore constructs a new RedisStore
func NewRedisStore(client *redis.Client, sessionDuration time.Duration) *RedisStore {
	return &RedisStore{client, sessionDuration}
}

//Save saves the provided `sessionState` and associated SessionID to the store.
//The `sessionState` parameter is typically a pointer to a struct containing
//all the data you want to associated with the given SessionID.
func (rs *RedisStore) Save(sid SessionID, sessionState interface{}) error {
	data, err := json.Marshal(sessionState)
	if err != nil {
		return err
	}
	rsErr := rs.Client.Set(sid.getRedisKey(), data, rs.SessionDuration).Err()
	if rsErr != nil {
		return rsErr
	}
	return nil
}

//Get populates `sessionState` with the data previously saved
//for the given SessionID
func (rs *RedisStore) Get(sid SessionID, sessionState interface{}) error {

	redisPipe := rs.Client.Pipeline()
	defer redisPipe.Close()
	prevSess := redisPipe.Get(sid.getRedisKey())
	redisPipe.Expire(sid.getRedisKey(), rs.SessionDuration)

	_, err := redisPipe.Exec()
	if err != nil {
		return ErrStateNotFound
	}

	result, err := prevSess.Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(result), sessionState)
}

//Delete deletes all state data associated with the SessionID from the store.
func (rs *RedisStore) Delete(sid SessionID) error {
	rs.Client.Del(sid.getRedisKey())
	return nil
}

//getRedisKey() returns the redis key to use for the SessionID
func (sid SessionID) getRedisKey() string {
	return "sid:" + sid.String()
}
