package benchmark

// Evolve represents an evolution logic for a given scenario parameters
type Evolve func(scenario *Scenario) bool

type evolution struct {
	next      Evolve
	iteration int
}

// Config carries the benchmark scenario evolution configuration
type Config struct {
	Num       int
	KeySize   int
	ValueSize int
}

// Scenario represents a benchmark scenario
type Scenario struct {
	name string
	evolution
	Config
}

// Benchmark creates a new benchmark scenario
func Benchmark(next Evolve, num, keySize, valueSize int) Scenario {
	return Scenario{
		evolution: evolution{
			next: next,
		},
		Config: Config{
			Num:       num,
			KeySize:   keySize,
			ValueSize: valueSize,
		},
	}
}

func (s *Scenario) Name(name string) Scenario {
	s.name = name
	return *s
}

// next evolves the given scenario for one iteration
func (s *Scenario) next() bool {
	return s.evolution.next(s)
}

// get returns the scenario config
func (s *Scenario) get() Config {
	s.iteration++
	return s.Config
}

// execute executes all the scenarios based on the corresponding configuration
func (s *Scenario) execute(exec func(scenario Config)) {
	for s.next() {
		exec(s.get())
	}
}

// Builder facilitates the creation of evolution logic for the becnhmark scenarios
type Builder []Evolve

// Evolution creates a new scenario evolution builder
func Evolution() Builder {
	return make([]Evolve, 0)
}

// add adds an evolution stage to the builder
func (eb Builder) add(stage Evolve) Builder {
	eb = append(eb, stage)
	return eb
}

// create creates a new scenario based on the builder properties
func (eb Builder) create() Evolve {
	if len(eb) == 0 {
		panic("cannot create evolution scenario without any instructions")
	}
	return func(scenario *Scenario) bool {
		hasNext := true
		for i := 0; i < len(eb); i++ {
			next := eb[i](scenario)
			hasNext = next && hasNext
		}
		return hasNext
	}
}

// limit specifies the  number of evolutions for a given scenario
func limit(iteration int) func(scenario *Scenario) bool {
	return func(scenario *Scenario) bool {
		return scenario.iteration < iteration
	}
}

// num specifies the number of elements as an evolution parameter
func num(o op) func(scenario *Scenario) bool {
	return func(scenario *Scenario) bool {
		if scenario.iteration > 0 {
			scenario.Num = o(scenario.Num)
		}
		return true
	}
}

// key specifies the elements keys size as an evolution parameter
func key(o op) func(scenario *Scenario) bool {
	return func(scenario *Scenario) bool {
		if scenario.iteration > 0 {
			scenario.KeySize = o(scenario.KeySize)
		}
		return true
	}
}

// value specifies the elements values size as an evolution parameter
func value(o op) func(scenario *Scenario) bool {
	return func(scenario *Scenario) bool {
		if scenario.iteration > 0 {
			scenario.ValueSize = o(scenario.ValueSize)
		}
		return true
	}
}

// op identifies an operation on a generic integer
type op func(n int) int

// add specifies the addition evolution logic for scenario values
func add(m int) op {
	return func(n int) int {
		return n + m
	}
}

// pow specifies the exponential logic for scenario values evolution
func pow(m int) op {
	return func(n int) int {
		return n * m
	}
}
