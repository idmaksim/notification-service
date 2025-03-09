package main

import iapp "github.com/idmaksim/notification-service/internal/app"

func main() {
	app := iapp.NewApp()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
