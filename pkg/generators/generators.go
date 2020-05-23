package generators

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Generator interface {
	Generate() string
}

type NumericIdGenerator struct {
	Prefix string
	Suffix string
}

func (gen *NumericIdGenerator) Generate() string {
	epoch := time.Now().Unix()
	rand.Seed(epoch)
	randomInt := rand.Intn(998) + 1
	return gen.Prefix + strconv.Itoa(int(epoch)) + fmt.Sprintf("%03d", randomInt) + gen.Suffix
}

type AlphaNumericIdGenerator struct {
	Prefix string
	Suffix string
}

func (gen *AlphaNumericIdGenerator) Generate() string {
	return "ABC123"
}