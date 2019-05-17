package random_test

import (
	"testing"

	random "github.com/GameComponent/economy-service/pkg/helper/random"
)

func TestGenerateRandomInt(t *testing.T) {
	totalCount := map[int64]int{
		0: 0,
		1: 0,
		2: 0,
		3: 0,
	}

	for i := 0; i < 1000; i++ {
		randomResult := random.GenerateRandomInt(0, 3)
		totalCount[randomResult] += 1
	}

	if totalCount[0] == 0 {
		t.Errorf("totalCount[0] is zero")
	}

	if totalCount[1] == 0 {
		t.Errorf("totalCount[1] is zero")
	}

	if totalCount[2] == 0 {
		t.Errorf("totalCount[2] is zero")
	}

	if totalCount[3] == 0 {
		t.Errorf("totalCount[3] is zero")
	}
}
func TestGenerateRandomInt2(t *testing.T) {
	totalCount := map[int64]int{
		0: 0,
		1: 0,
		2: 0,
		3: 0,
	}

	for i := 0; i < 1000; i++ {
		randomResult := random.GenerateRandomInt(1, 3)
		totalCount[randomResult] += 1
	}

	if totalCount[0] != 0 {
		t.Errorf("totalCount[0] is not zero")
	}

	if totalCount[1] == 0 {
		t.Errorf("totalCount[1] is zero")
	}

	if totalCount[2] == 0 {
		t.Errorf("totalCount[2] is zero")
	}

	if totalCount[3] == 0 {
		t.Errorf("totalCount[3] is zero")
	}
}

func TestGenerateRandomIntLowerMaxShouldReturnMin(t *testing.T) {
	result := random.GenerateRandomInt(3, 1)

	if result != 3 {
		t.Errorf("result should be 3")
	}
}

func TestGenerateRandomIntEqualMinMaxShouldReturnMin(t *testing.T) {
	result := random.GenerateRandomInt(3, 3)

	if result != 3 {
		t.Errorf("result should be 3")
	}
}
