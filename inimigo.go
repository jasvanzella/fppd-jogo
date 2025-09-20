// inimigo.go - Lógica e estrutura dos inimigos
package main

import (
	"math"
	"math/rand"
	"time"
)

// ModoInimigo define o estado de comportamento do inimigo (Patrulha ou Perseguição)
type ModoInimigo int

const (
	ModoPatrulha    ModoInimigo = iota // Movimento aleatório
	ModoPerseguicao                    // Movimento em direção ao jogador
)

// Inimigo representa um único inimigo no jogo, com seu estado completo
type Inimigo struct {
	ID   int
	X, Y int
	Modo ModoInimigo
}

// PedidoAtualizacao é a mensagem que a goroutine do inimigo envia ao "Guardião do Jogo"
type PedidoAtualizacao struct {
	IDInimigo int
	Resposta  chan Inimigo // Canal de resposta para devolver o estado atualizado
}

// rotinaInimigo é a função que roda em uma goroutine para cada inimigo.
// Ele apenas pede ao gerenciador central o que fazer.
func rotinaInimigo(id int, canalPedidos chan<- PedidoAtualizacao) {
	// O inimigo tenta se mover a cada 800ms (pode ajustar para deixá-lo mais rápido/lento)
	ticker := time.NewTicker(400 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		// Cria um canal de resposta para esta requisição específica
		canalResposta := make(chan Inimigo)

		// Envia o pedido de atualização para o Guardião com o seu ID e o canal de resposta
		canalPedidos <- PedidoAtualizacao{
			IDInimigo: id,
			Resposta:  canalResposta,
		}

		// Espera bloqueado até que o Guardião processe e responda
		<-canalResposta
		// O estado do inimigo já foi atualizado pelo Guardião. A goroutine não
		// precisa fazer mais nada, apenas esperar o próximo "tick" para pedir de novo.
	}
}

// calcularDistancia é uma função auxiliar para determinar a distância entre dois pontos
func calcularDistancia(x1, y1, x2, y2 int) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	return math.Sqrt(dx*dx + dy*dy)
}

// calcularMovimentoPerseguicao determina o melhor passo (dx, dy) para se aproximar do jogador
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

	// Evita movimento puramente diagonal se o movimento em uma direção já alinha o inimigo
	if dx != 0 && dy != 0 {
		if rand.Intn(2) == 0 {
			dx = 0
		} else {
			dy = 0
		}
	}

	return dx, dy
}