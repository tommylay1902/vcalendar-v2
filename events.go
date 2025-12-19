package main

import (
	"fmt"

	"vcalendar-v2/model"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func setupEvents(app *application.App, gcClient *model.GcClient) {
	app.Event.On("vcalendar-v2:auth-needed", func(e *application.CustomEvent) {
		isAuthenticated := model.HasAuth()
		app.Event.Emit("vcalendar-v2:token-needed", model.GoogleAuth{
			TokenNeeded: isAuthenticated,
		})
		if isAuthenticated {
			client, err := model.InitializeClientGC()
			if err != nil {
				fmt.Println("error initializing gc client")
				panic(err)
			}
			gcClient = client

		} else {
			gcClient = gcClient.OpenBrowser()
		}
	})

	app.Event.On("vcalendar-v2:auth-code-token", func(event *application.CustomEvent) {
		token := event.Data.(model.AuthCodeToken)
		gcClient.AddAuthCode(token.Token)
		app.Event.Emit("vcalendar-v2:token-needed", model.GoogleAuth{
			TokenNeeded: true,
		})
	})
}
