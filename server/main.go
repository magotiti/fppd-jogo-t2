package main

import (
	"log"
	"net"
	"net/rpc"
	"sync"

	shared "t2/common"
	"t2/jogo"
)

type Servidor struct {
	mu        sync.Mutex
	players   map[string]shared.EstadoPlayer
	sequences map[string]int
	jogo      *jogo.Jogo
}

// registra o jogador a partir do evento no client
func (server *Servidor) RegistrarJogador(id string, ack *bool) error {
	server.mu.Lock()
	defer server.mu.Unlock()

	// ja temos esse id
	if _, existe := server.players[id]; existe {
		*ack = false
		return nil
	}

	estadoAtual := shared.EstadoJogo{
		Players: server.players,
	}

	var pos [2]int
	found := false
	for y, linha := range server.jogo.Mapa {
		for x := range linha {
			if jogo.JogoPodeMoverPara(server.jogo, estadoAtual, x, y) {
				pos = [2]int{x, y}
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		*ack = false
		return nil
	}

	server.players[id] = shared.EstadoPlayer{
		ID : id,
		X  : pos[0],
		Y  : pos[1],
		Sequence: 0,
	}
	server.sequences[id] = 0

	*ack = true
	return nil
}

// processa a demanda de acao do cliente
func (server *Servidor) AtualizaPosicao(cmd shared.Movimento, ack *bool) error {
	server.mu.Lock()
	defer server.mu.Unlock()

	player, existe := server.players[cmd.ID]
	if !existe {
		*ack = false
		return nil
	}

	lastSeq := server.sequences[cmd.ID]
	if cmd.Sequence <= lastSeq {
		*ack = true
		return nil
	}

	nx := player.X + cmd.DeltaX
	ny := player.Y + cmd.DeltaY

	// monta o estado atual para passar a funcao de validacao
	estadoAtual := shared.EstadoJogo{
		Players: server.players,
	}

	if jogo.JogoPodeMoverPara(server.jogo, estadoAtual, nx, ny) {
		player.X = nx
		player.Y = ny
		player.Sequence = cmd.Sequence
		server.players[cmd.ID] = player
		server.sequences[cmd.ID] = cmd.Sequence
	}

	*ack = true
	return nil
}

// retorna o estado atual do jogo
func (server *Servidor) GetEstadoJogo(id string, state *shared.EstadoJogo) error {
	server.mu.Lock()
	defer server.mu.Unlock()

	var mapaVisual [][]rune
	for _, linha := range server.jogo.Mapa {
		var linhaRunas []rune
		for _, elem := range linha {
			linhaRunas = append(linhaRunas, elem.Simbolo())
		}
		mapaVisual = append(mapaVisual, linhaRunas)
	}

	*state = shared.EstadoJogo{
		Players: server.players,
		Mapa:    mapaVisual,
	}
	return nil
}

func main() {
	j := jogo.JogoNovo()
	if err := jogo.JogoCarregarMapa("mapa.txt", &j); err != nil {
		log.Fatalf("Erro ao carregar mapa: %v", err)
	}

	s := &Servidor{
		players:   make(map[string]shared.EstadoPlayer),
		sequences: make(map[string]int),
		jogo:      &j,
	}

	if err := rpc.Register(s); err != nil {
		log.Fatalf("Erro ao registrar servidor RPC: %v", err)
	}

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("Erro ao escutar na porta 1234: %v", err)
	}
	defer listener.Close()

	log.Println("Servidor escutando em :1234...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Erro ao aceitar conexao: %v", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
