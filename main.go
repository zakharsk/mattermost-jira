package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"log"
)

// Get it

type Data struct {
	WebhookEvent string
	User         struct {
		Name        string
		DisplayName string
	}
	Issue struct {
		Self   string
		Key    string
		Fields struct {
			Summary string
		}
	}
	Comment struct {
		Body string
	}
	Changelog []struct{
		Items struct{
			Field string
			FromString string
			ToString string
		      }
	}
}

// Send it

type Message struct {
	Text string
}

func index(w http.ResponseWriter, r *http.Request) {
	// Get mattermost URL
	mattermostHookURL := r.URL.Query().Get("mattermost_hook_url")

	// Parse JSON from JIRA
	decoder := json.NewDecoder(r.Body)
	var data Data
	decoder.Decode(&data)

	// Get JIRA URL
	u, _ := url.Parse(data.Issue.Self)

	// Select action
	var action, appendix string
	switch data.WebhookEvent {
	case "jira:issue_created":
		action = "created"
	case "jira:issue_updated":
		action = "updated"
	case "jira:issue_deleted":
		action = "deleted"
	}

	if len(data.Comment.Body) > 0 {
		action = "commented"
		appendix = fmt.Sprintf(" with\n>%s", data.Comment.Body)
	}

	// Create message for Mattermost
	text := fmt.Sprintf(
		//[UserFirstName UserSecondName](user_link) commented issue [[TSK-158]](issue_link) "Test task" with "Test comment"
		"[%s](%s://%s/secure/ViewProfile.jspa?name=%s) %s [%s](%s://%s/browse/%s) \"%s\"%s //\n%s//",
		data.User.DisplayName,
		u.Scheme,
		u.Host,
		data.User.Name,
		action,
		data.Issue.Key,
		u.Scheme,
		u.Host,
		data.Issue.Key,
		data.Issue.Fields.Summary,
		appendix,
		data.Changelog,


	)

	message := Message{
		Text: text,
	}

	byteMessage, _ := json.Marshal(message)


	if len(mattermostHookURL) > 0 {
		// Create http-client
		req, _ := http.NewRequest("POST", mattermostHookURL, bytes.NewBuffer(byteMessage))
		req.Header.Set("Content-Type", "application/json")

		// Send data to Mattermost
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	} else {
		log.Println(text)
	}


}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	http.HandleFunc("/", index)
	http.ListenAndServe(":"+port, nil)
}
