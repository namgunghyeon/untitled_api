# untitled_api
#### Test
```
How to make a HTTP request using cUrl
-------------------------------------
In `graphql-go-handler`, based on the GET/POST and the Content-Type header, it expects the input params differently.
This behaviour was ported from `express-graphql`.


1) using GET
$ curl -g -GET 'http://localhost:8080/graphql?query=mutation+M{newTodo:createTodo(text:"This+is+a+todo+mutation+example"){text+done}}'

2) using POST + Content-Type: application/graphql
$ curl -XPOST http://localhost:8080/graphql -H 'Content-Type: application/graphql' -d 'mutation M { newTodo: createTodo(text: "This is a todo mutation example") { text done } }'

3) using POST + Content-Type: application/json
$ curl -XPOST http://localhost:8080/graphql -H 'Content-Type: application/json' -d '{"query": "mutation M { newTodo: createTodo(text: \"This is a todo mutation example\") { text done } }"}'

4) curl -g -GET 'http://localhost:8080/graphql?query={search(project:"angular",version:"1.6.",type:"function",name:"a"){Project,Version,Name}}'


Any of the above would return the following output:
{
  "data": {
	  "newTodo": {
	   "done": false,
	   "text": "This is a todo mutation example"
	}
 }
}
```
