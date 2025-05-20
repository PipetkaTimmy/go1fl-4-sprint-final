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

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")

	if len(parts) != 2 {
		return 0, 0, errors.New("недопустимый формат данных")
	}

	if parts[0] == "" || parts[1] == "" {
		return 0, 0, errors.New("пустое значение шага или продолжительности")
	}

	if strings.TrimSpace(parts[0]) != parts[0] {
		return 0, 0, errors.New("пробелы в значении шагов")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, errors.New("недопустимое значение шагов")
	}

	if steps <= 0 {
		return 0, 0, errors.New("шаги должны быть больше нуля")
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, errors.New("неверный формат продолжительности")
	}

	if duration <= 0 {
		return 0, 0, errors.New("продолжительность должна быть больше нуля")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	if weight <= 0 || height <= 0 {
		log.Printf("вес и рост должны быть больше нуля")
		return ""
	}

	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Printf("%v", err)
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Printf("%v", err)
		return ""
	}

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, calories)
	return result
}
