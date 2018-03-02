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

//func MongoController(c *mgo.Database) string {
func MongoController(c *mgo.Collection, id string) string {
	err := c.Insert(&Person{"New", "+55 53 8116 9639", id},
		&Person{"nla", "+55 53 8402 8510", id})
	fmt.Println(err)
	return "yes"
}

type Person struct {
	Name  string
	Phone string
	ID    string
}
type User struct {
	Username  string
	Firstname string
	Lastname  string
	City      string
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

	app.Get("/user/{id:long}", func(ctx iris.Context) {
		userID, _ := ctx.Params().GetInt64("id")
		ctx.Writef("User ID: %d", userID)
	})

	//curl -H "Content-Type: application/json" -H "TEST:B" -d '{"Firstname":"aaa", "Lastname":"bbb", "Username":"ccc","City":"ddd"}' http://localhost:8080/test/2X3190 -v
	app.Post("/test/{id}", middleware, func(ctx iris.Context) {
		//id, _ := ctx.Params().Get("id")
		id := ctx.Params().Get("id")
		fmt.Println(id)

		var userRead User
		ctx.ReadJSON(&userRead)
		fmt.Println(userRead)
		fmt.Println("where the input user name is", userRead.Username)

		ctx.Header("TEST", "RETURN")
		ctx.Writef("User ID:data")
	})

	//curl -H "Content-Type: application/json" -H "TEST:B" -d @a.json http://localhost:8080/test/2X3190 -v
	app.Get("/test/{id}", middleware, func(ctx iris.Context) {
		//id, _ := ctx.Params().Get("id")
		id := ctx.Params().Get("id")
		fmt.Println(id)

		var userRead User
		ctx.ReadJSON(&userRead)
		fmt.Println(userRead)

		ret := MongoController(c1, id)
		fmt.Println(ret)
		doe := User{
			Username:  "Johndoe",
			Firstname: "John",
			Lastname:  "Doe",
			City:      "Neither FBI knows!!!",
		}

		//it's a Header reuturn
		ctx.Header("TEST", "RETURN")
		//it's a StatusCode return
		//ctx.StatusCode(iris.StatusInternalServerError)
		//it's a json return
		ctx.JSON(doe)
		//it's a body return
		//ctx.Writef("User ID: %s", ret)
	})
	//curl -X POST -F 'img_avatar=@/go/src/github.com/aliciproject/mongocmd'  -H "TEST:B"  http://localhost:8080/image -vi
	app.Post("/image", middleware, func(ctx iris.Context) {
		// it will keep the filename mongocmd so the dest is just a directory that we must create by myself
		dest := "/go/src/github.com/aliciproject/mongocmdcopy"
		n, err := ctx.UploadFormFiles(dest)
		fmt.Println("show read file number", n)
		fmt.Println("show error", err)
		//ctx.Header("Content-Type", "multipart/form-data")

	})

	//curl  -H "TEST:B"  http://localhost:8080/image -vi
	app.Get("/image", middleware, func(ctx iris.Context) {
		file := "/go/src/github.com/aliciproject/mongocmd"
		ctx.SendFile(file, "c.txt")
	})
	// Start the server using a network address.
	app.Run(iris.Addr("0.0.0.0:8080"))
	//app.Run(iris.AutoTLS("0.0.0.0:443", "example.com", "admin@example.com"))
}
