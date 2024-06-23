package engine

import (
	"crypto/rand"
	"fmt"
	mr "math/rand"
	"strings"
)

var (
	idChars = strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "")
	idCap   = len(idChars) - 1
)

func randomId() string {
	result := strings.Builder{}
	bytes := make([]byte, 12)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(fmt.Errorf("couldn't create random: %v", err))
	}
	for i := 0; i < len(bytes); i++ {
		index := int(bytes[i]) % idCap
		result.WriteString(idChars[index])
	}
	return result.String()
}

type IFunc func() int

func irandom(lower, upper int) int {
	return lower + int(float64(upper-lower)*mr.Float64())
}

func irandomF(lower, upper int) IFunc {
	return func() int { return irandom(lower, upper) }
}

type FFunc func() float64

func frandom(lower, upper float64) float64 {
	return lower + (upper-lower)*mr.Float64()
}

func frandomF(lower, upper float64) FFunc {
	return func() float64 { return frandom(lower, upper) }
}
