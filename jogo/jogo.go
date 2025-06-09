package jogo

import (
	"bufio"
	"os"

	shared "t2/common"
)

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	simbolo   rune
	cor       Cor
	corFundo  Cor
	tangivel  bool // Indica se o elemento bloqueia passagem
}

// Jogo contém o estado fixo do jogo (mapa)
type Jogo struct {
	Mapa [][] Elemento // grade 2D representando o mapa
	StatusMsg string
}

// Elementos visuais do jogo
var (
	Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false}
)

func (e Elemento) Simbolo() rune {
	return e.simbolo
}

// Cria e retorna uma nova instância do jogo (apenas mapa vazio)
func JogoNovo() Jogo {
	return Jogo{}
}

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
func JogoCarregarMapa(nome string, jogo *Jogo) error {
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
		for _, ch := range linha {
			e := Vazio
			switch ch {
			case Parede.simbolo:
				e = Parede
			case Inimigo.simbolo:
				e = Inimigo
			case Vegetacao.simbolo:
				e = Vegetacao
			case Personagem.simbolo:
				// Ignorar a posição do personagem no mapa, pois agora cada jogador tem posição própria
				e = Vazio
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// Verifica se pode mover para a posição (x,y), considerando o mapa fixo e as posições dos jogadores
func JogoPodeMoverPara(jogo *Jogo, estado shared.EstadoJogo, x, y int) bool {
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}
	// Se o mapa fixo bloqueia
	if jogo.Mapa[y][x].tangivel {
		return false
	}

	// Se a posição já está ocupada por algum jogador
	for _, p := range estado.Players {
		if p.X == x && p.Y == y {
			return false
		}
	}

	return true
}

// Atualiza a posição do jogador no estado compartilhado, se possível
func JogoMoverJogador(estado *shared.EstadoJogo, id string, dx, dy int, jogo *Jogo) bool {
	p, ok := estado.Players[id]
	if !ok {
		return false // jogador não existe
	}
	nx, ny := p.X+dx, p.Y+dy
	if JogoPodeMoverPara(jogo, *estado, nx, ny) {
		p.X = nx
		p.Y = ny
		estado.Players[id] = p
		return true
	}
	return false
}
