package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Yakumo-zi/gtool/pkg/slice"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

const GithubApiUrl = "https://api.github.com/users/%s/events"

type EventType string

const PushEvent EventType = "PushEvent"
const WatchEvent EventType = "WatchEvent"
const CreateEvent EventType = "CreateEvent"

type Event struct {
	Id        string    `json:"id"`
	Type      EventType `json:"type"`
	Actor     Actor     `json:"actor"`
	Repo      Repo      `json:"repo"`
	Payload   Payload   `json:"payload"`
	Public    bool      `json:"public"`
	CreatedAt time.Time `json:"created_at"`
	Org       Org       `json:"org,omitempty"`
}

func main() {
	userName := flag.String("user", "", "github username")
	flag.Parse()
	if *userName == "" {
		Usage()
		return
	}
	body, err := http.Get(fmt.Sprintf(GithubApiUrl, *userName))
	if err != nil {
		log.Fatal(err)
	}
	defer body.Body.Close()
	var events []Event
	json.NewDecoder(body.Body).Decode(&events)
	output := make([]string, len(events))
	for _, event := range events {
		switch event.Type {
		case PushEvent:
			s := fmt.Sprintf("Pushed %d commits to %s at %s", len(event.Payload.Commits), event.Repo.Name, event.CreatedAt.Format("2006/01/02"))
			output = append(output, s)
		case WatchEvent:
			s := fmt.Sprintf("Starred %s at %s", event.Repo.Name, event.CreatedAt.Format("2006/01/02"))
			output = append(output, s)
		case CreateEvent:
			s := fmt.Sprintf("Created repository %s at %s", event.Repo.Name, event.CreatedAt.Format("2006/01/02"))
			output = append(output, s)
		}
	}
	output = slice.Distinct(output)
	for _, v := range output {
		if v != "" {
			fmt.Println(v)
		}
	}
}

func Usage() {
	fmt.Fprintf(os.Stdin, "Usage of %s:\n", path.Base(os.Args[0]))
	fmt.Fprintf(os.Stdin, "\t%s <github username>\n", path.Base(os.Args[0]))
}
