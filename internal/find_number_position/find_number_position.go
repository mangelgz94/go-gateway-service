package find_number_position

import (
	"context"
)

type FindNumberPositionService struct {
	config  *Config
	numbers []int
}

func (f *FindNumberPositionService) FindNumberPosition(ctx context.Context, number int) int {
	numberPosition := -1
	start := 0
	end := len(f.numbers) - 1
	for start <= end {
		mid := (start + end) / 2
		if f.numbers[mid] == number {
			numberPosition = mid
			break
		} else if f.numbers[mid] < number {
			start = mid + 1
		} else if f.numbers[mid] > number {
			end = mid - 1
		}
	}

	return numberPosition
}

func (f *FindNumberPositionService) fillNumbers() {
	numbers := make([]int, 0, f.config.ArraySize)
	for i := 0; i < f.config.ArraySize; i++ {
		numbers = append(numbers, i)
	}

	f.numbers = numbers
}

func New(config *Config) *FindNumberPositionService {
	findNumberService := &FindNumberPositionService{config: config}
	findNumberService.fillNumbers()

	return findNumberService
}
