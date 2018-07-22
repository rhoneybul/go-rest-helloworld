package models

import "gopkg.in/mgo.v2/bson"

type User struct {
	Name     string        `json: "name"`
	Username string        `json: "username"`
	Email    string        `json: "email"`
	Password string        `json: "password"`
	Id       bson.ObjectId `json: "id" bson: "_id"`
}
