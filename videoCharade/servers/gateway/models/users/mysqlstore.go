package users

import (
	"database/sql"
	"errors"
	"final-project-crew/videoCharade/servers/gateway/indexes"
	"fmt"
	"strings"
	"time"
)

const sqlInsertUser = "insert into users (first_name,last_name,email,user_name,pass_hash,photo_url) values (?,?,?,?,?,?)"
const sqlSelectAll = "select * from users"
const sqlSelectByID = sqlSelectAll + " where id=?"
const sqlSelectByEmail = sqlSelectAll + " where email=?"
const sqlSelectByUserName = sqlSelectAll + " where user_name=?"
const sqlUpdateAll = "update users set first_name =? , last_name =? where id =?"
const sqlUpdateFName = "update users set first_name=? where id=?"
const sqlUpdateLName = "update users set last_name=? where id=?"
const sqlDelete = "delete from users where id=?"
const sqlDeleteAll = "delete from users"

const sqlSignInInsert = "insert into sign_ins (user_id, date_time, ip_address) values (?,?,?)"
const sqlSignInGetByID = "select * from sign_ins where id=?"
const sqlSignInGetByUserID = "select * from sign_ins where user_id=?"
const sqlSignInGetAll = "select * from sign_ins"
const sqlSignInDelete = "delete from sing_ins where id=?"

var InvalidUpdate = errors.New("invalid update")

//MySQLStore represents a tasks.Store backed by MySQL
type MySQLStore struct {
	//db is the open database object
	//this store will use to send queries
	//to the database
	db *sql.DB
}

type SignInLog struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userId"`
	DateTime  time.Time `json:"dateTime"`
	IPAddress string    `json:"ipAddress"`
}

//NewMySQLStore constructs a new MySQLStore
func NewMySQLStore(db *sql.DB) *MySQLStore {
	if db == nil {
		panic("nil database pointer")
	}
	return &MySQLStore{db}
}

//GetByID takes in a Id and returns a user & error.
func (s *MySQLStore) GetByID(id int64) (*User, error) {

	row := s.db.QueryRow(sqlSelectByID, id)
	user, err := scanUser(row)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//GetByEmail takes in an email string and returns a user & error
func (s *MySQLStore) GetByEmail(email string) (*User, error) {
	row := s.db.QueryRow(sqlSelectByEmail, email)
	user, err := scanUser(row)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//GetByUserName takes in a username string and returns a user & error
func (s *MySQLStore) GetByUserName(username string) (*User, error) {
	row := s.db.QueryRow(sqlSelectByUserName, username)
	user, err := scanUser(row)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//Insert takes in a User struct as a parameter which it will populate the user table with
func (s *MySQLStore) Insert(user *User) (*User, error) {

	results, err := s.db.Exec(sqlInsertUser, user.FirstName, user.LastName, user.Email, user.UserName, user.PassHash, user.PhotoURL)
	if err != nil {
		return nil, fmt.Errorf("executing insert: %v", err)
	}
	id, err := results.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("getting new ID: %v", err)
	}

	user.ID = id
	return user, nil
}

//Update takes in a id and updates struct, it finds the row with the given id and apply the updates in
//the given updates parameter
func (s *MySQLStore) Update(id int64, updates *Updates) (*User, error) {
	var results sql.Result
	var err error
	firstName := updates.FirstName
	lastName := updates.LastName
	if len(firstName) < 1 && len(lastName) < 1 {
		return nil, InvalidUpdate
	} else if len(firstName) > 1 && len(lastName) > 1 {
		results, err = s.db.Exec(sqlUpdateAll, firstName, lastName, id)
	} else if len(firstName) > 1 {
		results, err = s.db.Exec(sqlUpdateFName, firstName, id)
	} else if len(lastName) > 1 {
		results, err = s.db.Exec(sqlUpdateLName, lastName, id)
	}

	if err != nil {
		return nil, fmt.Errorf("updating: %v", err)
	}
	affected, err := results.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("getting rows affected: %v", err)
	}
	if affected == 0 {
		return nil, ErrUserNotFound
	}
	return s.GetByID(id)
}

//Delete takes an id and return an error. Finds a row to delete from the database with the
// given id
func (s *MySQLStore) Delete(id int64) error {
	_, err := s.db.Exec(sqlDelete, id)
	if err != nil {
		return fmt.Errorf("Error executing DELETE operation")
	}
	return nil
}

//SignInInsert inserts a new sign in into the signin table
func (s *MySQLStore) SignInInsert(log *SignInLog) (*SignInLog, error) {
	results, err := s.db.Exec(sqlSignInInsert, log.UserID, log.DateTime, log.IPAddress)
	if err != nil {
		return nil, fmt.Errorf("Error executing INSERT operation")
	}
	id, err := results.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("getting new ID: %v", err)
	}

	log.ID = id
	return log, nil
}

//GetAllSignInByUserID takes in a Id and returns a user & error.
func (s *MySQLStore) GetAllSignInByUserID(userID int64) ([]*SignInLog, error) {
	rows, err := s.db.Query(sqlSignInGetByUserID, userID)
	if err != nil {
		return nil, fmt.Errorf("selecting: %v", err)
	}
	var list []*SignInLog
	for rows.Next() {
		log := &SignInLog{}
		rows.Scan(&log.ID, &log.UserID, &log.DateTime, &log.IPAddress)
		list = append(list, log)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over rows: %v", err)
	}
	return list, nil
}

func scanUser(row *sql.Row) (*User, error) {
	user := &User{}
	if err := row.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.UserName, &user.PassHash, &user.PhotoURL); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("scanning: %v", err)
	}
	return user, nil
}

//LoadAllUsers iterates through all users and populates the trie with each users first name, last name and username
func (s *MySQLStore) LoadAllUsers() (*indexes.Trie, error) {
	time.Sleep(20 * time.Second)
	t := indexes.NewTrie()
	rows, err := s.db.Query(sqlSelectAll)
	if err != nil {
		return nil, fmt.Errorf("selecting: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		user := &User{}
		rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.UserName, &user.PassHash, &user.PhotoURL)

		t.Add(strings.ToLower(user.FirstName), user.ID)
		t.Add(strings.ToLower(user.LastName), user.ID)
		t.Add(strings.ToLower(user.UserName), user.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over rows: %v", err)
	}
	return t, nil
}
