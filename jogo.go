// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"math/rand"
	"os"
	"time"
)

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	simbolo     rune
	cor         Cor
	corFundo    Cor
	tangivel    bool // Indica se o elemento bloqueia passagem
	interagivel bool // Indica se o elemento pode ser interagido
}

type Posicao struct {
	X, Y int
}

type Movimento struct {
	DeX, DeY     int       // posição atual
	ParaX, ParaY int       // nova posição
	Result       chan bool // opcional: confirma se a movimentação deu certo
}

type EventoJogo struct {	//novo
    Tipo  string      // "Movimento", "Portal", "Status"
    Data any // Pode carregar Movimento, Mensagem ou string
}

var canalJogo = make(chan EventoJogo, 10) //novo

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa           [][]Elemento // grade 2D representando o mapa
	PosX, PosY     int          // posição atual do personagem
	UltimoVisitado Elemento     // elemento que estava na posição do personagem antes de mover
	StatusMsg      string       // mensagem para a barra de status
	Inimigos       []Posicao    // posições atuais dos inimigos
	Portais        []portal     // colecao para armazenar os portais ativos
	TemChave	   bool         // indica se o jogador possui a chave
}

// Elementos visuais do jogo
var (
	Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true, false}
	Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true, true}
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true, false}
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false, false}
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false, false}
	Portal     = Elemento{'O', CorAmarelo, CorPadrao, false, true}
	Chave 	   = Elemento{'🔑', CorAmarelo, CorPadrao, false, true}
)

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	return Jogo{
		UltimoVisitado: Vazio,
		Inimigos:       []Posicao{}, // garante que não tem inimigos extras
	}
}

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Parede.simbolo:
				e = Parede
			case Inimigo.simbolo:
				e = Vazio                                            // o inimigo será desenhado depois
				jogo.Inimigos = append(jogo.Inimigos, Posicao{x, y}) // adiciona a posição do inimigo
			case Vegetacao.simbolo:
				e = Vegetacao
			case Portal.simbolo:
				portal := NovoPortal()
				jogo.Portais = append(jogo.Portais, portal)
				e = Vazio // o portal será desenhado depois
			case Personagem.simbolo:
				jogo.PosX, jogo.PosY = x, y // registra a posição inicial do personagem
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++

		for i := range jogo.Portais { //novo
			p := &jogo.Portais[i]
			// desenha portal na posição inicial
			jogo.Mapa[p.Y][p.X] = Portal
			// desenha destino do portal
			elemDestino := Elemento{'╬', CorAmarelo, CorPadrao, false, true}
			jogo.Mapa[p.DestY][p.DestX] = elemDestino
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil

}

// Verifica se o personagem pode se mover para a posição (x, y)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	// Verifica se a coordenada Y está dentro dos limites verticais do mapa
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	// Verifica se a coordenada X está dentro dos limites horizontais do mapa
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	if jogo.Mapa[y][x].simbolo == Inimigo.simbolo {
		interfaceFinalizar()
		os.Exit(0)
		// termina o jogo
	}
	// Verifica se o elemento de destino é tangível (bloqueia passagem)
	if jogo.Mapa[y][x].tangivel {
		return false
	}

	// Pode mover para a posição
	return true
}

// Move um elemento para a nova posição
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy

	// Obtem elemento atual na posição
	elemento := jogo.Mapa[y][x] // guarda o conteúdo atual da posição

	jogo.Mapa[y][x] = jogo.UltimoVisitado   // restaura o conteúdo anterior
	jogo.UltimoVisitado = jogo.Mapa[ny][nx] // guarda o conteúdo atual da nova posição
	jogo.Mapa[ny][nx] = elemento            // move o elemento
}

// moverInimigo faz o inimigo se mover sozinho em intervalos de tempo
func moverInimigo(x, y int, canal chan Movimento) {
	for {
		time.Sleep(time.Second)

		dx, dy := 0, 0
		switch rand.Intn(4) {
		case 0:
			dx = 1
		case 1:
			dx = -1
		case 2:
			dy = 1
		case 3:
			dy = -1
		}
		nx, ny := x+dx, y+dy
		// envia pedido de movimentação
		result := make(chan bool)
		canal <- Movimento{
			DeX: x, DeY: y,
			ParaX: nx, ParaY: ny,
			Result: result,
		}

		// só atualiza as coordenadas se o movimento foi permitido
		if <-result {
			x, y = nx, ny
		}
	}
}

func gerenciarMapa(jogo *Jogo, canal chan Movimento) {
	for mov := range canal {
		// verifica se a posição de destino é válida
		if mov.ParaY >= 0 && mov.ParaY < len(jogo.Mapa) &&
			mov.ParaX >= 0 && mov.ParaX < len(jogo.Mapa[0]) &&
			!jogo.Mapa[mov.ParaY][mov.ParaX].tangivel {

			// move o inimigo
			jogo.Mapa[mov.DeY][mov.DeX] = Vazio
			jogo.Mapa[mov.ParaY][mov.ParaX] = Inimigo

			if mov.Result != nil {
				mov.Result <- true
			}
		} else if mov.Result != nil {
			mov.Result <- false
		}
	}
}


func posicionarChave(jogo *Jogo) {
	chaveX, chaveY := 5, 10
	jogo.Mapa[chaveY][chaveX] = Chave
}