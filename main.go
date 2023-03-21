package main

import (
	"net/http"
	"package30/lib30"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", lib30.Hello)
	r.Post("/create", lib30.CreateUser)
	r.Post("/make_friends", lib30.MakeFriends)
	r.Delete("/user/{id}", lib30.DeleteUser)
	r.Get("/friends/{id}", lib30.GetUserFriends)
	r.Put("/{id}", lib30.UpdateUserAge)

	http.ListenAndServe(":8080", r)

}
