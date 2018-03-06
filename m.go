package main

import (
	"fmt"
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

func TestController(t int) string {
	fmt.Println("into controller", t)
	return "haha"
}

type User struct {
	Username   string `json:"username" bson:"username"`
	City       string `json:"city" bson:"city"`
	ID         string `json:"id" bson:"id"`
	Status     string `json:"status" bson:"status"`
	ClientTime string `json:"clienttime" bson:"clienttime"`
	Timestamp  string `json:"timestamp" bson:"timestamp"`
}
type History struct {
	Users []User `json:"users" bson:"users"`
}
type TransID struct {
	ID string `json:"id" bson:"id"`
}

type ImageInfo struct {
	Name string `json:"name" bson:"name"`
	ID   string `json:"id" bson:"id"`
}

//curl -H "TEST:B" http://localhost:8080/test/2X3190
func middleware(ctx iris.Context) {
	shareInformation := "this is a sharable information between handlers"

	requestPath := ctx.Path()
	println("Before the mainHandler: " + requestPath)

	ctx.Values().Set("info", shareInformation)
	checkHeader := ctx.GetHeader("TEST")
	if checkHeader == "" {
		return
	}
	fmt.Println(checkHeader)
	ctx.Next() // execute the next handler, in this case the main one.
}

func main() {

	//testinput := 2
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c1 := session.DB("test").C("people")
	cimage := session.DB("test").C("image")
	chistory := session.DB("test").C("history")
	//c2 := session.DB("test").C("history")

	app := iris.New()
	// Load all templates from the "./views" folder
	// where extension is ".html" and parse them
	// using the standard `html/template` package.
	app.RegisterView(iris.HTML("./views", ".html"))

	// Method:    GET
	// Resource:  http://localhost:8080
	app.Get("/", func(ctx iris.Context) {
		// Bind: {{.message}} with "Hello world!"
		ctx.ViewData("message", "Hello world!")
		// Render template file: ./views/hello.html
		ctx.View("hello.html")
	})
	//curl -H "Content-Type: application/json" -H "TEST:B" -d '{"Firstname":"aaa", "Lastname":"bbb", "Username":"ccc","City":"ddd"}' http://localhost:8080/test/2X3190 -v
	//curl -X GET -H "Content-Type: application/json" -H "TEST:B" -d '{"Firstname":"aaa", "Lastname":"bbb", "Username":"ccc","City":"ddd"}' http://localhost:8080/getid -v
	app.Get("/getid", func(ctx iris.Context) {
		// Bind: {{.message}} with "Hello world!"

		transid, err := GetTransID(c1)
		if err != nil {
			fmt.Println("get tranlation id:", transid)
		}
		//it's ok
		var tid TransID
		tid.ID = transid
		//it's ok
		/*
			tid := make(map[string]interface{})
			tid["ID"] = transid
		*/
		//it's ok
		/*
			tid := TransID{
				Id: "Johndoe",
			}
		*/
		fmt.Println(tid)
		ctx.Header("TEST", "RETURN")
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(tid)
	})

	//curl -H "Content-Type: application/json" -H "TEST:B" -d '{"Firstname":"aaa", "Lastname":"bbb", "Username":"ccc","City":"ddd"}' http://localhost:8080/test/2X3190 -v
	//curl -X POST -H "Content-Type: application/json" -H "TEST:B" -d '{"Username":"cc","City":"ddd"}' http://localhost:8080/userinfo/e6afa881-ab96-46cc-a2bf-83bde68f010f -v
	//curl -X POST -H "Content-Type: application/json" -H "TEST:B" -d '{"Username":"山小","City":"ddd"}' http://localhost:8080/userinfo/59036923-79fc-4fc2-8b4d-3e39a6b7624a -v
	app.Post("/userinfo/{id}", middleware, func(ctx iris.Context) {
		DBInsertUserInfo(c1, ctx)
		ctx.Header("TEST", "RETURN")
		ctx.StatusCode(iris.StatusOK)
	})

	//curl -H "Content-Type: application/json" -H "TEST:B" -d @a.json http://localhost:8080/test/2X3190 -v
	//curl -X PUT -H "Content-Type: application/json" -H "TEST:B" -d @a.json http://localhost:8080/userinfo/e6afa881-ab96-46cc-a2bf-83bde68f010f -v
	//curl -X PUT -H "Content-Type: application/json" -H "TEST:B" -d @a.json http://localhost:8080/userinfo/59036923-79fc-4fc2-8b4d-3e39a6b7624a
	app.Put("/userinfo/{id}", middleware, func(ctx iris.Context) {
		if _, err := DBUserInfoStatusDone(c1, ctx); err != nil {
			fmt.Println(err)
			ctx.Writef("some error")
			ctx.StatusCode(iris.StatusExpectationFailed)
		}
		ctx.StatusCode(iris.StatusOK)
	})

	app.Get("/userinfo/{id}", middleware, func(ctx iris.Context) {
		//id, _ := ctx.Params().Get("id")
		id := ctx.Params().Get("id")
		fmt.Println(id)

		var userRead User
		ctx.ReadJSON(&userRead)
		fmt.Println(userRead)

		//ret := MongoController(c1, id)
		ret, _ := DBGetUserInfo(c1, ctx)
		fmt.Println(ret)

		//it's a Header reuturn
		ctx.Header("TEST", "RETURN")
		ctx.JSON(ret)

	})
	app.Delete("/userinfo/{id}", middleware, func(ctx iris.Context) {
		DBMoveUserInfo(c1, chistory, ctx)

	})
	//curl -X POST -F 'img_avatar=@/go/src/github.com/aliciproject/mongocmd'  -H "TEST:B"  http://localhost:8080/image -vi
	app.Post("/image/{id}", middleware, func(ctx iris.Context) {
		// it will keep the filename mongocmd so the dest is just a directory that we must create by myself
		urlPath := ctx.Params().Get("id")
		fmt.Println("show the path:", urlPath)
		dest := "/opt/data/aliceproject/image/" + urlPath
		cmdStr := "mkdir -p " + dest
		_ = RunCommand(cmdStr)
		fmt.Println("show dest:", dest)
		n, err := ctx.UploadFormFiles(dest)
		fmt.Println("show read file number", n)
		fmt.Println("show error", err)
		//ctx.Header("Content-Type", "multipart/form-data")

	})
	// this post can be change to put in image
	//curl -X POST  -H "TEST:B" -H "IMG:mongocmd"  http://localhost:8080/imageinfo/e6afa881-ab96-46cc-a2bf-83bde68f010f -v
	app.Post("/imageinfo/{id}", middleware, func(ctx iris.Context) {
		_, err := DBInsertImage(cimage, ctx)
		if err != nil {
			ctx.Writef("error")
		}
		ctx.Writef("ok")
	})

	//curl  -H "TEST:B"  http://localhost:8080/image -vi
	//where destFile stored in header and can be a hint name for client, here we use "default"
	//curl  -H "TEST:B"  http://localhost:8080/image/e6afa881-ab96-46cc-a2bf-83bde68f010f -v
	//curl  -H "TEST:B" -H "IMG:mongocmd"  http://localhost:8080/image/e6afa881-ab96-46cc-a2bf-83bde68f010f -v
	app.Get("/image/{id}", middleware, func(ctx iris.Context) {
		urlPath := ctx.Params().Get("id")
		imgName := ctx.GetHeader("IMG")
		dest := "/opt/data/aliceproject/image/" + urlPath
		file := dest + "/" + imgName
		fmt.Println("get file: ", file)
		ctx.SendFile(file, "default")
	})

	//curl  -H "Content-Type: application/json" -H "TEST:B"  http://localhost:8080/history
	app.Get("/history", middleware, func(ctx iris.Context) {
		ret, _ := DBGetHistory(chistory, ctx)
		ctx.JSON(ret)
	})
	// Start the server using a network address.
	app.Run(iris.Addr("0.0.0.0:8080"))
	//app.Run(iris.AutoTLS("0.0.0.0:443", "example.com", "admin@example.com"))
}
