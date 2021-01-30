package bench

type FilterType int

const (
	IN FilterType = iota + 1
	OUT
	LABEL
)

type Filter struct {
	Type   FilterType
	filter map[string]float64
	labels []string
}

func (f Filter) Apply(benchmark Benchmark, match *bool) {
	switch f.Type {
	case IN:
		for f, v := range f.filter {
			value, ok := benchmark.read(f)
			if ok && value == v {
				// all good here ...
			} else {
				*match = false
			}
		}
	case OUT:
		for f, v := range f.filter {
			value, ok := benchmark.read(f)
			if ok && value == v {
				*match = false
			}
		}
	case LABEL:
		// Note : labels are 'all or nothing'
		for _, v := range f.labels {
			if benchmark.hasLabel(v) {
				// all good ...
			} else {
				*match = false
			}
		}
	}

}

func Include(f map[string]float64) Filter {
	return Filter{
		Type:   IN,
		filter: f,
	}
}

func Exclude(f map[string]float64) Filter {
	return Filter{
		Type:   OUT,
		filter: f,
	}
}

func Label(f ...string) Filter {
	return Filter{
		Type:   LABEL,
		labels: f,
	}
}
