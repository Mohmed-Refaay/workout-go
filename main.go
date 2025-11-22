package main

import (
	"backend-go/internals/app"
	"backend-go/internals/routes"
	"flag"
	"fmt"
	"net/http"
	"time"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Set port that server runs on")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()

	r := routes.SetupRoutes(app)
	s := &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%d", port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.Logger.Printf("App is running on %s\n", s.Addr)
	app.Logger.Fatal(s.ListenAndServe())
}
