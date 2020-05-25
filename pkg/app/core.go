package app

import (
	"encoding/json"
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

func GenerateNewId(refillChannel chan int) *string {
	ids := redis.LPop(c.IdListKey)
	if ids.Err() != nil {
		refillChannel <- 0
		return nil
	}
	id := ids.Val()
	go checkForRefill(refillChannel)
	return &id
}

func checkForRefill(refillChannel chan int) {
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
	bytes, err := json.Marshal(cal)
	if err != nil {
		refillChannel <- c.CalibrationErrorJsonMarshall
	}
	redis.Set(c.CalibratedSettings, string(bytes))
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

func CheckIsCalibrated() bool {
	return redis.Get(c.CalibratedKey).Err() == nil
}

func ReturnAppInfo() (*CoreInfo, error) {
	var info = new(CoreInfo)
	status := redis.Get(c.CalibratedKey)
	if status.Err() != nil {
		info.Status = c.NotCalibrated
	} else {
		info.Status = c.Calibrated
	}
	info.Version = coreAppVersion
	if info.Status == c.NotCalibrated {
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

func retrieveCalibrationSettings() (*CalibrationSettings, error) {
	var settings CalibrationSettings
	settingsJson := redis.Get(c.CalibratedSettings)
	if settingsJson.Err() != nil {
		return nil, errors.New("unable to retrieve calibration settings from Redis")
	}
	bytes, bytesErr := settingsJson.Bytes()
	if bytesErr != nil {
		return nil, errors.New("unable to retrieve bytes[] settings from Redis")
	}
	err := json.Unmarshal(bytes, &settings)
	if err != nil {
		return nil, errors.New("unable to Marshal settings")
	}
	return &settings, nil
}