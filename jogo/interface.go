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

    // desenha todos os elementos do mapa base (com cor)
    if jogo != nil {
        for y, linha := range jogo.Mapa {
            for x, elem := range linha {
                InterfaceDesenharElemento(x, y, elem)
            }
        }
    } else {
        // fallback: desenha só símbolos se não houver mapa base
        for y, linha := range estado.Mapa {
            for x, elem := range linha {
                InterfaceDesenharElemento(x, y, Elemento{simbolo: elem, cor: CorPadrao, corFundo: CorPadrao})
            }
        }
    }

    // Desenha todos os jogadores sobre o mapa
    for _, p := range estado.Players {
        InterfaceDesenharElemento(p.X, p.Y, Personagem)
    }

    // desenha a barra de status APÓS o mapa
    InterfaceDesenharBarraDeStatus(jogo)
    InterfaceAtualizarTela()
}

func InterfaceLimparTela() {
    termbox.Clear(CorPadrao, CorPadrao)
}

func InterfaceAtualizarTela() {
    termbox.Flush()
}

func InterfaceDesenharElemento(x, y int, elem Elemento) {
    termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
}

func InterfaceDesenharBarraDeStatus(jogo *Jogo) {
    linhaStatus := 0
    if jogo != nil && len(jogo.Mapa) > 0 {
        linhaStatus = len(jogo.Mapa) + 1
    } else {
        linhaStatus = 1
    }

    // Linha de status dinâmica
    if jogo != nil {
        for i, c := range jogo.StatusMsg {
            termbox.SetCell(i, linhaStatus, c, CorTexto, CorPadrao)
        }
    }

    // Instruções fixas
    msg := "Use WASD para mover e E para interagir. ESC para sair."
    for i, c := range msg {
        termbox.SetCell(i, linhaStatus+2, c, CorTexto, CorPadrao)
    }
}