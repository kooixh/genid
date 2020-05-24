package app

import (
	"fmt"
	c "github.com/kooixh/genid/pkg/constants"
	g "github.com/kooixh/genid/pkg/generators"
	"github.com/kooixh/genid/utils"
	"github.com/kooixh/genid/utils/redis"
)

func GenerateNewId(idChannel chan string, refillChannel chan int) {
	ids := redis.LPop(c.IdListKey)
	refillThreshold, err := redis.Get(c.IdListKey).Int()
	if err != nil {
		refillThreshold = 10
	}
	if int(redis.Length(c.IdListKey).Val()) < refillThreshold {
		go Refill(refillChannel)
	} else {
		idChannel <- ids.Val()
		refillChannel <- 0
	}
}

func Calibrate(offset int, initial int, total int, idType string, refillChannel chan int) {
	redis.Set(c.OffsetKey, offset)
	redis.Set(c.InitialKey, initial)
	redis.Set(c.TotalKey, total)
	redis.Set(c.GeneratorTypeKey, idType)
	go Refill(refillChannel)

}

func Refill(finished chan int) {
	offset, err := redis.Get(c.OffsetKey).Int()
	if err != nil {
		finished <- c.RedisErrorCode
	}
	total, err := redis.Get(c.TotalKey).Int()
	if err != nil {
		finished <- c.RedisErrorCode
	}
	initial, err := redis.Get(c.InitialKey).Int()
	if err != nil {
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
		finished <- 1
	}

	generatedIds := generator.Generate(rawIds)
	for _, elem := range generatedIds {
		redis.RPush(c.IdListKey, elem)
	}
	fmt.Println("generation done, increasing offset")
	redis.Set(c.OffsetKey, offset + 1)
	finished <- 0
}

