package cmd

import (
	"fmt"
	"job-queuer/app/router"
	"job-queuer/app/util/database"
	"job-queuer/app/util/env"
	"log"
	"net/http"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "runjob",
	Short: "Run the application worker",
	Run: func(cmd *cobra.Command, args []string) {
		StartServer()
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}

func StartServer() {
	env.InitEnv()
	fmt.Println("Environment variables loaded.")
	db, err := database.InitDB()
	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("Database connection established.")
	c := cron.New(cron.WithSeconds())

	r := router.NewRouter(db, c)

	c.Start()
	defer c.Stop()

	log.Printf("Starting server on %s...", "8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
