package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	handler "job-queuer/app/apihandler"
	repo "job-queuer/app/repository"
	"job-queuer/app/service"
)

func NewRouter(db *gorm.DB) *mux.Router {
	r := mux.NewRouter()

	jobRepo := repo.NewJobRepo(db)
	jobService := service.NewJobService(jobRepo)
	jobHandler := handler.NewJobHandler(jobService)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")
	r.HandleFunc("/schedule", jobHandler.ScheduleJob).Methods("POST")
	r.HandleFunc("/status", jobHandler.CheckJobStatus).Methods("GET")
	return r
}
