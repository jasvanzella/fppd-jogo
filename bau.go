// bau.go - logica do bau
package main

import (
	"time"
)

type Bau struct {
	X, Y    int
	Ativo bool
	tempoAtivo   time.Duration
}

// NovoBau cria e retorna uma nova instância do Baú com as coordenadas fornecidas.
func NovoBau(x, y int) Bau {
	return Bau{
		X:       x,
		Y:       y,
		Ativo:      false, // começa fechado
		tempoAtivo: 3 * time.Second, // abre/fecha a cada 3 segundos (pode ajustar)
	}
}