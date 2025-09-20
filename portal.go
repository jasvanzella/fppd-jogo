package main

import (
	"time"
)
type Mensagem struct {
	Tipo string
	PosX int
	PosY int
}

type portal struct {
	X, Y         int
	Ativo        bool
	tempoAtivo   time.Duration
	DestX, DestY int
}

const ( // coordenadas fixas para o portal
	PortalX = 10
	PortalY = 5
	DestX   = 60
	DestY   = 2
)

// Cria um novo portal com destino aleatório válido
func NovoPortal() portal {
	return portal{
		X:          PortalX,
		Y:          PortalY,
		Ativo:      false, // começa fechado
		tempoAtivo: 3 * time.Second, // abre/fecha a cada 3 segundos (pode ajustar)
		DestX:      DestX,
		DestY:      DestY,
	}
}