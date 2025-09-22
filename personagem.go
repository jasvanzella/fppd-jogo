// personagem.go - Funções para movimentação e ações do personagem
package main

import (
	"os"
)

func personagemMover(tecla rune, jogo *Jogo) {
	dx, dy := 0, 0
	switch tecla {
	case 'w':
		dy = -1
	case 'a':
		dx = -1
	case 's':
		dy = 1
	case 'd':
		dx = 1
	}

	nx, ny := jogo.PosX+dx, jogo.PosY+dy
	// adicionado o terceiro argumento 'false', indicando que quem está movendo não é um inimigo.
	if jogoPodeMoverPara(jogo, nx, ny, false) {
		jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
		jogo.PosX, jogo.PosY = nx, ny
	}
}

// agora recebe 'canalJogo' para poder enviar eventos.
func personagemInteragir(jogo *Jogo) {

	for i := range jogo.Portais {
		p := &jogo.Portais[i]
		if p.Ativo && jogo.PosX == p.X && jogo.PosY == p.Y {
			// Teletransporta o jogador para o destino do portal
			jogo.PosX = p.DestX
			jogo.PosY = p.DestY
			jogo.StatusMsg = "Teleportado via portal!"
			return // Importante para sair da função após interagir
		}
	}

	// Interação com o baú
	if jogo.Bau.Ativo && jogo.PosX == jogo.Bau.X && jogo.PosY == jogo.Bau.Y {
		jogo.StatusMsg = "Você encontrou o baú!"
		interfaceFinalizar()
		println("Parabéns! Você encontrou o baú. Fim de jogo.")
		os.Exit(0) // termina o jogo aqui
	}

	jogo.StatusMsg = "Nada para interagir aqui."
}

// a função foi ajustada para receber 'canalJogo'.
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
	switch ev.Tipo {
	case "sair":
		return false
	case "interagir":
		personagemInteragir(jogo)
	case "mover":
		personagemMover(ev.Tecla, jogo)
	}
	return true
}