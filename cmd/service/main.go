package main

import (
	"fmt"
	app "l0/internal/app"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(16)

	App, err := app.CreateApp()

	if err != nil {
		fmt.Printf("error while creating app: %v", err)
		return
	}

	sc := *App.Sc

	if err := App.Run(); err != nil {
		fmt.Printf("error on app running: %v", err)
		return
	}

	sigChan := make(chan os.Signal, 1)

	<-sigChan

	App.DB.Close()
	sc.Close()
}
