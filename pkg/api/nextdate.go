package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat rule is empty")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid dstart date format: %w", err)
	}

	parts := strings.SplitN(repeat, " ", 2)
	rule := parts[0]

	switch rule {
	case "y":
		for !afterNow(date, now) {
			date = date.AddDate(1, 0, 0)
		}
		return date.Format(DateFormat), nil

	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("repeat d rule missing interval")
		}
		daysStr := strings.TrimSpace(parts[1])
		interval, err := strconv.Atoi(daysStr)
		if err != nil {
			return "", fmt.Errorf("invalid interval for d rule: %w", err)
		}
		if interval < 1 || interval > 400 {
			return "", fmt.Errorf("d interval must be between 1 and 400")
		}
		for !afterNow(date, now) {
			date = date.AddDate(0, 0, interval)
		}
		return date.Format(DateFormat), nil

	case "w":

		return "", fmt.Errorf("unsupported repeat rule: %s", repeat)

	default:
		return "", fmt.Errorf("unsupported repeat rule: %s", repeat)
	}
}

func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()

	if y1 > y2 {
		return true
	} else if y1 == y2 {
		if m1 > m2 {
			return true
		} else if m1 == m2 {
			return d1 > d2
		}
	}
	return false
}

// HTTP обработчик
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	if nowStr == "" {
		nowStr = time.Now().Format(DateFormat)
	}
	now, err := time.Parse(DateFormat, nowStr)
	if err != nil {
		http.Error(w, "invalid now parameter: "+err.Error(), http.StatusBadRequest)
		return
	}

	date := r.FormValue("date")
	if date == "" {
		http.Error(w, "missing date parameter", http.StatusBadRequest)
		return
	}

	repeat := r.FormValue("repeat")
	if repeat == "" {
		http.Error(w, "missing repeat parameter", http.StatusBadRequest)
		return
	}

	next, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, next)
}

//
