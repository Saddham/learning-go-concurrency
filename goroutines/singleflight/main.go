package main

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/singleflight"
)

var users = map[string][]byte{}
var sfGroup singleflight.Group

func main() {
	users["saddham"], _ = bcrypt.GenerateFromPassword([]byte("foo"), 10)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		un, pw, _ := r.BasicAuth()

		key := un + "." + pw
		_, err, _ := sfGroup.Do(key, func() (interface{}, error) {
			if err := bcrypt.CompareHashAndPassword(users[un], []byte(pw)); err != nil {
				return nil, err
			}

			return nil, nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not allowed in"))

			return
		}

		w.Write([]byte("Success"))
	})

	log.Fatal(http.ListenAndServe(":3001", nil))
}

// curl http://localhost:3001/login
// curl http://saddham:wrong@localhost:3001/login
// curl http://saddham:foo@localhost:3001/login
