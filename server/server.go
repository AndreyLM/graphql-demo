package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	graphql_demo "github.com/andreylm/graphql-demo"
	"github.com/andreylm/graphql-demo/api/dal"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := dal.Connect()
	if err != nil {
		panic(err)
	}
	resolver := graphql_demo.NewResolver(db)

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(graphql_demo.NewExecutableSchema(graphql_demo.Config{Resolvers: resolver})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
