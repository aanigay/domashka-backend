package types

import (
	"encoding/json"
	"time"
)

type Date time.Time

var _ json.Unmarshaler = &Date{}

func (mt *Date) UnmarshalJSON(bs []byte) error {
	var s string
	err := json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}
	t, err := time.ParseInLocation("2006-01-02", s, time.UTC)
	if err != nil {
		return err
	}
	*mt = Date(t)
	return nil
}

func MonthEngToRus(month string) string {
	switch month {
	case "January":
		return "Январь"
	case "February":
		return "Февраль"
	case "March":
		return "Март"
	case "April":
		return "Апрель"
	case "May":
		return "Май"
	case "June":
		return "Июнь"
	case "July":
		return "Июль"
	case "August":
		return "Август"
	case "September":
		return "Сентябрь"
	case "October":
		return "Октябрь"
	case "November":
		return "Ноябрь"
	case "December":
		return "Декабрь"
	default:
		return ""
	}
}
