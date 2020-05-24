package generators

import (
	"github.com/kooixh/genid/utils"
	"strconv"
)

type Generator interface {
	Generate(ids []int64) []string
}

type NumericIdGenerator struct {
}

func (gen *NumericIdGenerator) Generate(ids []int64) []string {
	var alphaNumericResult []string
	for _, elem := range ids {
		alphaNumericResult = append(alphaNumericResult, strconv.FormatInt(elem, 10))
	}
	return utils.Shuffle(alphaNumericResult)
}

type AlphaNumericIdGenerator struct {
}

func (gen *AlphaNumericIdGenerator) Generate(ids []int64) []string {
	var alphaNumericResult []string
	for _, elem := range ids {
		alphaNumericResult = append(alphaNumericResult, strconv.FormatInt(elem, 36))
	}
	return utils.Shuffle(alphaNumericResult)
}