package workers

import (
	"testing"
)

func BenchmarkWorkerPool(b *testing.B) {
	const numWorkers = 4
	const numJobs = 1000

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		jobs := make(chan Task, numJobs)
		results := make(chan Result, numJobs)

		CreatePool(numWorkers, jobs, results)

		for j := 1; j <= numJobs; j++ {
			jobs <- Task{ID: j, Data: "payload"}
		}
		close(jobs)

		for j := 0; j < numJobs; j++ {
			<-results
		}
	}
}