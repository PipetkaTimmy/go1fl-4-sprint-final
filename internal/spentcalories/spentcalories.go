package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
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
	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, errors.New("недопустимый формат данных")
	}

	if parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return 0, "", 0, errors.New("пустое значение шага, активности или продолжительности")
	}

	if strings.TrimSpace(parts[0]) != parts[0] {
		return 0, "", 0, errors.New("пробелы в значении шагов")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, errors.New("недопустимое значение шагов")
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("шаги должны быть больше нуля")
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, errors.New("неверный формат продолжительности")
	}

	if duration <= 0 {
		return 0, "", 0, errors.New("продолжительность должна быть больше нуля")
	}

	return steps, parts[1], duration, nil
}

func distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 {
		return 0.0
	}

	stepLength := height * stepLengthCoefficient
	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm
	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0.0
	}

	dist := distance(steps, height)
	durationInHours := duration.Hours()
	speed := dist / durationInHours

	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var dist, speed, calories float64

	switch activity {
	case "Ходьба":
		dist = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Бег":
		dist = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	durationHours := duration.Hours()
	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity, durationHours, dist, speed, calories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("шаги должны быть больше нуля")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть больше нуля")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть больше нуля")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше нуля")
	}

	speed := meanSpeed(steps, height, duration)

	if speed <= 0 {
		return 0, errors.New("скорость должна быть больше нуля")
	}

	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("шаги должны быть больше нуля")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть больше нуля")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть больше нуля")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше нуля")
	}

	speed := meanSpeed(steps, height, duration)

	if speed <= 0 {
		return 0, errors.New("скорость должна быть больше нуля")
	}

	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / minInH * walkingCaloriesCoefficient

	return calories, nil
}
