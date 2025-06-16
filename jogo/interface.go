package jogo

import (
	shared "t2/common"

	"github.com/nsf/termbox-go"
)

// Define um tipo Cor para encapsuladar as cores do termbox
type Cor = termbox.Attribute

// Definições de cores utilizadas no jogo
const (
    CorPadrao     Cor = termbox.ColorDefault
    CorCinzaEscuro    = termbox.ColorDarkGray
    CorVermelho       = termbox.ColorRed
    CorVerde          = termbox.ColorGreen
    CorParede         = termbox.ColorBlack | termbox.AttrBold | termbox.AttrDim
    CorFundoParede    = termbox.ColorDarkGray
    CorTexto          = termbox.ColorDarkGray
)

// EventoTeclado representa uma ação detectada do teclado (como mover, sair ou interagir)
type EventoTeclado struct {
    Tipo  string // "sair", "interagir", "mover"
    Tecla rune   // Tecla pressionada, usada no caso de movimento
}

// Inicializa a Interface gráfica usando termbox
func InterfaceIniciar() {
    if err := termbox.Init(); err != nil {
        panic(err)
    }
}

// Encerra o uso da Interface termbox
func InterfaceFinalizar() {
    termbox.Close()
}

// Lê um evento do teclado e o traduz para um EventoTeclado
func InterfaceLerEventoTeclado() EventoTeclado {
    ev := termbox.PollEvent()
    if ev.Type != termbox.EventKey {
        return EventoTeclado{}
    }
    if ev.Key == termbox.KeyEsc {
        return EventoTeclado{Tipo: "sair"}
    }
    if ev.Ch == 'e' {
        return EventoTeclado{Tipo: "interagir"}
    }
    return EventoTeclado{Tipo: "mover", Tecla: ev.Ch}
}

// Renderiza todo o estado atual do jogo na tela
func InterfaceDesenharJogo(jogo *Jogo, estado shared.EstadoJogo) {
    InterfaceLimparTela()

    // Desenha todos os elementos do mapa recebido do servidor
    for y, linha := range estado.Mapa {
        for x, elem := range linha {
            InterfaceDesenharElemento(x, y, Elemento{simbolo: elem, cor: CorPadrao, corFundo: CorPadrao})
        }
    }

    // Desenha todos os jogadores sobre o mapa
    for _, p := range estado.Players {
        InterfaceDesenharElemento(p.X, p.Y, Personagem)
    }

    // Desenha a barra de status (opcional, pode remover se não usar jogo.StatusMsg)
    InterfaceDesenharBarraDeStatus(jogo)

    // Força a atualização do terminal
    InterfaceAtualizarTela()
}

// Limpa a tela do terminal
func InterfaceLimparTela() {
    termbox.Clear(CorPadrao, CorPadrao)
}

// Força a atualização da tela do terminal com os dados desenhados
func InterfaceAtualizarTela() {
    termbox.Flush()
}

// Desenha um elemento na posição (x, y)
func InterfaceDesenharElemento(x, y int, elem Elemento) {
    termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
}

// Exibe uma barra de status com informações úteis ao jogador
func InterfaceDesenharBarraDeStatus(jogo *Jogo) {
    // Linha de status dinâmica
    if jogo != nil {
        for i, c := range jogo.StatusMsg {
            termbox.SetCell(i, 1, c, CorTexto, CorPadrao)
        }
    }

    // Instruções fixas
    msg := "Use WASD para mover e E para interagir. ESC para sair."
    for i, c := range msg {
        termbox.SetCell(i, 3, c, CorTexto, CorPadrao)
    }
}