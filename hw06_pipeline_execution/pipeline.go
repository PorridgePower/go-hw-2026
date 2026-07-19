package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	firstChan := make(Bi)
	var stageOutChan Out = firstChan
	var stageInChan In = firstChan

	for _, st := range stages {
		stageOutChan = st(stageInChan)
		stageInChan = In(stageOutChan)
	}

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
				firstChan <- val
			}
		}
	}
	go starter()

	return stageOutChan
	// return nil ????
}
