package cmd

import (
	"fmt"
	"job-queuer/app/router"
	"job-queuer/app/util/database"
	"job-queuer/app/util/env"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jobqueuer",
	Short: "Job Queuer is a CLI tool for managing jobs",
	Long:  `A simple CLI application to queue and manage jobs.`,
	Run: func(cmd *cobra.Command, args []string) {
		StartServer()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func StartServer() {
	env.InitEnv()

	db, err := database.InitDB()
	if err != nil {
		panic("Failed to connect to database")
	}
	r := router.NewRouter(db)
	log.Printf("Starting server on %s...", "8080")
	if err := http.ListenAndServe("localhost:8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
