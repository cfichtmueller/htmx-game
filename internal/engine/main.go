package engine

import (
	"crypto/rand"
	"fmt"
	"strings"
)

var (
	idChars = strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "")
	idCap   = len(idChars) - 1
)

func intersects(p *Player, cells []*Cell) (*Cell, bool) {
	for _, c := range cells {
		if p.Agent.Intersects(c.Agent) {
			return c, true
		}
	}
	return nil, false
}

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
