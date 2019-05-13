package main

import (
	"gopaytest/container"
	"gopaytest/routes"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

func main() {
	// We could use a library to unmarshal env vars (e.g. Netflix/go-env) but since I'm doing custom validation on these
	// and I only have two variables I decided to do it manually.
	// If you end up having more than two environment variables may be worth using a library.
	httpPortRe := "[0-9]{2,5}"
	httpPort := os.Getenv("HTTP_PORT")
	if !regexp.MustCompile(httpPortRe).MatchString(httpPort) {
		log.Fatalf("Invalid HTTP port supplied %q (%v)", httpPort, httpPortRe)
	}

	baseURL, err := url.ParseRequestURI(os.Getenv("BASE_URL"))
	if err != nil {
		log.Fatal("BASE_URL environment variable is not valid")
	}

	c := container.NewContainer(baseURL.String())

	router := routes.NewRouter(c)

	log.Fatal(http.ListenAndServe(":"+httpPort, router))
}
