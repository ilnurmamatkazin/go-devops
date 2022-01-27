package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
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
		err = errors.New("не определен тип duration")
	}

	return
}

func SetHesh(id, metricType, key string, delta *int64, value *float64) []byte {
	if key == "" {
		return nil
	}

	var hash []byte
	if metricType == "gauge" {
		hash = []byte(fmt.Sprintf("%s:gauge:%f", id, *value))
	} else {
		hash = []byte(fmt.Sprintf("%s:counter:%d", id, *delta))
	}

	h := hmac.New(sha256.New, []byte(key))
	h.Write(hash)
	return h.Sum(nil)

}

func SetEncodeHesh(id, metricType, key string, delta *int64, value *float64) string {
	return hex.EncodeToString(SetHesh(id, metricType, key, delta, value))
}
