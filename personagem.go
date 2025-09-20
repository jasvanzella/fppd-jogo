// personagem.go - Funções para movimentação e ações do personagem
package main

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
	if jogoPodeMoverPara(jogo, nx, ny, false) {
		jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
		jogo.PosX, jogo.PosY = nx, ny
	}
}

func personagemInteragir(jogo *Jogo) {
	for i := range jogo.Portais {
		p := &jogo.Portais[i]
		if p.Ativo && jogo.PosX == p.X && jogo.PosY == p.Y {
			jogo.PosX = p.DestX
			jogo.PosY = p.DestY
			jogo.StatusMsg = "Teleportado via portal!"
			return
		}
	}
	
	// NOVO: Lógica de interação com o baú.
	if jogo.Bau.Visivel && jogo.PosX == jogo.Bau.X && jogo.PosY == jogo.Bau.Y {
		jogo.StatusMsg = "Você encontrou o baú!"
		return
	}

	jogo.StatusMsg = "Nada para interagir aqui."
}

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