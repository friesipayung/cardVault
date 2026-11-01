package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/friesipayung/cardVault/internal/db"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	dbPath := db.DefaultPath()
	if v := os.Getenv("CARDVAULT_DB"); v != "" {
		dbPath = v
	}
	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	if err := database.Migrate(); err != nil {
		log.Fatal(err)
	}

	cardSvc := NewCardService(database)
	providerSvc := NewProviderService(database)
	nfcSvc := NewNFCService(database)

	app := application.New(application.Options{
		Name:        "cardVault",
		Description: "Card loyalty management vault",
		Services: []application.Service{
			application.NewService(cardSvc),
			application.NewService(providerSvc),
			application.NewService(nfcSvc),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "cardVault",
		Width:            1100,
		Height:           720,
		BackgroundColour: application.NewRGB(15, 17, 26),
		URL:              "/",
	})

	fmt.Println("cardVault DB:", dbPath)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
