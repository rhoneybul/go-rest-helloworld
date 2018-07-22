package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	// "io"
	"encoding/json"
	"gomongo/controllers"
	"gopkg.in/mgo.v2"
	"net/http"
	"os"
)

var HOST = fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	greeting := map[string]string{"Welcome": "Hello, World"}
	data, _ := json.Marshal(greeting)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func getSession() *mgo.Session {
	s, err := mgo.Dial(os.Getenv("MONGO_URI"))
	if err != nil {
		panic(err)
	}
	return s
}

func main() {

	mongoSession := getSession()

	router := httprouter.New()

	router.GET("/", index)

	// users
	uc := controllers.NewUserController(mongoSession)
	router.POST("/user", uc.CreateUser)
	router.GET("/user/:id", uc.GetUser)
	router.DELETE("/user/:id", uc.DeleteUser)
	router.GET("/user", uc.GetUsers)

	fmt.Printf("Server Running on %s\n", HOST)

	err := http.ListenAndServe(HOST, router)

	if err != nil {
		fmt.Println(err)
	}
}
