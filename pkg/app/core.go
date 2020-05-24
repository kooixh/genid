package app

import (
	"fmt"
	c "github.com/kooixh/genid/pkg/constants"
	g "github.com/kooixh/genid/pkg/generators"
	"github.com/kooixh/genid/utils"
	"github.com/kooixh/genid/utils/redis"
)

type CalibrationSettings struct {
	Initial int
	Offset int
	Total int
	IdType string
}

func (cal *CalibrationSettings) Validate() (bool, string) {
	if cal.IdType != c.GeneratorTypeAlphaNum && cal.IdType != c.GeneratorTypeNum {
		return false, "Error value in --type, (alphanum or num)"
	}
	if cal.Total > cal.Initial {
		return false, "Error value in --total, cannot be greater than initial"
	}
	return true, ""
}

func GenerateNewId(idChannel chan string, refillChannel chan int) {
	ids := redis.LPop(c.IdListKey)
	refillThreshold, err := redis.Get(c.IdListKey).Int()
	if err != nil {
		refillThreshold = 10
	}
	if int(redis.Length(c.IdListKey).Val()) < refillThreshold {
		go Refill(refillChannel)
	} else {
		refillChannel <- 0
	}
	idChannel <- ids.Val()
}

func Calibrate(cal CalibrationSettings, refillChannel chan int) {
	redis.Set(c.OffsetKey, cal.Offset)
	redis.Set(c.InitialKey, cal.Initial)
	redis.Set(c.TotalKey, cal.Total)
	redis.Set(c.GeneratorTypeKey, cal.IdType)
	go Refill(refillChannel)

}

func Refill(finished chan int) {

	offset, errOffset := redis.Get(c.OffsetKey).Int()
	total, errTotal := redis.Get(c.TotalKey).Int()
	initial, errInitial := redis.Get(c.InitialKey).Int()
	if errOffset != nil || errTotal != nil || errInitial != nil {
		finished <- c.RedisErrorCode
	}

	rawIds := utils.GenerateNewIdSet(total, offset, initial)

	generatorType := redis.Get(c.GeneratorTypeKey).Val()
	var generator g.Generator
	if generatorType == c.GeneratorTypeAlphaNum {
		generator = new(g.AlphaNumericIdGenerator)
	} else if generatorType == c.GeneratorTypeNum {
		generator = new(g.NumericIdGenerator)
	} else {
		fmt.Println("error, generator type is unknown")
		finished <- c.RedisErrorUnknownTypeCode
	}

	generatedIds := generator.Generate(rawIds)
	for _, elem := range generatedIds {
		redis.RPush(c.IdListKey, elem)
	}
	fmt.Println("generation done, increasing offset")
	redis.Set(c.OffsetKey, offset + 1)
	finished <- 0
}
