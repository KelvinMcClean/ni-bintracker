package service

import (
	"bintracker/internal/bintracker"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const credentialsFile = "credentials.json"

func AddBinsToCalendar(cfg bintracker.Config, bins []bintracker.Bin) {
	client := getClient()
	calendarService := getCalendar(cfg, client)
	addBinsToCalendar(calendarService, bins, cfg)
}

func getClient() *http.Client {
	// 2. Read the OAuth 2.0 credentials
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Request exactly the scope we need (Events scope allows reading and writing events)
	config, err := google.ConfigFromJSON(b, calendar.CalendarEventsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// 3. Initialize Google Calendar Service using user auth
	return getClientFromConfig(config)
}

func getClientFromConfig(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		log.Printf("Token not found or invalid, requesting new token: %v", err)
		tok = getTokenFromWeb(config)
		log.Printf("Saving new token to file: %s", tokFile)
		saveToken(tokFile, tok)

	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// Point redirect to a local server
	config.RedirectURL = "http://localhost:8080"

	codeChan := make(chan string)

	// Start a temporary HTTP server to receive the callback
	server := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			fmt.Fprintf(w, "<h1>Authorization Successful!</h1><p>You can close this tab and return to your terminal.</p>")
			codeChan <- code
		} else {
			fmt.Fprintf(w, "<h1>Authorization Failed.</h1>")
		}
	})

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser to authorize:\n\n%v\n\n", authURL)

	// Wait for the redirect to hit our local server
	authCode := <-codeChan

	// Shut down the server cleanly
	_ = server.Shutdown(context.Background())

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	log.Printf("Token saved to %s", path)
}

func getCalendar(config bintracker.Config, client *http.Client) *calendar.Service {
	ctx := context.Background()
	calendarService, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		panic(err)
	}
	return calendarService
}

func addBinsToCalendar(calendarService *calendar.Service, bins []bintracker.Bin, config bintracker.Config) {
	for _, bin := range bins {
		err := addEventToCalendar(calendarService, bin, config)
		if err != nil {
			fmt.Printf("Error adding event to calendar: %v\n", err)
		}
	}
}

func addEventToCalendar(srv *calendar.Service, collection bintracker.Bin, config bintracker.Config) error {
	councilLayout := "2006-01-02T15:04:05"
	parsedDate, err := time.Parse(councilLayout, collection.Date)
	if err != nil {
		parsedDate, err = time.Parse("2006-01-02", collection.Date)
		if err != nil {
			return fmt.Errorf("could not parse date %s: %v", collection.Date, err)
		}
	}

	colourEmoji, icon := getColourEmoji(collection.Colour)

	dateStr := parsedDate.Format("2006-01-02")

	event := &calendar.Event{
		Summary:     fmt.Sprintf("%s %s %s Collection %s", icon, colourEmoji, collection.Name, icon),
		Description: "Automated bin collection reminder.",
		Start: &calendar.EventDateTime{
			Date: dateStr,
		},
		End: &calendar.EventDateTime{
			Date: dateStr,
		},
		// Optional: Add a reminder for 6.30 PM the night before
		Reminders: &calendar.EventReminders{
			UseDefault: false,
			Overrides: []*calendar.EventReminder{
				{
					Method:  "popup",
					Minutes: ((5 * 60) + 30), // 6.30 PM the night before
				},
			},
			ForceSendFields: []string{"UseDefault"}, // Forces Go to serialize "useDefault": false
		},
	}
	exists, err := eventExists(srv, config.Calendar.ID, dateStr, event.Summary)
	if err != nil {
		return fmt.Errorf("error checking if event exists: %v", err)
	}
	if exists {
		fmt.Printf("Event already exists for %s on %s, skipping.\n", event.Summary, dateStr)
		return nil
	}
	fmt.Printf("Adding event to calendar: %s on %s\n", event.Summary, dateStr)
	resp, err := srv.Events.Insert(config.Calendar.ID, event).Do()
	if err != nil {
		return fmt.Errorf("could not create event: %v", err)
	}
	fmt.Printf("Event created: %s\n", resp.HtmlLink)
	return nil
}

func eventExists(srv *calendar.Service, calendarID string, dateStr string, summary string) (bool, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false, fmt.Errorf("invalid date format: %v", err)
	}

	// Set timeframe boundaries for the entire day (RFC3339 format)
	timeMin := t.Format(time.RFC3339)
	timeMax := t.Add(24 * time.Hour).Format(time.RFC3339)

	// List events on that day
	events, err := srv.Events.List(calendarID).
		TimeMin(timeMin).
		TimeMax(timeMax).
		SingleEvents(true).
		ShowDeleted(false).
		Do()
	if err != nil {
		return false, err
	}

	// Check if any existing event matches the summary
	for _, item := range events.Items {
		if item.Summary == summary {
			return true, nil
		}
	}

	return false, nil
}

func getColourEmoji(colour string) (string, string) {
	switch colour {
	case "Grey", "Black":
		return "🚮", "🗑️"
	case "Blue":
		return "♻️", "🗑️"
	case "Green", "Brown":
		return "🌍", "🗑️"
	case "Red":
		return "❤️", "🗑️"
	case "Yellow":
		return "🥛", "🗑️"
	default:
		return "🗑️", "🗑️"
	}
}
