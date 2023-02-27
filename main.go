package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// hello is a component that displays a simple "Hello World!". A component is a
// customizable, independent, and reusable UI element. It is created by
// embedding app.Compo into a struct.

type Comment struct {
	Postid int    `json:"postId,omitempty"`
	Id     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	Body   string `json:"body,omitempty"`
}

type hello struct {
	app.Compo
	comments []Comment
}

type parent struct {
	app.Compo
}

type child struct {
	app.Compo
	name string
}

func (p *parent) Render() app.UI {
	firstChild := &child{}
	secondChild := &child{}

	firstChild.name = "alex"
	secondChild.name = "bob"

	return app.Div().Body(
		app.Div().Text("My children:"),
		firstChild,
		secondChild,
	)
}

func (c *child) Render() app.UI {
	return app.Div().Text(c.name)
}

// The Render method is where the component appearance is defined. Here, a
// "Hello World!" is displayed as a heading.
func (h *hello) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Comments"),
		app.Range(h.comments).Slice(func(i int) app.UI {
			return app.Li().Text(h.comments[i])
		}),
	)

}

func (h *hello) OnMount(ctx app.Context) {
	ctx.Async(func() {
		resp, err := http.Get("https://jsonplaceholder.typicode.com/comments/")
		if err != nil {
			log.Fatal(err)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		var comments []Comment
		if err = json.Unmarshal(body, &comments); err != nil {
			log.Fatal(err)
			return
		}

		ctx.Dispatch(func(ctx app.Context) {
			h.comments = comments
			h.Update()
		})
	})

}

// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.
func main() {
	// The first thing to do is to associate the hello component with a path.
	//
	// This is done by calling the Route() function,  which tells go-app what
	// component to display for a given path, on both client and server-side.
	app.Route("/", &hello{})
	app.Route("/parent", &parent{})

	// Once the routes set up, the next thing to do is to either launch the app
	// or the server that serves the app.
	//
	// When executed on the client-side, the RunWhenOnBrowser() function
	// launches the app,  starting a loop that listens for app events and
	// executes client instructions. Since it is a blocking call, the code below
	// it will never be executed.
	//
	// When executed on the server-side, RunWhenOnBrowser() does nothing, which
	// lets room for server implementation without the need for precompiling
	// instructions.
	app.RunWhenOnBrowser()

	// Finally, launching the server that serves the app is done by using the Go
	// standard HTTP package.
	//
	// The Handler is an HTTP handler that serves the client and all its
	// required resources to make it work into a web browser. Here it is
	// configured to handle requests with a path that starts with "/".
	http.Handle("/", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
