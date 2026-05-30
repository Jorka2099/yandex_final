package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	
	return date.After(now)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if repeat == "" {
		return "", fmt.Errorf("repeat is empty")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("wrong dstart format %v", err)
	}

	repeatSplitted := strings.Split(repeat, " ")
	rule := repeatSplitted[0]

	switch rule {
	case "d":
		if len(repeatSplitted) < 2 {
			return "", fmt.Errorf("days amount not mentioned")
		}

		daysAmount, err := strconv.Atoi(repeatSplitted[1])
		if err != nil {
			return "", fmt.Errorf("invalid days format %v", err)
		}

		if daysAmount <= 0 {
			return "", fmt.Errorf("days amount must be positive")
		}

		if daysAmount > 400 {
			return "", fmt.Errorf("Days amount exceed 400")
		}

		for {
			date = date.AddDate(0, 0, daysAmount)
			if afterNow(date, now)	 {
				break
			}
		}

		return date.Format(DateFormat), nil

	case "y":

		if len(repeatSplitted) > 1 {
			return "", fmt.Errorf("More than one parameter")
		}

		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	default:
		return "", fmt.Errorf("unsupported rule")
	}

}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse(DateFormat, nowStr)
	if err != nil {
		http.Error(w, "wrong data format", http.StatusBadRequest)
		return
	}

	res, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(res))
}
