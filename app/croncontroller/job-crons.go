package croncontroller

import (
	"job-queuer/app/service"
	"time"
)

func StartJobCron(svc service.JobService) {
	ticker := time.NewTicker(20 * time.Second)
	go func() {
		for range ticker.C {
			jobs, err := svc.GetQueuedJobs()
			if err != nil || len(jobs) == 0 {
				continue
			}
			// Process jobs with quota and bias
			svc.ProcessJobsWithQuota(jobs)
		}
	}()
}
