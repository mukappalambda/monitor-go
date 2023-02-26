package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// hello is a component that displays a simple "Hello World!". A component is a
// customizable, independent, and reusable UI element. It is created by
// embedding app.Compo into a struct.
type hello struct {
	app.Compo
	comment interface{}
	a       int
}

// OnPreRender method for the hello component.
func (h *hello) OnNav(ctx app.Context) {
	// The OnPreRender method is called before the component is rendered. It is
	// used to initialize the component state.
	//
	// Here, the component state is initialized with the current time.

	// fetch data from jsonplaceholder comment's API

	// client := &http.Client{}
	// req, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/comments", nil)

	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// resp, err := client.Do(req)

	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// var data interface{}
	// err = json.NewDecoder(resp.Body).Decode(&data)

	// defer resp.Body.Close()

	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	fmt.Println("OnNav", h.comment)

	ctx.Async(func() {
		fmt.Println("dispatch", h.comment)
		ctx.Dispatch(func(ctx app.Context) {
			h.comment = 123
		})
	})

}

func (h *hello) OnMount(ctx app.Context) {
	// The OnMount method is called after the component is rendered for the
	// first time. It is used to initialize the component state.
	//
	// Here, the component state is initialized with the current time.
	h.a = 1
	fmt.Println("OnMount", h.comment)
	h.Update()
}

// The Render method is where the component appearance is defined. Here, a
// "Hello World!" is displayed as a heading.
func (h *hello) Render(ctx app.Context) app.UI {

	ctx.Async(func() {
		fmt.Println("dispatch", h.comment)
		ctx.Dispatch(func(ctx app.Context) {
			h.comment = 123
		})
	})

	fmt.Println("Render", h.comment)
	fmt.Println("Render", h.a)

	return app.Div().Body(
		app.Div().Body(
			app.A().Href("/sad").Text("Home"),
			app.If(h.comment == 123,
				app.Div().Text("yes"),
				// app.H1().Text("Hello World!"),
				// app.Range(h.comment.([]interface{})).Slice(
				// 	func(i int) app.UI {
				// 		return app.Div().Text(h.comment.([]interface{})[i].(map[string]interface{})["name"].(string))
				// 	},
				// ),
			).Else(app.Div().Text("Loading...")),
		),
	)
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
