package benchmark

type Evolve func(scenario *BenchmarkScenario) bool

type evolution struct {
	next      Evolve
	iteration int
}

type Config struct {
	Num       int
	KeySize   int
	ValueSize int
}

type BenchmarkScenario struct {
	evolution
	Config
}

func Benchmark(next Evolve, num, keySize, valueSize int) BenchmarkScenario {
	return BenchmarkScenario{
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

func (s *BenchmarkScenario) Next() bool {
	return s.evolution.next(s)
}

func (s *BenchmarkScenario) Get() Config {
	s.iteration++
	return s.Config
}

func (s *BenchmarkScenario) execute(exec func(scenario Config)) {

	for s.Next() {
		exec(s.Get())
	}

}

type EvolutionBuilder []Evolve

func Evolution() EvolutionBuilder {
	return make([]Evolve, 0)
}

func (eb EvolutionBuilder) Add(stage Evolve) EvolutionBuilder {
	eb = append(eb, stage)
	return eb
}

func (eb EvolutionBuilder) Create() Evolve {
	if len(eb) == 0 {
		panic("cannot create evolution scenario without any instructions")
	}
	return func(scenario *BenchmarkScenario) bool {
		hasNext := true
		for i := 0; i < len(eb); i++ {
			next := eb[i](scenario)
			hasNext = next && hasNext
		}
		return hasNext
	}
}

func Limit(iteration int) func(scenario *BenchmarkScenario) bool {
	return func(scenario *BenchmarkScenario) bool {
		return scenario.iteration < iteration
	}
}

type op func(n int) int

func Add(m int) op {
	return func(n int) int {
		return n + m
	}
}

func Pow(m int) op {
	return func(n int) int {
		return n * m
	}
}

func Num(o op) func(scenario *BenchmarkScenario) bool {
	return func(scenario *BenchmarkScenario) bool {
		if scenario.iteration > 0 {
			scenario.Num = o(scenario.Num)
		}
		return true
	}
}

func Key(o op) func(scenario *BenchmarkScenario) bool {
	return func(scenario *BenchmarkScenario) bool {
		if scenario.iteration > 0 {
			scenario.KeySize = o(scenario.KeySize)
		}
		return true
	}
}

func Value(o op) func(scenario *BenchmarkScenario) bool {
	return func(scenario *BenchmarkScenario) bool {
		if scenario.iteration > 0 {
			scenario.ValueSize = o(scenario.ValueSize)
		}
		return true
	}
}
