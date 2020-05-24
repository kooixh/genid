package app

import (
	"errors"
	"fmt"
	"github.com/kooixh/genid/utils"
	"github.com/kooixh/genid/utils/redis"
	c "github.com/kooixh/genid/pkg/constants"
)

func Calibrate(offset int, initial int, total int) {
	redis.Set(c.OffsetKey, offset)
	redis.Set(c.InitialKey, initial)
	redis.Set(c.TotalKey, total)
}

func Refill(finished chan int) (int, error) {
	offset, err := redis.Get(c.OffsetKey).Int()
	if err != nil {
		finished <- c.RedisErrorCode
		return 2, errors.New("unable to retrieve offset")
	}
	total, err := redis.Get(c.TotalKey).Int()
	if err != nil {
		finished <- c.RedisErrorCode
		return 2, errors.New("unable to retrieve total")
	}
	initial, err := redis.Get(c.InitialKey).Int()
	if err != nil {
		finished <- c.RedisErrorCode
		return 2, errors.New("unable to retrieve initial number")
	}

	rawIds := utils.GenerateNewIdSet(total, offset, initial)
	alphaNumIds := utils.ConvertAlphaNumeric(rawIds)
	for _, elem := range alphaNumIds {
		redis.RPush(c.IdListKey, elem)
	}
	fmt.Println("generation done, increasing offset")
	redis.Set(c.OffsetKey, offset + 1)

	finished <- 0
	return 0, nil
}

