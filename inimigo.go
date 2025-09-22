// inimigo.go - Lógica e estrutura dos inimigos
package main

import (
	"math"
	"math/rand"
	"time"
)

// define o estado de comportamento do inimigo (Patrulha ou Perseguição)
type ModoInimigo int

const (
	ModoPatrulha    ModoInimigo = iota // Movimento aleatório
	ModoPerseguicao                    // Movimento em direção ao jogador
)

// representa um único inimigo no jogo
type Inimigo struct {
	ID   int
	X, Y int
	Modo ModoInimigo
}

// é a mensagem que a goroutine do inimigo envia 
type PedidoAtualizacao struct {
	IDInimigo int
	Resposta  chan Inimigo // canal de resposta para devolver o estado atualizado
}

// função que roda uma goroutine para cada inimigo
// apenas pede ao gerenciador central o que fazer.
func rotinaInimigo(id int, canalPedidos chan<- PedidoAtualizacao) {
	// o inimigo tenta se mover a cada 800ms (pode ajustar para deixá-lo mais rápido/lento)
	ticker := time.NewTicker(400 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		// cria um canal de resposta para esta requisição específica
		canalResposta := make(chan Inimigo)

		// envia o pedido de atualização para o loop com o seu ID e o canal de resposta
		canalPedidos <- PedidoAtualizacao{
			IDInimigo: id,
			Resposta:  canalResposta,
		}

		// espera bloqueado até que o loop principal processe e responda
		<-canalResposta
		// o estado do inimigo já foi atualizado pelo loop principal. a goroutine não
		// precisa fazer mais nada, apenas esperar o próximo "tick" para pedir de novo.
	}
}

// função auxiliar para determinar a distância entre dois pontos
func calcularDistancia(x1, y1, x2, y2 int) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	return math.Sqrt(dx*dx + dy*dy)
}

// determina o melhor passo (dx, dy) para se aproximar do jogador
func calcularMovimentoPerseguicao(inimigoX, inimigoY, jogadorX, jogadorY int) (int, int) {
	dx, dy := 0, 0
	if jogadorX > inimigoX {
		dx = 1
	} else if jogadorX < inimigoX {
		dx = -1
	}

	if jogadorY > inimigoY {
		dy = 1
	} else if jogadorY < inimigoY {
		dy = -1
	}

	// evita movimento puramente diagonal se o movimento em uma direção já alinha o inimigo
	if dx != 0 && dy != 0 {
		if rand.Intn(2) == 0 {
			dx = 0
		} else {
			dy = 0
		}
	}

	return dx, dy
}