package pool

type Pool struct {
	opts PoolOptions
	ch   chan func()
}

type PoolOptions struct {
	Max int
}

func New(opts *PoolOptions) *Pool {
	pool := &Pool{
		ch:   make(chan func()),
		opts: *opts,
	}
	go pool.schedule()
	return pool
}

func (p *Pool) schedule() {
	for i := 0; i < p.opts.Max; i++ {
		go func() {
			for f := range p.ch {
				f()
			}
		}()
	}
}

func (p *Pool) Go(fn func()) {
	p.ch <- fn
}
