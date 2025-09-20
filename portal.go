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
	canalMapa    chan Mensagem
	tempoAtivo   time.Duration
	DestX, DestY int
}

const (	// coordenadas fixas para o portal 
    PortalX = 10
    PortalY = 5
    DestX   = 20
    DestY   = 15
)

// Cria um novo portal com destino aleatório válido
func NovoPortal() portal {
    return portal{
        X:          PortalX,
        Y:          PortalY,
        Ativo:      false,                // começa fechado
        canalMapa:  make(chan Mensagem, 1),
        tempoAtivo: 3 * time.Second,      // abre/fecha a cada 3 segundos (pode ajustar)
        DestX:      DestX,
        DestY:      DestY,
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