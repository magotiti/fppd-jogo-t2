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

    // carrega o mapa base local (para desenhar com cor)
    jogoLocal := jogoCore.JogoNovo()
    if err := jogoCore.JogoCarregarMapa("mapa.txt", &jogoLocal); err != nil {
        panic(err)
    }

    jogoCore.InterfaceIniciar()
    defer jogoCore.InterfaceFinalizar()

    // canal para sinalizar atualização do estado
    atualizaTela := make(chan struct{}, 1)

    // go routine para atualizar o estado do jogo constantemente
    go atualizarEstado(client, id, atualizaTela)

    // go routine para ler eventos do teclado
    eventos := make(chan jogoCore.EventoTeclado, 1)
    go func() {
        for {
            eventos <- jogoCore.InterfaceLerEventoTeclado()
        }
    }()

    // desenha a tela inicialmente
    mu.Lock()
    estado := estadoAtual
    mu.Unlock()
    jogoCore.InterfaceDesenharJogo(&jogoLocal, estado)

    for {
        select {
        case <-atualizaTela:
            mu.Lock()
            estado := estadoAtual
            mu.Unlock()
            jogoCore.InterfaceDesenharJogo(&jogoLocal, estado)
        case evento := <-eventos:
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
            if evento.Tipo == "sair" {
                var ack bool
                err := client.Call("Servidor.DesconectarJogador", id, &ack)
                if err != nil {
                    log.Println("Erro ao desconectar jogador:", err)
                }
                return
            }
        }
    }
}

func atualizarEstado(client *rpc.Client, id string, atualizaTela chan<- struct{}) {
    for {
        var estado shared.EstadoJogo
        err := client.Call("Servidor.GetEstadoJogo", id, &estado)
        if err != nil {
            time.Sleep(time.Second)
            continue
        }
        mu.Lock()
        estadoAtual = estado
        mu.Unlock()
        // sinaliza para redesenhar a tela
        select {
        case atualizaTela <- struct{}{}:
        default:
        }
        time.Sleep(100 * time.Millisecond)
    }
}