package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDate(now time.Time, dateStr, repeat string) (string, error) {
	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", err
	}

	// обрезаем время
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// без повтора
	if repeat == "" {
		if !date.After(now) {
			return "", nil
		}
		return date.Format("20060102"), nil
	}

	if repeat == "y" {
		date = date.AddDate(1, 0, 0) // всегда прибавляем хотя бы 1 раз
		for !date.After(now) {
			date = date.AddDate(1, 0, 0)
		}
		return date.Format("20060102"), nil
	}

	if strings.HasPrefix(repeat, "d ") {
		parts := strings.Fields(repeat)
		if len(parts) != 2 {
			return "", nil
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", nil
		}
		date = date.AddDate(0, 0, days) // всегда сдвиг
		for !date.After(now) {
			date = date.AddDate(0, 0, days)
		}
		return date.Format("20060102"), nil
	}

	return "", errors.New("неподдерживаемый repeat")
}

// Переаботал и упростил функцию afterNow
func afterNow(date, now time.Time) bool {
	return date.After(now)
}

// HTTP обработчик
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")
	nowStr := r.FormValue("now")

	var now time.Time
	var err error

	if nowStr != "" {
		now, err = time.Parse("20060102", nowStr)
		if err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "неверный параметр now"})
			return
		}
	} else {
		now = time.Now()
	}

	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if next == "" {
		writeJson(w, http.StatusOK, map[string]string{"next": ""})
		return
	}

	writeJson(w, http.StatusOK, map[string]string{"next": next})
}
