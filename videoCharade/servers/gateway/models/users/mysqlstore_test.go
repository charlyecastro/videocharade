package users

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestMySQLStore_GetByID(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("error creating sql mock: %v", err)
	}

	defer db.Close()

	expectedUser := &User{
		ID:        1,
		Email:     "johndoe@email.com",
		FirstName: "John",
		LastName:  "Doe",
		UserName:  "johndoe123",
		PassHash:  []byte("password123"),
		PhotoURL:  "https://www.gravatar.com/avatar/sakjfddslkjfei",
	}

	// Initialize a MySQLStore struct to allow us to interface with the SQL client
	store := NewMySQLStore(db)

	// Create a row with the appropriate fields in your SQL database
	// Add the actual values to the row
	row := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "user_name", "pass_hash", "photo_url"})
	row.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.FirstName, expectedUser.LastName, expectedUser.UserName, expectedUser.PassHash, expectedUser.PhotoURL)

	// Expecting a successful "query"
	// This tells our db to expect this query (id) as well as supply a certain response (row)

	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectByID)).
		WithArgs(expectedUser.ID).WillReturnRows(row)

	// Since we know our query is successful, we want to test whether there happens to be
	// any expected error that may occur.
	user, err := store.GetByID(expectedUser.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Again, since we are assuming that our query is successful, we can test for when our
	// function doesn't work as expected.
	if err == nil && !reflect.DeepEqual(user, expectedUser) {
		t.Errorf("User Result: %v", user)
		t.Errorf("Expected Result: %v", expectedUser)
		t.Errorf("User queried does not match expected user")
	}

	// Expecting a unsuccessful "query"
	// Attempting to search by an id that doesn't exist. This would result in a
	// sql.ErrNoRows error
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectByID)).
		WithArgs(-1).WillReturnError(sql.ErrNoRows)

	// Since we are expecting an error here, we create a condition opposing that to see
	// if our GetById is working as expected
	if _, err = store.GetByID(-1); err == nil {
		t.Errorf("Expected error: %v, but recieved nil", sql.ErrNoRows)
	}

	// Attempting to trigger a DBMS querying error
	queryingErr := fmt.Errorf("DBMS error when querying")
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectByID)).
		WithArgs(expectedUser.ID).WillReturnError(queryingErr)

	if _, err = store.GetByID(expectedUser.ID); err == nil {
		t.Errorf("Expected error: %v, but recieved nil", queryingErr)
	}

	// This attempts to check if there are any expectations that we haven't met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet sqlmock expectations: %v", err)
	}

}

func TestMySQLStore_GetByEmail(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("error creating sql mock: %v", err)
	}

	defer db.Close()

	expectedUser := &User{
		ID:        1,
		Email:     "johndoe@email.com",
		FirstName: "John",
		LastName:  "Doe",
		UserName:  "johndoe123",
		PassHash:  []byte("password123"),
		PhotoURL:  "https://www.gravatar.com/avatar/sakjfddslkjfei",
	}

	// Initialize a MySQLStore struct to allow us to interface with the SQL client
	store := NewMySQLStore(db)

	// Create a row with the appropriate fields in your SQL database
	// Add the actual values to the row
	row := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "user_name", "pass_hash", "photo_url"})
	row.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.FirstName, expectedUser.LastName, expectedUser.UserName, expectedUser.PassHash, expectedUser.PhotoURL)

	// Expecting a successful "query"
	// This tells our db to expect this query (id) as well as supply a certain response (row)

	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectByEmail)).
		WithArgs(expectedUser.Email).WillReturnRows(row)

	// Since we know our query is successful, we want to test whether there happens to be
	// any expected error that may occur.
	user, err := store.GetByEmail(expectedUser.Email)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Again, since we are assuming that our query is successful, we can test for when our
	// function doesn't work as expected.
	if err == nil && !reflect.DeepEqual(user, expectedUser) {
		t.Errorf("User Result: %v", user)
		t.Errorf("Expected Result: %v", expectedUser)
		t.Errorf("User queried does not match expected user")
	}

	// Expecting a unsuccessful "query"
	// Attempting to search by an id that doesn't exist. This would result in a
	// sql.ErrNoRows error
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectByEmail)).
		WithArgs("wrong@email.com").WillReturnError(sql.ErrNoRows)

	// Since we are expecting an error here, we create a condition opposing that to see
	// if our GetById is working as expected
	if _, err = store.GetByEmail("wrong@email.com"); err == nil {
		t.Errorf("Expected error: %v, but recieved nil", sql.ErrNoRows)
	}

	// Attempting to trigger a DBMS querying error
	queryingErr := fmt.Errorf("DBMS error when querying")
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectByEmail)).
		WithArgs(expectedUser.Email).WillReturnError(queryingErr)

	if _, err = store.GetByEmail(expectedUser.Email); err == nil {
		t.Errorf("Expected error: %v, but recieved nil", queryingErr)
	}

	// This attempts to check if there are any expectations that we haven't met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet sqlmock expectations: %v", err)
	}
}

func TestMySQLStore_GetByUserName(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("error creating sql mock: %v", err)
	}

	defer db.Close()

	expectedUser := &User{
		ID:        1,
		Email:     "johndoe@email.com",
		FirstName: "John",
		LastName:  "Doe",
		UserName:  "johndoe123",
		PassHash:  []byte("password123"),
		PhotoURL:  "https://www.gravatar.com/avatar/sakjfddslkjfei",
	}

	store := NewMySQLStore(db)

	row := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "user_name", "pass_hash", "photo_url"})
	row.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.FirstName, expectedUser.LastName, expectedUser.UserName, expectedUser.PassHash, expectedUser.PhotoURL)

	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectByUserName)).
		WithArgs(expectedUser.UserName).WillReturnRows(row)

	user, err := store.GetByUserName(expectedUser.UserName)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if err == nil && !reflect.DeepEqual(user, expectedUser) {
		t.Errorf("User Result: %v", user)
		t.Errorf("Expected Result: %v", expectedUser)
		t.Errorf("User queried does not match expected user")
	}

	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectByUserName)).
		WithArgs("johndoe321").WillReturnError(sql.ErrNoRows)

	if _, err = store.GetByUserName("johndoe321"); err == nil {
		t.Errorf("Expected error: %v, but recieved nil", sql.ErrNoRows)
	}

	queryingErr := fmt.Errorf("DBMS error when querying")
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectByUserName)).
		WithArgs(expectedUser.UserName).WillReturnError(queryingErr)

	if _, err = store.GetByUserName(expectedUser.UserName); err == nil {
		t.Errorf("Expected error: %v, but recieved nil", queryingErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet sqlmock expectations: %v", err)
	}

}

func TestMySQLStore_Insert(t *testing.T) {
	//create a new sql mock
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("error creating sql mock: %v", err)
	}
	defer db.Close()

	inputUser := &User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@email.com",
		UserName:  "johndoe123",
		PassHash:  []byte("password123"),
		PhotoURL:  "https://www.gravatar.com/avatar/sakjfddslkjfei",
	}

	// Initialize a MySQLStore struct to allow us to interface with the SQL client
	store := NewMySQLStore(db)

	//insert an invalid user, should create a new row and should increment the lastinsertid and rows affected
	mock.ExpectExec(regexp.QuoteMeta(sqlInsertUser)).
		WithArgs(inputUser.FirstName, inputUser.LastName, inputUser.Email, inputUser.UserName, inputUser.PassHash, inputUser.PhotoURL).
		WillReturnResult(sqlmock.NewResult(2, 1))

	user, err := store.Insert(inputUser)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if err == nil && !reflect.DeepEqual(user, inputUser) {
		t.Errorf("User returned does not match input user")
	}

	// Inserting an invalid user should not update the table and should return an error
	invalidUser := &User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@email.com",
		UserName:  "johndoe123",
		PassHash:  []byte("password123"),
		PhotoURL:  "https://www.gravatar.com/avatar/sakjfddslkjfei",
	}
	insertErr := fmt.Errorf("Error executing INSERT operation")
	mock.ExpectExec(regexp.QuoteMeta(sqlInsertUser)).
		WithArgs(inputUser.FirstName, inputUser.LastName, inputUser.Email, inputUser.UserName, inputUser.PassHash, inputUser.PhotoURL).
		WillReturnError(insertErr)

	if _, err = store.Insert(invalidUser); err == nil {
		t.Errorf("Expected error: %v but recieved nil", insertErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet sqlmock expectations: %v", err)
	}
}

func TestMySQLStore_Update(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("error creating sql mock: %v", err)
	}

	defer db.Close()

	defaultUser := &User{
		ID:        1,
		Email:     "johndoe@email.com",
		FirstName: "John",
		LastName:  "Doe",
		UserName:  "johndoe123",
		PassHash:  []byte("password123"),
		PhotoURL:  "https://www.gravatar.com/avatar/sakjfddslkjfei",
	}

	// Initialize a MySQLStore struct to allow us to interface with the SQL client
	store := NewMySQLStore(db)

	row := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "user_name", "pass_hash", "photo_url"})
	row.AddRow(defaultUser.ID, defaultUser.Email, defaultUser.FirstName, defaultUser.LastName, defaultUser.UserName, defaultUser.PassHash, defaultUser.PhotoURL)

	//passing in normal update which has a first name and last name, should update the users first and last name
	updatesAll := &Updates{
		FirstName: "charlye",
		LastName:  "castro",
	}

	expectedUserAll := &User{
		ID:        1,
		Email:     "johndoe@email.com",
		FirstName: "charlye",
		LastName:  "castro",
		UserName:  "johndoe123",
		PassHash:  []byte("password123"),
		PhotoURL:  "https://www.gravatar.com/avatar/sakjfddslkjfei",
	}

	mock.ExpectExec(regexp.QuoteMeta(sqlUpdateAll)).
		WithArgs(updatesAll.FirstName, updatesAll.LastName, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	rowAllExp := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "user_name", "pass_hash", "photo_url"})
	rowAllExp.AddRow(defaultUser.ID, defaultUser.Email, expectedUserAll.FirstName, expectedUserAll.LastName, defaultUser.UserName, defaultUser.PassHash, defaultUser.PhotoURL)
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectByID)).WithArgs(1).WillReturnRows(rowAllExp)

	userAll, err := store.Update(1, updatesAll)
	if err != nil {
		t.Errorf("error is right here")
		t.Errorf("Unexpected error: %v", err)
	}

	if err == nil && !reflect.DeepEqual(userAll, expectedUserAll) {
		t.Errorf("User Result: %v", userAll)
		t.Errorf("Rows: %v", row)
		t.Errorf("Expected Result: %v", expectedUserAll)
		t.Errorf("User queried does not match expected user")
	}
	//passing update with only first name, should only update the users first name and leave the last name alone
	updatesFName := &Updates{
		FirstName: "charlye",
		LastName:  "",
	}

	expectedUserFName := &User{
		ID:        1,
		Email:     "johndoe@email.com",
		FirstName: "charlye",
		LastName:  "Doe",
		UserName:  "johndoe123",
		PassHash:  []byte("password123"),
		PhotoURL:  "https://www.gravatar.com/avatar/sakjfddslkjfei",
	}

	mock.ExpectExec(regexp.QuoteMeta(sqlUpdateFName)).
		WithArgs(updatesFName.FirstName, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	rowFNameExp := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "user_name", "pass_hash", "photo_url"})
	rowFNameExp.AddRow(defaultUser.ID, defaultUser.Email, expectedUserFName.FirstName, expectedUserFName.LastName, defaultUser.UserName, defaultUser.PassHash, defaultUser.PhotoURL)
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectByID)).WithArgs(1).WillReturnRows(rowFNameExp)

	userFName, err := store.Update(1, updatesFName)
	if err != nil {
		t.Errorf("error is right here")
		t.Errorf("Unexpected error: %v", err)
	}

	if err == nil && !reflect.DeepEqual(userFName, expectedUserFName) {
		t.Errorf("User Result: %v", userFName)
		t.Errorf("Rows: %v", row)
		t.Errorf("Expected Result: %v", expectedUserFName)
		t.Errorf("User queried does not match expected user")
	}

	//passing update with only lastname, should only update the users last name and leave the first name alone
	updatesLName := &Updates{"", "castro"}

	expectedUserLName := &User{
		ID:        1,
		Email:     "johndoe@email.com",
		FirstName: "John",
		LastName:  "castro",
		UserName:  "johndoe123",
		PassHash:  []byte("password123"),
		PhotoURL:  "https://www.gravatar.com/avatar/sakjfddslkjfei",
	}

	mock.ExpectExec(regexp.QuoteMeta(sqlUpdateLName)).
		WithArgs(updatesLName.LastName, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	rowLNameExp := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "user_name", "pass_hash", "photo_url"})
	rowLNameExp.AddRow(defaultUser.ID, defaultUser.Email, expectedUserLName.FirstName, expectedUserLName.LastName, defaultUser.UserName, defaultUser.PassHash, defaultUser.PhotoURL)
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectByID)).WithArgs(1).WillReturnRows(rowLNameExp)

	userLName, err := store.Update(1, updatesLName)
	if err != nil {
		t.Errorf("error is right here")
		t.Errorf("Unexpected error: %v", err)
	}

	if err == nil && !reflect.DeepEqual(userLName, expectedUserLName) {
		t.Errorf("User Result: %v", userLName)
		t.Errorf("Rows: %v", row)
		t.Errorf("Expected Result: %v", expectedUserLName)
		t.Errorf("User queried does not match expected user")
	}

	//passing in empty updates, which should return an error
	// updatesEmpty := &Updates{
	// 	FirstName: "",
	// 	LastName:  "",
	// }

	// noEmptyFieldsErr := fmt.Errorf("Cannot have both fields be empty")
	// mock.ExpectExec(regexp.QuoteMeta(sqlUpdateAll)).
	// 	WithArgs(updatesEmpty.FirstName, updatesEmpty.LastName, 1).
	// 	WillReturnError(noEmptyFieldsErr)

	// _, emptyErr := store.Update(1, updatesEmpty)

	// if emptyErr == nil {
	// 	t.Errorf("Expected error: %v, but recieved nil", noEmptyFieldsErr)
	// }

	// Expecting a unsuccessful "query"
	updateErr := fmt.Errorf("Error executing Update operation")
	mock.ExpectExec(regexp.QuoteMeta(sqlUpdateAll)).
		WithArgs(updatesAll.FirstName, updatesAll.LastName, 5).
		WillReturnError(updateErr)

	if _, err = store.Update(5, updatesAll); err == nil {
		t.Errorf("Expected error: %v but recieved nil", updateErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet sqlmock expectations: %v", err)
	}

}

func TestMySQLStore_Delete(t *testing.T) {
	//create a new sql mock
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("error creating sql mock: %v", err)
	}
	//ensure it's closed at the end of the test
	defer db.Close()

	// Initialize a user struct we will use as a test variable
	user := &User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@email.com",
		UserName:  "johndoe123",
		PassHash:  []byte("password123"),
		PhotoURL:  "https://www.gravatar.com/avatar/sakjfddslkjfei",
	}

	// Initialize a MySQLStore struct to allow us to interface with the SQL client
	store := NewMySQLStore(db)

	row := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "user_name", "pass_hash", "photo_url"})
	row.AddRow(user.ID, user.FirstName, user.LastName, user.Email, user.UserName, user.PassHash, user.PhotoURL)

	// This tells our db to expect an insert query with certain arguments with a certain
	// return result
	mock.ExpectExec(regexp.QuoteMeta(sqlDelete)).
		WithArgs(user.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = store.Delete(user.ID)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	deleteErr := fmt.Errorf("Error executing DELETE operation")
	mock.ExpectExec(regexp.QuoteMeta(sqlDelete)).
		WithArgs(5).
		WillReturnError(deleteErr)

	if err = store.Delete(5); err == nil {
		t.Errorf("Expected error: %v but recieved nil", deleteErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet sqlmock expectations: %v", err)
	}
}
