package main

import (
	"log"
	"os"

	"mf-take-home-task/internal/database"
	"mf-take-home-task/internal/stock"

	"github.com/urfave/cli/v2"
)

func main() {
	// 1. Connect
	dbUrl := os.Getenv("DATABASE_URL") // looks unsafe in a plain string tho; magically reads the dockercomposeyaml
	if dbUrl == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	db, err := database.Connect(dbUrl)
	if err != nil {
		log.Fatalf("DB Connection failed: %v", err)
	}
	defer db.Close()

	// 2. Migrate
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// 3. Setup Logic
	repo := stock.NewRepository(db)

	// 4. CLI Definition
	app := &cli.App{
		Name:  "stock-cli",
		Usage: "Manage Mobilfox Stock",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all products",
				Action: func(c *cli.Context) error {
					return repo.ListProducts()
				},
			},
			{
				Name:  "increase",
				Usage: "Increase stock",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "id", Required: true},
					&cli.StringFlag{Name: "sku", Required: true},
					&cli.IntFlag{Name: "qty", Required: true},
					&cli.StringFlag{Name: "reason"},
				},
				Action: func(c *cli.Context) error {
					return repo.UpdateStock(
						c.String("id"),
						c.String("sku"),
						c.Int("qty"),
						true, // isIncrement = true
						c.String("reason"),
					)
				},
			},
			{
				Name:  "decrease",
				Usage: "Decrease stock",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "id", Required: true},
					&cli.StringFlag{Name: "sku", Required: true},
					&cli.IntFlag{Name: "qty", Required: true},
					&cli.StringFlag{Name: "reason"},
				},
				Action: func(c *cli.Context) error {
					return repo.UpdateStock(
						c.String("id"),
						c.String("sku"),
						c.Int("qty"),
						false, // isIncrement = false
						c.String("reason"),
					)
				},
			},
			{
				Name:  "report",
				Usage: "Print inventory report",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "top",
						Value: 5, // Default value
						Usage: "Number of top products to show",
					},
					&cli.IntFlag{
						Name:  "low-stock",
						Value: 10, // Default value
						Usage: "Threshold for low stock warning",
					},
				},
				Action: func(c *cli.Context) error {
					return repo.Report(
						c.Int("top"),
						c.Int("low-stock"),
					)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
