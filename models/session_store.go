package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type SessionStore interface {
	Find(string) (*Session, error)
	Save(*Session) error
	Delete(*Session) error
}

type FileSessionStore struct {
	filename string
	Sessions map[string]Session
}

func NewFileSessionStore(name string) (*FileSessionStore, error) {
	store := &FileSessionStore{
		filename: name,
		Sessions: map[string]Session{},
	}

	contents, err := ioutil.ReadFile(name)

	if err != nil {
		if os.IsNotExist(err) {
			return store, err
		}

		return nil, err
	}

	err = json.Unmarshal(contents, store)

	if err != nil {
		return nil, err
	}

	return store, nil
}

func (store *FileSessionStore) Find(id string) (*Session, error){
	session, exists := store.Sessions[id]
	if !exists {
		return nil, nil
	}

	return &session, nil
}

func (store *FileSessionStore) Save(session *Session) error{
	store.Sessions[session.ID] = *session

	contents, err := json.MarshalIndent(store, "", " ")

	if err != nil {
		return  err
	}

	return ioutil.WriteFile(store.filename, contents, 0660)
}

func (store *FileSessionStore) Delete(session *Session) error {
	delete(store.Sessions, session.ID)

	contents, err := json.MarshalIndent(store, "", " ")

	if err != nil {
		return  err
	}

	return ioutil.WriteFile(store.filename, contents, 0660)
}

var GlobalSessionStore SessionStore

func init() {
	store, err := NewFileSessionStore("./../assets/sessions.json")

	if err != nil {
		panic(fmt.Errorf("Error creating session store: %s", err))
	}

	GlobalSessionStore = store
}