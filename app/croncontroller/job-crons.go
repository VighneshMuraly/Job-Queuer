package croncontroller

import (
	"fmt"
	"job-queuer/app/service"
)

func StartJobCron(svc service.JobService) {
	jobs, err := svc.GetQueuedJobs()
	if err != nil || len(jobs) == 0 {
		fmt.Println("No queued jobs to process")
		return
	}
	svc.ProcessNextJob(jobs)
}
