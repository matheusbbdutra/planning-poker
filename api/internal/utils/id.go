package utils

import (
	gonanoid "github.com/jaevor/go-nanoid"
)

func GenerateID() (string, error) {
	// Usar `Standard` garante uma distribuição uniforme e é mais seguro contra colisões.
	// O tamanho 10 é mantido, mas considere usar 21 (Canonic) para maior segurança.
	id, err := gonanoid.Standard(10)
	if err != nil {
		return "", err
	}
	return id(), nil
}
