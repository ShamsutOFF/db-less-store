package auth

import (
	"math/rand/v2"
)

type CodeGenerator struct{}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{}
}

func (g *CodeGenerator) GenerateCode() string {
	// Генерируем 4-значный код
	return string(randomDigit()) + string(randomDigit()) + string(randomDigit()) + string(randomDigit())
}

func randomDigit() byte {
	return byte(rand.IntN(10)) + '0'
}
