package main

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"net/http"
	"net/http/httputil"
	"fmt"
	"log"
	"os"
)

type Todo struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello World</h1><div>Welcome to whereever you are</div>")
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		requestDump, _ := httputil.DumpRequest(r, true)
		log.Printf("%s\n",string(requestDump))
		fmt.Println(string(requestDump))
		handler.ServeHTTP(w, r)
	})
}

func openLogFile(logfile string) {
	if logfile != "" {
		lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

		if err != nil {
			log.Fatal("OpenLogfile: os.OpenFile:", err)
		}

		log.SetOutput(lf)
	}
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

	logPath := "development.log"

	openLogFile(logPath)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// serve HTTP
	http.HandleFunc("/", rootHandler)
	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", logRequest(http.DefaultServeMux))

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
