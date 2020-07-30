package main

import (
	"chooseYourAdventure/cyoa"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := flag.Int("port", 3000, "the port to start the CYOA eb application on")
	filename := flag.String("file", "gopher.json", "the JSON file with CYOA story")
	flag.Parse()
	fmt.Printf("Using the story in %s.\n", *filename)

	f, err := os.Open(*filename)
	if err != nil {
		fmt.Println("Here is the error")
		panic(err)
	}

	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(err)
	}

	tpl := template.Must(template.New("").Parse(storyTmpl))

	h := cyoa.NewHandler(story, cyoa.WithTemplate(tpl), cyoa)
	fmt.Printf("Starting the server at %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}

var storyTmpl = `
	<html>
	<head>
		<meta charset="utf-8">
		<title>Chose Your Own Adventure</title>
	</head>
	<body>
		<div class="context">
			<h1>{{.Title}}</h1>
				{{range .Paragraphs}}
					<p>{{.}}</p>
				{{end}}
			<ul>
				{{range .Options}}
					<li class="links" ><a href="/story/{{.Chapter}}">{{.Text}}</a></li>
				{{end}}
			</ul>  
		</div>
		<style>
			body {
				background-color: #708090;
				display: flex;
				justify-content: center;
				align-items: center;
			}
		
			.context {
				max-width: 60rem;
				margin: auto;
				padding: 10px;
				border: 20px solid #6699cc;
				background-color: #66ccff;
			}
		
			h1 {
				text-align: center;
			}
		
			a {
				padding: 1rem;
				text-decoration: none;
			}
		
			.links {
				font-size: 1rem;
			}
	   </style>
	</body>
	</html>
`
