package main

import (
	"net/http"
	"os"
	"log"
	"encoding/json"
	"bytes"
	"fmt"
)

type User struct {
	DisplayName string
	Self string
}

type Issue struct {
	User User
}

type Message struct {
	Text string
}

func index(w http.ResponseWriter, r *http.Request) {
	// Get mattermost URL
	mattermostHookURL := r.URL.Query().Get("mattermost_hook_url")
	log.Println(mattermostHookURL) // FIXME Delete it

	// Parse JSON from JIRA
	decoder := json.NewDecoder(r.Body)
	var issue Issue
	decoder.Decode(&issue)
	log.Println(issue) // FIXME Delete it

	// Send data to Mattermost

	// Create message
	text := fmt.Sprintf(
		"[%s](%s)",
		issue.User.DisplayName,
		issue.User.Self,
	)

	message := Message{
		Text: text,
	}

	byteMessage, _ := json.Marshal(message)

	// Create http-client
	req, _ := http.NewRequest("POST", mattermostHookURL, bytes.NewBuffer(byteMessage))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	http.HandleFunc("/", index)
	http.ListenAndServe(":" + port, nil)
}
