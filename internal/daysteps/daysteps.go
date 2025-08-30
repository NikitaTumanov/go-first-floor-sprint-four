package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

var (
	errIncorrectInput  = errors.New("некорректный ввод данных о прогулке")
	errDurationLessOne = errors.New("длительность тренировки должна быть больше 0")
	errStepsLessOne    = errors.New("количество шагов должно быть больше 0")
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	input := strings.Split(data, ",")
	if len(input) != 2 {
		return 0, 0, errIncorrectInput
	}

	steps, err := strconv.Atoi(input[0])
	if err != nil {
		return 0, 0, err
	}
	if !(steps > 0) {
		return 0, 0, errStepsLessOne
	}

	duration, err := time.ParseDuration(input[1])
	if err != nil {
		return 0, 0, err
	}
	if !(duration > 0) {
		return 0, 0, errDurationLessOne
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	if !(steps > 0) {
		log.Println(err)
		return ""
	}

	// Дистанция в метрах
	distance := float64(steps) * stepLength
	// Дистанция в километрах
	distance /= mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calories)
}
