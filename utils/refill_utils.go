package utils

import (
	"math/rand"
	"time"
)

func GenerateNewIdSet(totalToGenerate int, offset int, initial int) []int64 {
	var inputs []int64
	for i := 1; i <= totalToGenerate; i++ {
		inputs = append(inputs, int64((offset * initial) + i))
	}
	return inputs
}

func Shuffle(inputs []string) []string {
	rand.Seed(time.Now().Unix())
	for i := len(inputs) - 1; i > 0; i-- {
		j := rand.Intn(i)
		swap(inputs, i , j)
	}
	return inputs
}

func swap(inputs []string, i int, j int) {
	temp := inputs[i]
	inputs[i] = inputs[j]
	inputs[j] = temp
}