package main

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"net/http"
	"util/logger"
	"util/db"
)

type Todo struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type Search struct {
	Project   string `json:"Project"`
	Version string `json:"Version"`
	Type string   `json:"Type"`
	Name string   `json:"Name"`
}

func main() {
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

	searchType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Search",
		Fields: graphql.Fields{
			"Project": &graphql.Field{
				Type: graphql.String,
			},
			"Version": &graphql.Field{
				Type: graphql.String,
			},
			"Type": &graphql.Field{
				Type: graphql.String,
			},
			"Name": &graphql.Field{
				Type: graphql.String,
			},
			"Path": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	keywordIndexType := graphql.NewObject(graphql.ObjectConfig{
		Name: "KeywordIndex",
		Fields: graphql.Fields{
			"Keyword": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	keywordType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Keyword",
		Fields: graphql.Fields{
			"Project": &graphql.Field{
				Type: graphql.String,
			},
			"Version": &graphql.Field{
				Type: graphql.String,
			},
			"Type": &graphql.Field{
				Type: graphql.String,
			},
			"KeywordIndex": &graphql.Field{
				Type: graphql.String,
			},
			"Path": &graphql.Field{
				Type: graphql.String,
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
					return newTodo, nil
				},
			},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"search": &graphql.Field{
				Type: graphql.NewList(searchType),
				Args: graphql.FieldConfigArgument{
					"project": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"version": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"type": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					project, _ := params.Args["project"].(string)
					version, _ := params.Args["version"].(string)
					searchType, _ := params.Args["type"].(string)
					name, _ := params.Args["name"].(string)
					searches := db.Search(project, version, searchType, name, 10)
					return searches, nil
				},
      },
			"keywordIndex": &graphql.Field{
				Type: graphql.NewList(keywordIndexType),
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					name, _ := params.Args["name"].(string)
					keywordIndexes := db.KeywordIndex(name, 10)
					return keywordIndexes, nil
				},
			},
			"keyword": &graphql.Field{
				Type: graphql.NewList(keywordType),
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					name, _ := params.Args["name"].(string)
					keywords := db.Keyword(name, 10)
					return keywords, nil
				},
      },
		},
	})
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
	// $ curl -XPOST http://localhost:8080/graphql -H 'Content-Type: application/json' -d '{"query": "mutation M { newTodo: createTodo(text: \"This is a todo mutation example\") { text done } }"}'
	// Any of the above would return the following output:
	// {
	//   "data": {
	// 	   "newTodo": {
	// 	     "done": false,
	// 	     "text": "This is a todo mutation example"
	// 	   }
	//   }
	// }
	//curl -g -GET 'http://localhost:8080/graphql?query={lastTodo{text+done}}'
	//curl -g -GET 'http://localhost:8080/graphql?query={search(project:"angular",version:"1.6.0",type:"function",name:"a"){Project,Version,Name,Path,Type}}'
	//curl -g -GET 'http://localhost:8080/graphql?query={keywordIndex(name:"get"){Keyword}}'
	//curl -g -GET 'http://localhost:8080/graphql?query={keyword(name:"getAttributesObject"){Project,Version,KeywordIndex,Path,Type}}'
}
