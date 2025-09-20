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

	// --- Cria o canal de movimento aqui ---
	var canalMovimento = make(chan Movimento)

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	posicionarChave(&jogo) //posiciona a chave no mapa

	// Inicializa os portais
	inicializarPortais(&jogo)

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

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

	canal := make(chan EventoJogo, 10)

	// goroutine para processar eventos
	go func() {
		for ev := range canal {
			switch ev.Tipo {
			case "pegarChave":
				jogo.TemChave = true
				jogo.StatusMsg = "Você pegou a chave! O portal vai começar a abrir e fechar."

				// ativa o portal
				for i := range jogo.Portais {
					go rotinaPortal(&jogo, &jogo.Portais[i])
				}
			}
		}
	}()

}