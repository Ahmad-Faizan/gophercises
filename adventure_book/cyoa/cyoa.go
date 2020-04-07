package cyoa

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
)

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

const tmpl = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Choose Your Own Adventure</title>
	<style>
      body {
        font-family: helvetica, arial;
      }
      h1 {
        text-align:center;
      }
      .page {
        width: 80%;
        max-width: 500px;
        margin: auto;
        margin-bottom: 40px;
		padding: 50px;
		padding-top: 5px;
        background: #fff9c4;
        border: 1px solid #eee;
        box-shadow: 0 10px 6px -6px #777;
      }
      ul {
        border-top: 1px dotted #ccc;
        padding: 10px 0 0 0;
        -webkit-padding-start: 0;
      }
      li {
        padding-top: 10px;
      }
      a,
      a:visited {
        text-decoration: none;
        color: #6295b5;
      }
      a:active,
      a:hover {
        color: #f7a2a2;
      }
      p {
        text-indent: 1em;
      }
    </style>
</head>
<body>
	<section class="page">
		<h1>{{.Title}}</h1>
		<div>
			{{range .Paragraphs}}
				<p>{{.}}</p>
			{{end}}
		</div>
		<div>
			{{range .Options}}
				<ul>
					<li><a href="/{{.Chapter}}">{{.Text}}</a></li>
				</ul>
			{{end}}
		</div>
	</section>
</body>
</html>`

type storyHandler struct {
	s Story
}

var tpl *template.Template

func NewStoryHandler(story Story) http.Handler {
	return storyHandler{story}
}

func (h storyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.New("").Parse(tmpl))

	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:]

	if chapter, ok := h.s[path]; ok {
		err := tpl.Execute(w, chapter)
		if err != nil {
			fmt.Print(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	}
}

func JSONstory(reader io.Reader) (Story, error) {
	jsond := json.NewDecoder(reader)
	var story Story
	if err := jsond.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}
