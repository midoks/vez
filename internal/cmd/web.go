package cmd

import (
	"github.com/urfave/cli"

	"github.com/kataras/iris/v12"
)

var Service = cli.Command{
	Name:        "web",
	Usage:       "This command starts all services",
	Description: `Start Web services`,
	Action:      runWebService,
	Flags: []cli.Flag{
		stringFlag("config, c", "", "Custom configuration file path"),
	},
}

// Book example.
type Book struct {
	Title string `json:"title"`
}

func list(ctx iris.Context) {
	books := []Book{
		{"Mastering Concurrency in Go"},
		{"Go Design Patterns"},
		{"Black Hat Go"},
	}

	ctx.JSON(books)
	// 提示: 在服务器优先级和客户端请求中进行响应协商，
	// 以此来代替 ctx.JSON:
	// ctx.Negotiation().JSON().MsgPack().Protobuf()
	// ctx.Negotiate(books)
}

func create(ctx iris.Context) {
	var b Book
	err := ctx.ReadJSON(&b)
	// 提示: 使用 ctx.ReadBody(&b) 代替，来绑定所有类型的入参
	if err != nil {
		// ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("Book creation failure").DetailErr(err))
		// 提示: 如果仅有纯文本（plain text）错误响应，
		// 可使用 ctx.StopWithError(code, err)
		return
	}

	println("Received Book: " + b.Title)

	ctx.StatusCode(iris.StatusCreated)
}

func runWebService(c *cli.Context) error {

	app := iris.New()

	booksAPI := app.Party("/books")
	{
		// booksAPI.Use(iris.Compression)

		// GET: http://localhost:8080/books
		booksAPI.Get("/", list)
		// POST: http://localhost:8080/books
		booksAPI.Post("/", create)
	}

	app.Listen(":8080")
	return nil
}
