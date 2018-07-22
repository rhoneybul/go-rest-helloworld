package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gomongo/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"os"
)

type UserController struct {
	session *mgo.Session
}

func NewUserController(s *mgo.Session) *UserController {
	return &UserController{session: s}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	bid := bson.ObjectIdHex(id)

	u := models.User{}

	err := uc.session.DB(os.Getenv("MONGO_DB")).C("users").FindId(bid).One(&u)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	uj, _ := json.Marshal(u)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	w.Write(uj)
}

func (uc UserController) GetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	users := []models.User{}

	err := uc.session.DB(os.Getenv("MONGO_DB")).C("users").Find(bson.M{}).All(&users)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	uj, _ := json.Marshal(users)

	w.Header().Set("Content-Type", "application.json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(uj)

}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id := p.ByName("id")

	bid := bson.ObjectIdHex(id)

	err := uc.session.DB(os.Getenv("MONGO_DB")).C("users").RemoveId(bid)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	u := models.User{}

	json.NewDecoder(r.Body).Decode(&u)

	count, err := uc.session.DB(os.Getenv("MONGO_DB")).C("users").Find(bson.M{"username": u.Username}).Count()

	if err != nil {
		fmt.Println("couldn't execute")
	}

	if count > 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u.Id = bson.NewObjectId()

	uc.session.DB(os.Getenv("MONGO_DB")).C("users").Insert(u)

	uj, _ := json.Marshal(u)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201

	w.Write(uj)

}
