package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for i := range stages {
		in = func(in In) Out {
			dataCh := make(Bi)
			go func() {
				defer close(dataCh)
				for {
					select {
					case <-done:
						return
					case v, ok := <-in:
						if !ok {
							return
						}
						dataCh <- v
					}
				}
			}()

			return stages[i](dataCh)
		}(in)
	}

	return in
}
