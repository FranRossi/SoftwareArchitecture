package pipes_and_filters

import "sync"

type Filter func(data any) error

type Pipeline struct {
	Filters []Filter
}

func (p *Pipeline) Use(f ...Filter) {
	p.Filters = append(p.Filters, f...)
}

func (p Pipeline) Run(input any) []error {

	var wg sync.WaitGroup
	wg.Add(len(p.Filters))
	out := make(chan error, len(p.Filters))

	// Runs each filters and saves the error to the out channel
	for _, f := range p.Filters {
		filter := f
		go func() {
			out <- filter(input)
			wg.Done()
		}()
	}

	// Waits for all the filters to finish
	wg.Wait()
	close(out)

	// Saves errors of all the filters
	var errors []error
	for err := range out {
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
