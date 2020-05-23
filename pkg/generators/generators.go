package generators

type Generator interface {
	Generate() string
}

type StringIdGenerator struct {
}

func (gen *StringIdGenerator) Generate() string {
	return "10001"
}
