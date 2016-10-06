package main

import (
	"net/http"
	"os"
	"encoding/json"
	"bytes"
	"fmt"
	"net/url"
)

// Get it

type Data struct {
	WebhookEvent string
	User User
	Issue Issue
}

type User struct {
	Name string
	DisplayName string
}

type Issue struct {
	Self string
	Key string
	Fields Fields
}

type Fields struct {
	Summary string
}

// Send it

type Message struct {
	Text string
}

func index(w http.ResponseWriter, r *http.Request) {
	// Get mattermost URL
	mattermostHookURL := r.URL.Query().Get("mattermost_hook_url")
	//log.Println(mattermostHookURL) // FIXME Delete it

	// Parse JSON from JIRA
	decoder := json.NewDecoder(r.Body)
	var data Data
	decoder.Decode(&data)
	//log.Println(data) // FIXME Delete it

	// Get JIRA URL

	u, _ := url.Parse(data.Issue.Self)
	//log.Println(u.Scheme + u.Host)

	// Send data to Mattermost

	// Create message
	text := fmt.Sprintf(
		"[[%s] %s](%s://%s/browse/%s) %s by [%s](%s://%s/secure/ViewProfile.jspa?name=%s)",
		data.Issue.Key,
		data.Issue.Fields.Summary,
		u.Scheme,
		u.Host,
		data.Issue.Key,
		data.WebhookEvent,
		data.User.DisplayName,
		u.Scheme,
		u.Host,
		data.User.Name,
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
