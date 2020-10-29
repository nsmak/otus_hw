package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil {
		ch := make(chan interface{})
		close(ch)
		return ch
	}

	// Воспользовался паттерном, который проходили на лекции
	take := func(done <-chan interface{}, valueStream <-chan interface{}) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for {
				select {
				case <-done:
					return
				default:
				}
				select {
				case <-done:
					return
				case val, ok := <-valueStream:
					if !ok { // эта проверка для первого тест-кейса, когда нет done канала
						return
					}
					select {
					case <-done:
						return
					case takeStream <- val:
					}
				}
			}
		}()
		return takeStream
	}

	inOutCh := in
	for _, stage := range stages {
		inOutCh = stage(take(done, inOutCh))
	}
	return inOutCh
}
