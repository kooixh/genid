package app

import (
	"errors"
	"fmt"
	c "github.com/kooixh/genid/pkg/constants"
	g "github.com/kooixh/genid/pkg/generators"
	"github.com/kooixh/genid/utils"
	"github.com/kooixh/genid/utils/redis"
)

const coreAppVersion = "0.1"

type CalibrationSettings struct {
	Initial int
	Offset int
	Total int
	IdType string
}

type CoreInfo struct {
	Settings *CalibrationSettings
	Status string
	Version string
	NextId string
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
	if ids.Err() != nil {
		idChannel <- ""
		refillChannel <- 0
		return
	}
	idChannel <- ids.Val()

	refillThreshold, err := redis.Get(c.IdListKey).Int()
	if err != nil {
		refillThreshold = 10
	}
	if int(redis.Length(c.IdListKey).Val()) < refillThreshold {
		go Refill(refillChannel)
	} else {
		refillChannel <- 0
	}
}

func Calibrate(cal CalibrationSettings, refillChannel chan int) {
	redis.Set(c.OffsetKey, cal.Offset)
	redis.Set(c.InitialKey, cal.Initial)
	redis.Set(c.TotalKey, cal.Total)
	redis.Set(c.GeneratorTypeKey, cal.IdType)
	redis.Set(c.CalibratedKey, true)
	go Refill(refillChannel)

}

func Refill(refillChannel chan int) {

	settings, err := retrieveCalibrationSettings()
	if err != nil {
		refillChannel <- c.RedisErrorCode
		return
	}
	rawIds := utils.GenerateNewIdSet(settings.Total, settings.Offset, settings.Initial)

	var generator g.Generator
	if settings.IdType == c.GeneratorTypeAlphaNum {
		generator = new(g.AlphaNumericIdGenerator)
	} else if settings.IdType == c.GeneratorTypeNum {
		generator = new(g.NumericIdGenerator)
	} else {
		fmt.Println("error, generator type is unknown")
		refillChannel <- c.RedisErrorUnknownTypeCode
		return
	}

	generatedIds := generator.Generate(rawIds)
	for _, elem := range generatedIds {
		redis.RPush(c.IdListKey, elem)
	}
	redis.Set(c.OffsetKey, settings.Offset + 1)
	refillChannel <- 0

}

func ReturnAppInfo() (*CoreInfo, error) {
	var info *CoreInfo
	info = new(CoreInfo)
	status := redis.Get(c.CalibratedKey)
	var statValue string
	if status.Err() != nil {
		statValue = c.NotCalibrated
	} else {
		statValue = c.Calibrated
	}
	info.Status = statValue
	info.Version = coreAppVersion
	if statValue == c.NotCalibrated {
		return info, nil
	}
	settings, err := retrieveCalibrationSettings()
	if err != nil {
		return nil, errors.New("unable to retrieve calibration settings from Redis")
	}
	info.Settings = settings
	nextId := redis.Peek(c.IdListKey)
	if nextId.Err() != nil {
		info.NextId = ""
	} else {
		info.NextId = nextId.Val()[0]
	}
	return info, nil
}

func retrieveCalibrationSettings() (*CalibrationSettings, error){
	offset, errOffset := redis.Get(c.OffsetKey).Int()
	total, errTotal := redis.Get(c.TotalKey).Int()
	initial, errInitial := redis.Get(c.InitialKey).Int()
	generatorType := redis.Get(c.GeneratorTypeKey).Val()
	if errOffset != nil || errTotal != nil || errInitial != nil {
		return nil, errors.New("unable to retrieve calibration settings from Redis")
	}
	return &CalibrationSettings{
		Initial: initial,
		Offset: offset,
		Total: total,
		IdType: generatorType,
	}, nil
}