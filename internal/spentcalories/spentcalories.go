package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	errIncorrectInput      = errors.New("некорректный ввод данных о тренировке")
	errDurationLessOne     = errors.New("длительность тренировки должна быть больше 0")
	errStepsLessOne        = errors.New("количество шагов должно быть больше 0")
	errUnknownTrainingType = errors.New("неизвестный тип тренировки")
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	input := strings.Split(data, ",")
	if len(input) != 3 {
		return 0, "", 0, errIncorrectInput
	}

	steps, err := strconv.Atoi(input[0])
	if err != nil {
		return 0, "", 0, err
	}
	if !(steps > 0) {
		return 0, "", 0, errStepsLessOne
	}

	duration, err := time.ParseDuration(input[2])
	if err != nil {
		return 0, "", 0, err
	}
	if !(duration > 0) {
		return 0, "", 0, errDurationLessOne
	}

	return steps, input[1], duration, nil
}

func distance(steps int, height float64) float64 {
	// Длина шага
	stepDistance := height * stepLengthCoefficient
	// Дистанция в метрах
	distance := float64(steps) * stepDistance
	// Дистанция в километрах
	return distance / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if !(duration > 0) {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if !(steps > 0 && weight > 0 && height > 0 && duration > 0) {
		return 0, errIncorrectInput
	}
	meanSpeed := meanSpeed(steps, height, duration)
	return (weight * meanSpeed * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if !(steps > 0 && weight > 0 && height > 0 && duration > 0) {
		return 0, errIncorrectInput
	}
	meanSpeed := meanSpeed(steps, height, duration)
	return ((weight * meanSpeed * duration.Minutes()) / minInH) * walkingCaloriesCoefficient, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64
	distance := distance(steps, height)
	meanSpeed := meanSpeed(steps, height, duration)
	switch trainingType {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", errUnknownTrainingType
	}
	if err != nil {
		log.Println(err)
		return "", err
	}
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", trainingType, float64(duration.Hours()), distance, meanSpeed, calories), err
}
