package main

import (
	"math/rand"
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
	canalMapa    chan Mensagem
	tempoAtivo   time.Duration
	DestX, DestY int
}

// Cria um novo portal com destino aleatório válido
func NovoPortal(x, y int, jogo *Jogo) portal {
	min := 15
	max := 30
	destX, destY := encontrarDestinoValido(jogo)
	return portal{
		X:          x,
		Y:          y,
		Ativo:      true,
		canalMapa:  make(chan Mensagem, 1),
		tempoAtivo: time.Duration(rand.Intn(max-min+1)+min) * time.Second,
		DestX:      destX,
		DestY:      destY,
	}
}

// Inicializa todos os portais do mapa
func inicializarPortais(jogo *Jogo) {
	for i := range jogo.Portais {
		p := &jogo.Portais[i]
		go rotinaPortal(jogo, p)
	}
}

// Rotina de ativação/desativação do portal e teleporte
func rotinaPortal(jogo *Jogo, p *portal) {
	elemDestino := Elemento{'╬', CorAmarelo, CorPadrao, false, true}
	for {
		select {
		case msg := <-p.canalMapa:
			if msg.Tipo == "Teleporte!" {
				jogo.PosX = p.DestX
				jogo.PosY = p.DestY
				jogo.StatusMsg = "Teleportado via portal!"
			}
		case <-time.After(p.tempoAtivo):
			if p.Ativo {
				// desativa portal
				jogo.Mapa[p.Y][p.X] = Vazio
				jogo.Mapa[p.DestY][p.DestX] = Vazio
				p.Ativo = false
			} else {
				// reativa portal
				jogo.Mapa[p.Y][p.X] = Portal
				jogo.Mapa[p.DestY][p.DestX] = elemDestino
				p.Ativo = true
			}
			interfaceDesenharJogo(jogo)
		}
	}
}

// Encontra posição válida para o destino do portal
func encontrarDestinoValido(jogo *Jogo) (int, int) {
	for {
		y := rand.Intn(len(jogo.Mapa))
		x := rand.Intn(len(jogo.Mapa[y]))
		if !jogo.Mapa[y][x].tangivel {
			return x, y
		}
	}
}