package util

import (
	"fmt"
	"math/rand"
)

func GeneratePIN() string {
	return fmt.Sprintf("%04d", rand.Intn(10000))
}
