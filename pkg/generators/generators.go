package generators

import (
	"github.com/kooixh/genid/utils"
	"strconv"
)

const base10 = 10
const base36 = 36

type Generator interface {
	Generate(ids []int64) []string
}

type NumericIdGenerator struct {
}

func (gen *NumericIdGenerator) Generate(ids []int64) []string {
	var alphaNumericResult []string
	for _, elem := range ids {
		alphaNumericResult = append(alphaNumericResult, strconv.FormatInt(elem, base10))
	}
	return utils.Shuffle(alphaNumericResult)
}

type AlphaNumericIdGenerator struct {
}

func (gen *AlphaNumericIdGenerator) Generate(ids []int64) []string {
	var alphaNumericResult []string
	for _, elem := range ids {
		alphaNumericResult = append(alphaNumericResult, strconv.FormatInt(elem, base36))
	}
	return utils.Shuffle(alphaNumericResult)
}