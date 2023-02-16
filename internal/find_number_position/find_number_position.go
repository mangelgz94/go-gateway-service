package find_number_position

import (
	"context"
	"sort"
)

type FindNumberPositionService struct {
	config  *Config
	numbers []int
}

func (f *FindNumberPositionService) FindNumberPosition(ctx context.Context, number int) int {
	//Go has already a binary search function for a slice of ints
	index := sort.SearchInts(f.numbers, number)
	if index == len(f.numbers) {
		return -1
	}

	return index
}

func (f *FindNumberPositionService) fillNumbers() {
	numbers := make([]int, 0, f.config.ArraySize)
	for i := 0; i < f.config.ArraySize; i++ {
		numbers = append(numbers, i)
	}

	f.numbers = numbers
}

func New(config *Config) *FindNumberPositionService {
	findNUmberService := &FindNumberPositionService{config: config}
	findNUmberService.fillNumbers()

	return findNUmberService
}
