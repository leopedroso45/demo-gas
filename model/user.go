package model

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Id       string `json:"_id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

func (n User) ToJSON() []byte {
	data, err := json.Marshal(n)

	if err != nil {
		fmt.Println("Serialization of Json failed")
		return nil
	}

	return data
}

func (n *User) ToBSON() []byte {
	data, err := bson.Marshal(n)

	if err != nil {
		fmt.Println("Serialization of Bson failed")
		return nil
	}

	return data
}

func (n *User) ClearUserDetails() {
	n.Password = ""
}

func (n *User) PassToSha1() {
	n.Password = toSha1(n.Password)
}

func toSha1(value string) string {
	h := sha1.New()
	h.Write([]byte(value))
	bs := h.Sum(nil)
	return string(bs)
}
