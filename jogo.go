// jogo.go - Funções para manipular os elementos do jogo
package main

import (
	"bufio"
	"os"
)

type Elemento struct {
	simbolo     rune
	cor         Cor
	corFundo    Cor
	tangivel    bool
	interagivel bool
}

type Posicao struct {
	X, Y int
}

type EventoJogo struct {
	Tipo string
	Data any
}

type Jogo struct {
	Mapa           [][]Elemento
	PosX, PosY     int
	UltimoVisitado Elemento
	StatusMsg      string
	Inimigos       []*Inimigo
	Portais        []portal
	Bau	           Bau
}

var (
	Personagem      = Elemento{'☺', CorCinzaEscuro, CorPadrao, true, false}
	ElementoInimigo = Elemento{'☠', CorVermelho, CorPadrao, true, true} // Renomeado
	Parede          = Elemento{'▤', CorParede, CorFundoParede, true, false}
	Vegetacao       = Elemento{'♣', CorVerde, CorPadrao, false, false}
	Vazio           = Elemento{' ', CorPadrao, CorPadrao, false, false}
	Portal          = Elemento{'O', CorAmarelo, CorPadrao, false, true}
	BauJogo 		= Elemento{'★', CorRoxo, CorPadrao, false, true}
)

func jogoNovo() Jogo {
	return Jogo{
		UltimoVisitado: Vazio,
		Inimigos:       []*Inimigo{},
	}
}

func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	inimigoIDCounter := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Parede.simbolo:
				e = Parede
			case ElementoInimigo.simbolo: // Usando a variável correta
				e = Vazio
				novoInimigo := &Inimigo{
					ID:   inimigoIDCounter,
					X:    x,
					Y:    y,
					Modo: ModoPatrulha,
				}
				jogo.Inimigos = append(jogo.Inimigos, novoInimigo)
				inimigoIDCounter++
			case Vegetacao.simbolo:
				e = Vegetacao
			case Portal.simbolo:
				portal := NovoPortal()
				jogo.Portais = append(jogo.Portais, portal)
				e = Vazio
			case BauJogo.simbolo:
				jogo.Bau = NovoBau(x, y) 
				e = Vazio // bau não é um elemento do mapa inicialmente
			case Personagem.simbolo:
				jogo.PosX, jogo.PosY = x, y
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}

	for i := range jogo.Portais {
		p := &jogo.Portais[i]
		jogo.Mapa[p.Y][p.X] = Portal
		elemDestino := Elemento{'╬', CorAmarelo, CorPadrao, false, true}
		jogo.Mapa[p.DestY][p.DestX] = elemDestino
	}
	return scanner.Err()
}

func jogoPodeMoverPara(jogo *Jogo, x, y int, ehInimigo bool) bool {
	if y < 0 || y >= len(jogo.Mapa) || x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	if ehInimigo {
		if jogo.Mapa[y][x].interagivel {
			return false
		}
		if x == jogo.PosX && y == jogo.PosY {
			interfaceFinalizar()
			println("Fim de Jogo! O inimigo te pegou.")
			os.Exit(0)
		}
	}

	if !ehInimigo {
		for _, inimigo := range jogo.Inimigos {
			if x == inimigo.X && y == inimigo.Y {
				interfaceFinalizar()
				println("Fim de Jogo! Você bateu em um inimigo.")
				os.Exit(0)
			}
		}
	}

	return !jogo.Mapa[y][x].tangivel
}

func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy
	elemento := jogo.Mapa[y][x]
	jogo.Mapa[y][x] = jogo.UltimoVisitado
	jogo.UltimoVisitado = jogo.Mapa[ny][nx]
	jogo.Mapa[ny][nx] = elemento
}