package api

import (
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/St-Ivanov/todo-list-project/internal/models"
)

var (
	errIncFormatDStart = errors.New("incorrect format dstart.")
	errIncFormatRepeat = errors.New("incorrect format repeat.")
	errParseTime       = errors.New("error parse time.")
	errIncNumOfPar     = errors.New("incorrect number of parameters.")
	errParseStrToInt   = errors.New("error parse string to int.")
	errDowloadPage     = errors.New("error dowload the page")
)

// Function for finding the next date according to a specified rule
func nextDate(now time.Time, dstart string, repeat string) (string, error) {
	if len(dstart) != 8 {
		return "", errIncFormatDStart
	}
	if len(repeat) == 0 {
		return "", errIncFormatRepeat
	}

	dataSplit := strings.Split(repeat, " ")

	date, err := time.Parse(models.DataFormat, dstart)
	if err != nil {
		return "", errParseTime
	}

	switch dataSplit[0] {
	case "y":
		if len(dataSplit) != 1 {
			return "", errIncNumOfPar
		}
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}
		}
	case "d":
		if len(dataSplit) != 2 {
			return "", errIncNumOfPar
		}
		step, err := strconv.Atoi(dataSplit[1])
		if err != nil || step > 400 {
			return "", errParseStrToInt
		}
		for {
			date = date.AddDate(0, 0, step)
			if date.After(now) {
				break
			}
		}
	case "w":
		if len(dataSplit) != 2 {
			return "", errIncFormatRepeat
		}
		daysString := strings.Split(dataSplit[1], ",")
		days := make([]int, len(daysString))
		for i, el := range daysString {
			num, err := strconv.Atoi(el)
			if num <= 0 || num > 7 || err != nil {
				return "", errIncFormatRepeat
			}
			days[i] = num
		}
		for {
			date = date.AddDate(0, 0, 1)
			day := date.Weekday()
			if day == 0 {
				day = 7
			}

			if date.After(now) && slices.Contains(days, int(day)) {
				break
			}
		}
	case "m":
		if len(dataSplit) != 2 && len(dataSplit) != 3 {
			return "", errIncFormatRepeat
		}
		daysString := strings.Split(dataSplit[1], ",")
		days := make([]int, len(daysString))
		for i, m := range daysString {
			temp, err := strconv.Atoi(m)
			if err != nil || temp > 31 || temp == 0 || temp < -2 {
				return "", errIncFormatRepeat
			}
			days[i] = temp
		}
		var (
			monthsString []string
		)
		if len(dataSplit) == 3 {
			monthsString = strings.Split(dataSplit[2], ",")
		}
		months := make([]int, len(monthsString))
		for i, m := range monthsString {
			temp, err := strconv.Atoi(m)
			if err != nil || temp < 1 || temp > 12 {
				return "", errIncFormatRepeat
			}
			months[i] = temp
		}

		for {
			date = date.AddDate(0, 0, 1)

			flag := false
			for _, m := range days {
				if m > 0 && date.Day() == m {
					flag = true
					break
				}
				if m < 0 {
					tempDate := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
					tempDate = tempDate.AddDate(0, 1, m)
					if tempDate.Day() == date.Day() {
						flag = true
						break
					}
				}
			}

			if date.After(now) && slices.Contains(months, int(date.Month())) && flag {
				break
			}
		}
	default:
		return "", errIncFormatRepeat
	}
	return date.Format(models.DataFormat), nil
}

// A handler for nextdate
func handlerNextDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "error loading the page.", http.StatusInternalServerError)
		return
	}

	nowString := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var (
		now time.Time
		err error
	)

	if len(nowString) == 0 {
		now = time.Now().UTC()
	} else {
		now, err = time.Parse(models.DataFormat, nowString)
		if err != nil {
			http.Error(w, errDowloadPage.Error(), http.StatusBadRequest)
			return
		}
	}

	ans, err := nextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, errDowloadPage.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ans))
}
