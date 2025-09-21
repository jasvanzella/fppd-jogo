//bau.go - logica e estrutura do bau
package main

import (
	"time"
)

// type mensagem ja declarada em portal, nao precisa repetir

type bau struct {
	X, Y         int
	Ativo        bool
	tempoAtivo   time.Duration
	DestX, DestY int
}

const ( // coordenadas fixas para o portal
	BauX = 10
	BauY = 5
)

// Cria um novo portal com destino aleatório válido
func NovoBau() bau {
	return bau{
		X:          BauX,
		Y:          BauY,
		Ativo:      false, // começa fechado
		tempoAtivo: 3 * time.Second, // abre/fecha a cada 3 segundos (pode ajustar)
	}
}