// decode json into a map of string to struct
// for each json value(struct), display title, display series of story, display options
// Present options(present option text and option arc as link)
// If link is clicked to route, repeat process.

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
)


type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type StoryItem struct {
	Title string `json:"title"`
	Story []string `json:"story"`
	Options []Option `json:"options"`
}

type StoryMap map[string]StoryItem

func parseJsonStory(path string) (StoryMap, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var cyoa StoryMap
	err = json.NewDecoder(file).Decode(&cyoa)
	if err != nil {
		return nil, err
	}
	return cyoa, nil
}



func handleStory (storyMap StoryMap) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		arc := path.Base(r.URL.Path)
		if arc == "/" {
			arc = "intro"
		}
		fmt.Println("On arc: ", arc)
		storyItem, ok := storyMap[arc]
		if !ok {
			tmpl := template.Must(template.ParseFiles("templates/404.html"))
			tmpl.Execute(w, nil)
			return
		}
		tmpl := template.Must(template.ParseFiles("templates/story.html"))
		tmpl.Execute(w, storyItem)
	}

}

func main() {
	story, err := parseJsonStory("./story.json")
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	http.HandleFunc("/", handleStory(story))
	fmt.Println("Running on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}