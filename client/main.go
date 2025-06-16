package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
	"sync"
	"time"

	shared "t2/common"
	jogoCore "t2/jogo"
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

    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Digite o endereco do servidor: ")
    endereco, _ := reader.ReadString('\n')
    endereco = strings.TrimSpace(endereco)

    client, err := rpc.Dial("tcp", endereco)
    if err != nil {
        log.Fatal("Erro ao conectar ao servidor RPC:", err)
    }

    // registra o jogador apos conectar
    var ok bool
    err = client.Call("Servidor.RegistrarJogador", id, &ok)
    if err != nil || !ok {
        log.Fatalf("Erro ao registrar jogador. ID já em uso ou posição inicial indisponível.")
    }

    jogoCore.InterfaceIniciar()
    defer jogoCore.InterfaceFinalizar()

    // go routine para atualizar o estado do jogo constantemente
    go atualizarEstado(client, id)

    for {
        mu.Lock()
        estado := estadoAtual
        mu.Unlock()

        // Agora desenha apenas o estado recebido do servidor
        jogoCore.InterfaceDesenharJogo(nil, estado)

        evento := jogoCore.InterfaceLerEventoTeclado()
        // anda
        if evento.Tipo == "mover" {
            dx, dy := jogoCore.PersonagemMover(evento.Tecla)
            if dx == 0 && dy == 0 {
                continue
            }
            mov := shared.Movimento{
                ID:       id,
                DeltaX:   dx,
                DeltaY:   dy,
                Sequence: sequence + 1,
            }
            var ack bool
            err := client.Call("Servidor.AtualizaPosicao", mov, &ack)
            if err != nil {
                log.Println("Erro ao enviar movimento:", err)
            }
            if ack {
                sequence++
            }
        }
        // quita
        if evento.Tipo == "sair" {
            var ack bool
            err := client.Call("Servidor.DesconectarJogador", id, &ack)
            if err != nil {
                log.Println("Erro ao desconectar jogador:", err)
            }
            break
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