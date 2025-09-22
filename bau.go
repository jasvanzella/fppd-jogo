// bau.go - logica do bau
//segue a mesma logica do portal
package main

import (
	"time"
)

type Bau struct {
	X, Y    int
	Ativo bool
	tempoAtivo   time.Duration
}

// cria e retorna uma nova instância do bau 
func NovoBau(x, y int) Bau {
	return Bau{
		X:       x,
		Y:       y,
		Ativo:      false, // começa fechado
		tempoAtivo: 3 * time.Second, // abre/fecha a cada 3 segundos (pode ajustar)
	}
}