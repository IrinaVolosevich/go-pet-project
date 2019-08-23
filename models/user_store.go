package models

import (
	"database/sql"
	"go-layouts/pkg/mysql"
)

type UserStore interface {
	Find(string) (*User, error)
	FindByEmail(string) (*User, error)
	FindByUsername(string) (*User, error)
	Save(User) error
}

type DBUserStore struct {
	db *sql.DB
}

func (store *DBUserStore) Find(id string) (*User, error) {
	row := store.db.QueryRow(
		"SELECT id, email, password, name, gender "+
			"FROM users WHERE id = ? ORDER BY created_at",
		id)

	user := User{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.HashedPassword,
		&user.Username,
		&user.Gender,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func (store *DBUserStore) FindByEmail(email string) (*User, error) {
	row := store.db.QueryRow(
		"SELECT id, email, password, name, gender "+
			"FROM users WHERE email = ? ORDER BY created_at",
		email)

	user := User{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.HashedPassword,
		&user.Username,
		&user.Gender,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func (store *DBUserStore) FindByUsername(username string) (*User, error) {
	row := store.db.QueryRow(
		"SELECT id, email, password, name, gender "+
			"FROM users WHERE name = ? ORDER BY created_at",
		username)

	user := User{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.HashedPassword,
		&user.Username,
		&user.Gender,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func NewDBUserStore() UserStore {
	return &DBUserStore{
		db: mysql.GlobalMySQLDB,
	}
}

func (store *DBUserStore) Save(user User) error {
	_, err := store.db.Exec(
		"REPLACE INTO users (id, name, password, email, gender) "+
			"VALUES(?, ?, ?, ?, ?)",
			user.ID,
		user.Username,
		user.HashedPassword,
		user.Email,
		*user.Gender)

	return err
}

var GlobalUserStore UserStore

func init() {
	db, err := mysql.NewMySQLDB("dev_main:kr4e65a2xJ9@tcp(localhost:3312)/dev_main")

	if err != nil {
		panic(err)
	}

	mysql.GlobalMySQLDB = db

	GlobalUserStore = NewDBUserStore()
}
