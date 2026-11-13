package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func wrapper(done In, stage Stage) Stage {
	wrapped := func(in In) (out Out) {
		proxyChan := make(Bi)
		go func() {
			for {
				select {
				case <-done:
					close(proxyChan)
					for range in { //nolint:revive
					}
					return
				case val, ok := <-in:
					if !ok {
						close(proxyChan)
						return
					}
					select {
					case proxyChan <- val:
					case <-done:
						close(proxyChan)
						for range in { //nolint:revive
						}
						return
					}
				}
			}
		}()
		return stage(proxyChan)
	}
	return wrapped
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	firstChan := make(Bi)
	var stageOutChan Out = firstChan
	var stageInChan In = firstChan

	starter := func() {
		defer close(firstChan)
		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				select {
				case firstChan <- val:
				case <-done:
					return
				}
			}
		}
	}
	go starter()

	for _, st := range stages {
		stageOutChan = wrapper(done, st)(stageInChan)
		stageInChan = stageOutChan
	}

	return stageOutChan
}
