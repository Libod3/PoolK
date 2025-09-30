package pool

import "fmt"

func worker(p *WorkerPool) {
	for task := range p.taskQueue {
		p.freeWorkersCount.Add(^uint32(0))

		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("worker task panicked: %v", r)
				}

				p.freeWorkersCount.Add(1)

			}()

			task()

			if p.doneCallback != nil {
				p.doneCallback()
			}
		}()
	}
}
