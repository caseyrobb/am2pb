package am2pb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// Alertmanager - https://prometheus.io/docs/alerting/latest/configuration/#webhook_config
type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     string            `json:"startsAt"`
	EndsAt       string            `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
} // Alert represents a single alert from Alertmanager

type AlertmanagerWebhook struct {
	Status            string   `json:"status"`
	Version           string   `json:"version"`
	GroupKey          string   `json:"groupKey"`
	Receiver          string   `json:"receiver"`
	GroupLabels       struct{} `json:"groupLabels"`
	CommonLabels      struct{} `json:"commonLabels"`
	CommonAnnotations struct{} `json:"commonAnnotations"`
	ExternalURL       string   `json:"externalURL"`
	Alerts            []Alert  `json:"alerts"`
} // AlertmanagerWebhook represents the structure of an Alertmanager webhook payload

// Pushbullet note - https://docs.pushbullet.com/#parameters-for-different-types-of-pushes
type PushbulletNote struct {
	Body  string `json:"body"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

func StartServer(bearerToken string) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var amWebhook AlertmanagerWebhook
		if err := json.Unmarshal(body, &amWebhook); err != nil {
			http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
			return
		}

		var pbNote PushbulletNote
		for _, alert := range amWebhook.Alerts {
			fmt.Printf("Received alert: Status %s, Label: %s, Annotation %v\n",
				alert.Status, alert.Labels, alert.Annotations)

			pbNote.Body = strings.Join([]string{
				alert.Annotations["summary"], "\n",
				alert.Annotations["description"]}, "")

			pbNote.Title = strings.Join([]string{
				alert.Labels["severity"],
				alert.Labels["alertname"],
				alert.Status,
				alert.StartsAt}, " ")

			pbNote.Type = "note"

			err = postData(pbNote, bearerToken)
			if err != nil {
				http.Error(w, "Error making POST request: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Webhook processed successfully")
		}

	})

	log.Println("Server starting on port 5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func postData(data PushbulletNote, bearerToken string) error {
	client := &http.Client{}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.pushbullet.com/v2/pushes", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("POST request response: %s\n", string(body))

	return nil
}
