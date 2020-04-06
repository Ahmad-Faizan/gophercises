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
</head>
<body>
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
