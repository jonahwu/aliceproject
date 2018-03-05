package main

import (
	"errors"
	"fmt"
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetTransID(c *mgo.Collection) (string, error) {
	u1 := uuid.Must(uuid.NewV4()).String()
	fmt.Printf("UUIDv4: %s\n", u1)
	num, _ := c.Find(bson.M{"id": u1}).Count()
	fmt.Println(num)
	if num == 0 {
		return u1, nil
	} else {
		return "", errors.New("found repeated uuid, try again")
	}
}

func DBInsertImage(c *mgo.Collection, ctx iris.Context) (string, error) {
	var img ImageInfo
	img.Name = ctx.GetHeader("IMG")
	img.ID = ctx.Params().Get("id")
	//err := c.Insert(&ImageInfo{"IMG", "+55 53 8116 9639", id})
	if err := c.Insert(&img); err != nil {
		fmt.Println(err)
		return "", err
	}
	return img.Name, nil

}

func DBUserInfoStatusDone(c *mgo.Collection, ctx iris.Context) (string, error) {
	id := ctx.Params().Get("id")
	updateTarget := bson.M{"id": id}
	change := bson.M{"$set": bson.M{"status": "done"}}
	if err := c.Update(updateTarget, change); err != nil {
		fmt.Println("update error:", err)
		return "", err
	}
	return id, nil
}
func DBInsertUserInfo(c *mgo.Collection, ctx iris.Context) (string, error) {

	id := ctx.Params().Get("id")
	fmt.Println(id)

	var userRead User
	ctx.ReadJSON(&userRead)
	fmt.Println(userRead)
	fmt.Println("where the input user name is", userRead.Username)
	userRead.Status = "processing"
	userRead.ID = id
	if err := c.Insert(&userRead); err != nil {
		fmt.Println(err)
		return "", err
	}
	return id, nil
}
