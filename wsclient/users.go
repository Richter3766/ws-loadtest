package wsclient

import (
	"encoding/json"
	"os"
)

func LoadUsers(filename string) ([]User, error) {
    f, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    var users []User
    err = json.NewDecoder(f).Decode(&users)
    return users, err
}