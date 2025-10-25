package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	handler "job-queuer/app/apihandler"
	"job-queuer/app/croncontroller"
	repo "job-queuer/app/repository"
	"job-queuer/app/service"
)

func NewRouter(db *gorm.DB, c *cron.Cron) *mux.Router {
	r := mux.NewRouter()

	jobRepo := repo.NewJobRepo(db)
	jobService := service.NewJobService(jobRepo)
	jobHandler := handler.NewJobHandler(jobService)

	c.AddFunc("*/20 * * * * *", func(jobService service.JobService) func() {
		return func() {
			croncontroller.StartJobCron(jobService)
		}
	}(jobService))

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")
	r.HandleFunc("/schedule", jobHandler.ScheduleJob).Methods("POST")
	r.HandleFunc("/status", jobHandler.CheckJobStatus).Methods("GET")
	return r
}
