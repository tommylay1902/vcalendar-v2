package model

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GcClient struct {
	gcService  *calendar.Service
	httpClient *http.Client
	app        *application.App
	config     *oauth2.Config
	tokFile    string
}

func HasAuth() bool {
	tokFile := "token.json"

	_, err := tokenFromFile(tokFile)
	return err == nil
}

func InitializeClientGC() (*GcClient, error) {
	ctx := context.Background()
	b, err := os.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	gc := GcClient{}

	gc.config = config
	client := gc.getClient(config)
	gc.httpClient = client

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
		return nil, err
	}
	gc.gcService = srv
	gc.tokFile = "token.json"
	return &gc, nil
}

func (gc *GcClient) OpenBrowser() *GcClient {
	app := application.Get()
	fmt.Println("opening browser from native funciton")
	b, csErr := os.ReadFile("client_secret.json")
	if csErr != nil {
		fmt.Println("error reading getting client secret")
		panic(csErr)
	}
	config, cErr := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if cErr != nil {
		fmt.Println("error getting config")
		panic(cErr)
	}

	gc.config = config
	fmt.Println(gc.config.Endpoint)

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	openBrowserErr := app.Browser.OpenURL(authURL)
	if openBrowserErr != nil {
		fmt.Println("error opening browser")
		panic(openBrowserErr)
	}
	return gc
}

func (gc *GcClient) AddAuthCode(authCode string) error {
	// Exchange auth code for token
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("------------------------------%#v\n", err)
		}
	}()
	if gc.config == nil || gc.config.Exchange == nil {
		fmt.Println("HELLLOIFDOIFDSIFDJDSJEOJEJOE")
		panic("PANIIICCC")
	}
	tok, err := gc.config.Exchange(context.TODO(), authCode)
	if err != nil {
		return fmt.Errorf("unable to retrieve token from web: %v", err)
	}

	// Read client secret
	b, err := os.ReadFile("client_secret.json")
	if err != nil {
		return fmt.Errorf("unable to read client secret file: %v", err)
	}

	// Create config from JSON
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		return fmt.Errorf("unable to parse client secret file: %v", err)
	}
	gc.tokFile = "token.json"
	fmt.Println(gc.tokFile)
	// Save token first
	saveToken(gc.tokFile, tok)

	// Get HTTP client with the token
	client := config.Client(context.TODO(), tok)
	gc.httpClient = client

	// Create Calendar service
	srv, err := calendar.NewService(context.TODO(), option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve Calendar client: %v", err)
	}

	gc.gcService = srv
	fmt.Println("Successfully authenticated and created calendar service!")
	return nil
}

func (gc *GcClient) GetEventsForTheDay(date *time.Time) {
	var now time.Time
	var endOfDay time.Time
	if date == nil {
		now = time.Now().Add(time.Hour)

		endOfDay = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	} else {
		fmt.Println(date)
		now = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endOfDay = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	}

	fmt.Printf("Current time: %s\n", now)

	fmt.Printf("End of day: %s\n", endOfDay)

	events, err := gc.gcService.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(now.Format(time.RFC3339)).TimeMax(endOfDay.Format(time.RFC3339)).
		EventTypes("default").MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%v (%v)\n", item.Summary, date)
		}
	}
}

// Retrieve a token, saves the token, then returns the generated client.
func (gc *GcClient) getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	gc.tokFile = "token.json"
	tok, err := tokenFromFile(gc.tokFile)
	if err != nil {
		saveToken(gc.tokFile, tok)
	}

	return config.Client(context.Background(), tok)
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
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
