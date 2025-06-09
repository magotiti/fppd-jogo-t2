package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"sync"
	shared "t2/common"
	jogoCore "t2/jogo"
	"time"
)

var sequence = 0

var (
	mu          sync.Mutex
	estadoAtual shared.EstadoJogo
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <playerID>")
		return
	}
	id := os.Args[1]

	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Erro ao conectar ao servidor RPC:", err)
	}

	// registra o jogador apos conectar
	var ok bool
	err = client.Call("Servidor.RegistrarJogador", id, &ok)
	if err != nil || !ok {
		log.Fatalf("Erro ao registrar jogador. ID já em uso ou posição inicial indisponível.")
	}

	// carreca o mapa local (somente no inicio do jogo)
	jogoLocal := jogoCore.JogoNovo()
	if err := jogoCore.JogoCarregarMapa("mapa.txt", &jogoLocal); err != nil {
		panic(err)
	}

	jogoCore.InterfaceIniciar()
	defer jogoCore.InterfaceFinalizar()

	// go routine para atualizar o estado do jogo constantemente
	go atualizarEstado(client, id)

	for {
		mu.Lock()
		estado := estadoAtual
		mu.Unlock()

		jogoCore.InterfaceDesenharJogo(&jogoLocal, estado)

		evento := jogoCore.InterfaceLerEventoTeclado()
		if evento.Tipo == "sair" {
			break
		}
		if evento.Tipo == "mover" {
			dx, dy := jogoCore.PersonagemMover(evento.Tecla)
			if dx == 0 && dy == 0 {
				continue
			}
			sequence++
			mov := shared.Movimento{
				ID:       id,
				DeltaX:   dx,
				DeltaY:   dy,
				Sequence: sequence,
			}
			var ack bool
			err := client.Call("Servidor.AtualizaPosicao", mov, &ack)
			if err != nil {
				log.Println("Erro ao enviar movimento:", err)
			}
		}
	}
}

func atualizarEstado(client *rpc.Client, id string) {
	for {
		var estado shared.EstadoJogo
		err := client.Call("Servidor.GetEstadoJogo", id, &estado)
		if err != nil {
			log.Println("Erro ao obter estado:", err)
			time.Sleep(time.Second)
			continue
		}
		mu.Lock()
		estadoAtual = estado
		mu.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}
