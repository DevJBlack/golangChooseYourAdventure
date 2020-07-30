package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultHanlderTmpl))
}

var tpl *template.Template

var defaultHanlderTmpl = `
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
					<li class="links" ><a href="/{{.Chapter}}">{{.Text}}</a></li>
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

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}

}

func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tpl, defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s      Story
	t      *template.Template
	pathFn func(r *http.Request) string
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)

	//                   ["intro"]
	if chapter, ok := h.s[path]; ok {
		err := tpl.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
}

func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}
