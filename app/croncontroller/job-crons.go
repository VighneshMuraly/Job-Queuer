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
	fmt.Println("Found " + fmt.Sprint(len(jobs)) + "queued jobs to process")
	svc.ProcessNextJob(jobs)
}
