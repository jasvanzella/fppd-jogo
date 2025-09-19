// main.go - Loop principal do jogo
package main

import (
	"os"
	"time"
)

func main() {
	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// --- Cria o canal de movimento aqui ---
	var canalMovimento = make(chan Movimento)

	// Inicia a goroutine que gerencia todas as alterações do mapa
	go gerenciarMapa(&jogo, canalMovimento)

	for _, inimigo := range jogo.Inimigos {
		go moverInimigo(inimigo.X, inimigo.Y, canalMovimento)
	}

	// --- Inicia a goroutine que redesenha a tela continuamente ---
	go func() {
		for {
			time.Sleep(100 * time.Millisecond) // atualiza a cada 0.1s
			interfaceDesenharJogo(&jogo)
		}
	}()

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			break
		}
		interfaceDesenharJogo(&jogo)
	}
}
