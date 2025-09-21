// main.go - Guardião do Jogo e loop principal
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const raioDeVisao = 12.0

// Arquivo: main.go

func gerenciarEventos(
	jogo *Jogo,
	canalTeclado <-chan EventoTeclado,
	canalPedidosInimigos <-chan PedidoAtualizacao,
	canalTickPortal <-chan struct{},
	canalTickBau <-chan struct{}, // canal p bau
) {
	// Loop de jogo principal.
	for {
		// --- LÓGICA DE ATUALIZAÇÃO DE ESTADO ---
		select { // O bloco select começa aqui
		case evento := <-canalTeclado:
			if continuar := personagemExecutarAcao(evento, jogo); !continuar {
				return
			}
		case pedido := <-canalPedidosInimigos:
			inimigo := jogo.Inimigos[pedido.IDInimigo]
			if jogo.Mapa[inimigo.Y][inimigo.X].simbolo == ElementoInimigo.simbolo {
				jogo.Mapa[inimigo.Y][inimigo.X] = Vazio
			}

			// --- LÓGICA DE MUDANÇA DE ESTADO COM MENSAGENS ---
			modoAnterior := inimigo.Modo
			dist := calcularDistancia(inimigo.X, inimigo.Y, jogo.PosX, jogo.PosY)

			if dist < raioDeVisao {
				inimigo.Modo = ModoPerseguicao
				if modoAnterior == ModoPatrulha {
					jogo.StatusMsg = fmt.Sprintf("Cuidado! O inimigo %d começou a te perseguir!", inimigo.ID)
				}
			} else {
				inimigo.Modo = ModoPatrulha
				if modoAnterior == ModoPerseguicao {
					jogo.StatusMsg = fmt.Sprintf("Ufa! O inimigo %d desistiu de te perseguir.", inimigo.ID)
				}
			}

			dx, dy := 0, 0
			if inimigo.Modo == ModoPerseguicao {
				dx, dy = calcularMovimentoPerseguicao(inimigo.X, inimigo.Y, jogo.PosX, jogo.PosY)
			} else { // Modo Patrulha
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
			}

			nx, ny := inimigo.X+dx, inimigo.Y+dy
			if jogoPodeMoverPara(jogo, nx, ny, true) {
				inimigo.X = nx
				inimigo.Y = ny
			}

			pedido.Resposta <- *inimigo

		// NOVO CASE PARA GERENCIAR O BAÚ
		case <-canalTickBau:
			jogo.Bau.Ativo = !jogo.Bau.Ativo

		// NOVO CASE ADICIONADO DENTRO DO SELECT
		case <-canalTickPortal:
			elemDestino := Elemento{'╬', CorAmarelo, CorPadrao, false, true}
			for i := range jogo.Portais {
				p := &jogo.Portais[i]
				p.Ativo = !p.Ativo // Alterna o estado
				if p.Ativo {
					jogo.Mapa[p.Y][p.X] = Portal
					jogo.Mapa[p.DestY][p.DestX] = elemDestino
				} else {
					jogo.Mapa[p.Y][p.X] = Vazio
					jogo.Mapa[p.DestY][p.DestX] = Vazio
				}
			}

		default:
			// Não faz nada, apenas permite que o loop continue
		} // <<<<---- O 'select' TERMINA AQUI, depois de todos os 'case' e 'default'

		// --- LÓGICA DE DESENHO (RENDERING) ---
		// (Essa parte já estava correta, mas agora não dará mais erro)

		if jogo.Bau.Ativo {
			jogo.Mapa[jogo.Bau.Y][jogo.Bau.X] = BauJogo
		} else {
			jogo.Mapa[jogo.Bau.Y][jogo.Bau.X] = Vazio
		} 

		for _, inimigo := range jogo.Inimigos {
			if jogo.Mapa[inimigo.Y][inimigo.X].simbolo == ElementoInimigo.simbolo {
				jogo.Mapa[inimigo.Y][inimigo.X] = Vazio
			}
		}
		for _, inimigo := range jogo.Inimigos {
			jogo.Mapa[inimigo.Y][inimigo.X] = ElementoInimigo
		}

		interfaceDesenharJogo(jogo)
		time.Sleep(60 * time.Millisecond)
	}
}

func main() {
	interfaceIniciar()
	defer interfaceFinalizar()

	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	canalTeclado := make(chan EventoTeclado)
	canalPedidosInimigos := make(chan PedidoAtualizacao)
	canalTickPortal := make(chan struct{})
	canalTickBau := make(chan struct{})

	// CORREÇÃO AQUI: Adicionado () no final para chamar a função
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			canalTickPortal <- struct{}{}
		}
	}() // <<--- PARÊNTESES ADICIONADOS AQUI

		go func() {
		ticker := time.NewTicker(4 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			canalTickBau <- struct{}{}
		}
	}()

	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// CORREÇÃO AQUI TAMBÉM: Adicionado () no final para chamar a função
	go func() {
		for {
			canalTeclado <- interfaceLerEventoTeclado()
		}
	}() // <<--- PARÊNTESES ADICIONADOS AQUI

	for _, inimigo := range jogo.Inimigos {
		go rotinaInimigo(inimigo.ID, canalPedidosInimigos)
	}

	gerenciarEventos(&jogo, canalTeclado, canalPedidosInimigos, canalTickPortal, canalTickBau)
}