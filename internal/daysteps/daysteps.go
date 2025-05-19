package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	stepLength = 0.65
	mInKm      = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")

	if len(parts) != 2 {
		return 0, 0, errors.New("недопустимый формат данных")
	}

	if parts[0] == "" || parts[1] == "" {
		return 0, 0, errors.New("пустое значение шага или продолжительности")
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

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {

	if weight <= 0 || height <= 0 {
		fmt.Println("вес и рост должны быть больше нуля")
		return ""
	}

	steps, duration, err := parsePackage(data)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.", steps, distanceKm, calories)
	return result
}
