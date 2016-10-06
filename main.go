package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"log"
	"strings"
)

// Get it
type Data struct {
	WebhookEvent string
	User         struct {
		Name        string
		AvatarUrls map[string]string
		DisplayName string
	}
	Issue struct {
		Self   string
		Key    string
		Fields struct {
			Issuetype struct{
				IconUrl string
				Name string
				  }
			Summary string
		}
	}
	Comment struct {
		Body string
	}
	Changelog struct{
		Items []struct{
			Field string
			FromString string
			ToString string
		      }
	}
}

// Send it
type Message struct {
	Text     string
	Username string
	Icon_url string
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
	var action, comment string
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
		comment = fmt.Sprintf("\n>%s", data.Comment.Body)
	}

	// Process changelog
	var changelog string
	if len(data.Changelog.Items) > 0 {
		for _, item := range data.Changelog.Items {
			itemName := strings.ToUpper(string(item.Field[0])) + item.Field[1:]
			if item.FromString == "" {
				item.FromString = "None"
			}
			changelog += fmt.Sprintf(
				"\n%s: ~~%s~~ %s",
				itemName,
				item.FromString,
				item.ToString,
			)
		}
	}

	// Create message for Mattermost
	text := fmt.Sprintf(
		//![user_icon](user_icon_link)[UserFirstName UserSecondName](user_link) commented task ![task_icon](task_icon link)[TSK-42](issue_link) "Test task"
		//Status: ~~Done~~ Finished
		//>Comment text
		"![user_icon](%s) [%s](%s://%s/secure/ViewProfile.jspa?name=%s) %s %s ![task_icon](%s) [%s](%s://%s/browse/%s) \"%s\"%s%s",
		data.User.AvatarUrls["16x16"],
		data.User.DisplayName,
		u.Scheme,
		u.Host,
		data.User.Name,
		action,
		strings.ToLower(data.Issue.Fields.Issuetype.Name),
		data.Issue.Fields.Issuetype.IconUrl,
		data.Issue.Key,
		u.Scheme,
		u.Host,
		data.Issue.Key,
		data.Issue.Fields.Summary,
		changelog,
		comment,
	)

	message := Message{
		Text: text,
		Username: "JIRA",
		Icon_url: "https://design.atlassian.com/images/logo/favicon.png",
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
		log.Println(string(byteMessage))
	}


}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	http.HandleFunc("/", index)
	http.ListenAndServe(":" + port, nil)
}
