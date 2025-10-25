package cmd

import (
	"job-queuer/app/croncontroller"
	repo "job-queuer/app/repository"
	"job-queuer/app/service"
	"job-queuer/app/util/database"
	"job-queuer/app/util/env"

	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "runjob",
	Short: "Run the job worker",
	Run: func(cmd *cobra.Command, args []string) {
		RunProcess()
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}

func RunProcess() {

	env.InitEnv()

	db, err := database.InitDB()
	if err != nil {
		panic("Failed to connect to database")
	}

	jobRepo := repo.NewJobRepo(db)
	jobService := service.NewJobService(jobRepo)

	croncontroller.StartJobCron(jobService)
}
