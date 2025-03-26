package main

import (
	"gone_playing/apps/exhibition"
	"sync"
)

func run(exhibitions []func()) {
	var wg sync.WaitGroup
	wg.Add(len(exhibitions))

	for _, run := range exhibitions {
		go func(runFunc func()) {
			defer wg.Done()
			runFunc()
		}(run)
	}

	wg.Wait()

}

func main() {
	exhibitions := []func(){
		exhibition.NewMocaTaipei().Run,
		exhibition.NewCKSMH().Run,
	}
	run(exhibitions)
}
