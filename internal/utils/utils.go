package utils

import (
	"errors"
	"strconv"
	"time"
)

func GetDataForTicker(value string) (interval int, duration time.Duration, err error) {
	strDuration := value[len(value)-1:]
	strInterval := value[0 : len(value)-1]

	if interval, err = strconv.Atoi(strInterval); err != nil {
		return
	}

	switch strDuration {
	case "s":
		duration = time.Second
	case "m":
		duration = time.Minute
	case "h":
		duration = time.Hour
	default:
		err = errors.New("ошибка создания тикера")
	}

	return
}
