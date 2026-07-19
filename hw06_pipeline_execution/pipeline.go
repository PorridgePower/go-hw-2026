package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	firstChan := make(Bi)
	stageOutChan := firstChan
	stageInChan := firstChan

	for _, st := range stages {
		stageOutChan := st(stageInChan)
		stageInChan = stageOutChan
	}

	starter := func() {
		for {
			defer close(stageOutChan)
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

	// background := func() {
	// 	defer close(in)
	// 	<-done
	// }
	// go background()

	return stageOutChan
	// return nil ????
}
