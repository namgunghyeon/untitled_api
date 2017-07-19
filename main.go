package main

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"net/http"
	"util/logger"
)

type Todo struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func main() {

	// define custom GraphQL ObjectType `todoType` for our Golang struct `Todo`
	// Note that
	// - the fields in our todoType maps with the json tags for the fields in our struct
	// - the field type matches the field type in our struct
	todoType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Todo",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"text": &graphql.Field{
				Type: graphql.String,
			},
			"done": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	})

	// root mutation
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createTodo": &graphql.Field{
				Type: todoType, // the return type for this field
				Args: graphql.FieldConfigArgument{
					"text": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"text2": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {

          text, _ := params.Args["text"].(string)
					text2, _ := params.Args["text2"].(string)
					newTodo := &Todo{
						ID:   "id0001",
						Text: text + "_" + text2,
						Done: false,
					}

					// return the new Todo object that we supposedly save to DB
					// Note here that
					// - we are returning a `Todo` struct instance here
					// - we previously specified the return Type to be `todoType`
					// - `Todo` struct maps to `todoType`, as defined in `todoType` ObjectConfig`
					return newTodo, nil
				},
			},
		},
	})

	// root query
	// we just define a trivial example here, since root query is required.
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"lastTodo": &graphql.Field{
				Type: todoType,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					todo := &Todo{
						ID:   "id0001",
						Text: "12345",
						Done: false,
					}
					return todo, nil
				},
			},
      "latestPost": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "world", nil
				},
      },
		},
	})

	// define schema
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	if err != nil {
		panic(err)
	}

	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	logPath := "./logs/development.log"

	logger.OpenLogFile(logPath)

	// serve HTTP
	http.HandleFunc("/", logger.RootHandler)
	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", logger.LogRequest(http.DefaultServeMux))

	// How to make a HTTP request using cUrl
	// -------------------------------------
	// In `graphql-go-handler`, based on the GET/POST and the Content-Type header, it expects the input params differently.
	// This behaviour was ported from `express-graphql`.
	//
	//
	// 1) using GET
	// $ curl -g -GET 'http://localhost:8080/graphql?query=mutation+M{newTodo:createTodo(text:"This+is+a+todo+mutation+example"){text+done}}'
	//
	// 2) using POST + Content-Type: application/graphql
	// $ curl -XPOST http://localhost:8080/graphql -H 'Content-Type: application/graphql' -d 'mutation M { newTodo: createTodo(text: "This is a todo mutation example") { text done } }'
	//
	// 3) using POST + Content-Type: application/json
	// $ curl -XPOST http://localhost:8080/graphql -H 'Content-Type: application/json' -d '{"query": "mutation M { newTodo: createTodo(text: \"This is a todo mutation example\") { text done } }"}'
	//
	// Any of the above would return the following output:
	// {
	//   "data": {
	// 	   "newTodo": {
	// 	     "done": false,
	// 	     "text": "This is a todo mutation example"
	// 	   }
	//   }
	// }
}
