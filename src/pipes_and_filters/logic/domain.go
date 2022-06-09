package logic

import "sync"

// TODO revisar si esta linea esta bien
type Any = interface{}

type Filter func(data Any, params map[string]string) error

type Pipeline struct {
	Filters []Filter
}

func (p *Pipeline) Use(f ...Filter) {
	p.Filters = append(p.Filters, f...)
}

func (p Pipeline) Run(input Any) []error {

	var wg sync.WaitGroup
	wg.Add(len(p.Filters))
	out := make(chan error, len(p.Filters))

	// Runs each filters and saves the error to the out channel
	for _, f := range p.Filters {
		filter := f
		go func() {
			m := make(map[string]string)
			out <- filter(input, m)
			wg.Done()
		}()
	}

	// Waits for all the filters to finish
	wg.Wait()
	close(out)
	var errors []error

	// Saves errors of all the filters
	for err := range out {
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
