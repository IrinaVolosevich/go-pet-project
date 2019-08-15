package models

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type User struct {
	ID             string
	Email          string
	HashedPassword string
	Username       string
}

const (
	passwordLength = 8
	hashCost       = 10
	userIDLength
)

func NewUser(username, email, password string) (User, error) {
	user := User{
		Email:    email,
		Username: username,
	}

	if username == "" {
		return user, ErrNoUsername
	}

	if email == "" {
		return user, ErrNoEmail
	}

	if password == "" {
		return user, ErrNoPassword
	}

	if len(password) < passwordLength {
		return user, ErrPasswordTooShort
	}

	existingUser, err := UserStore.FindByUsername(GlobalUserStore, username)

	if err != nil {
		return user, err
	}

	if existingUser != nil {
		return user, ErrUsernameExists
	}

	existingUser, err = UserStore.FindByEmail(GlobalUserStore, email)

	if err != nil {
		return user, err
	}

	if existingUser != nil {
		return user, ErrEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)

	user.HashedPassword = string(hashedPassword)
	user.ID = GenerateID("usr", userIDLength)

	err = UserStore.Save(GlobalUserStore, user)

	if err != nil {
		panic(err)
	}

	return user, err
}

func FindUser(username, password string) (*User, error) {
	out := &User{
		Username: username,
	}

	existingUser, error := UserStore.FindByUsername(GlobalUserStore, username)

	if error != nil {
		return out, error
	}

	if existingUser == nil {
		return nil, ErrCredentialsIncorrect
	}

	if bcrypt.CompareHashAndPassword([]byte(existingUser.HashedPassword), []byte(password)) != nil {
		return out, ErrCredentialsIncorrect
	}

	return existingUser, nil
}

func FindOrCreateSession(w http.ResponseWriter, r *http.Request) *Session {
	session := RequestSession(r)

	if session == nil {
		session = NewSession(w)
	}

	return session
}

func UpdateUser(user *User, email, currentPassword, newPassword string) (User, error) {
	out := *user
	out.Email = email

	existingUser,err := UserStore.FindByUsername(GlobalUserStore, email)

	if err != nil {
		return out, err
	}

	if existingUser != nil && existingUser.ID != user.ID {
		return out, err
	}

	user.Email = email

	if currentPassword == "" {
		return out, nil
	}

	if bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(currentPassword)) != nil {
		return out, ErrPasswordIncorrect
	}

	if newPassword == "" {
		return out, ErrNoPassword
	}

	if len(newPassword) < passwordLength {
		return out, ErrPasswordTooShort
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), hashCost)
	user.HashedPassword = string(hashedPassword)

	return out, err
}