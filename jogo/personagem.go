package jogo

import (
	"fmt"
	shared "t2/common"
)

// personagemMover cria um comando de movimento baseado na tecla, retornando o delta e a tecla
func PersonagemMover(tecla rune) (int, int) {
	dx, dy := 0, 0
	switch tecla {
	case 'w':
		dy = -1
	case 'a':
		dx = -1
	case 's':
		dy = 1
	case 'd':
		dx = 1
	}
	return dx, dy
}

// personagemInteragir poderia criar um comando de interação ou atualizar status local
func PersonagemInteragir(jogo *Jogo, estado shared.EstadoJogo, playerID string) {
	jogo.StatusMsg = fmt.Sprintf("Interagindo no jogador %s em (%d, %d)", playerID, estado.Players[playerID].X, estado.Players[playerID].Y)
}

// personagemExecutarAcao monta o comando movimento e retorna os dados para enviar ao servidor
func PersonagemExecutarAcao(ev EventoTeclado, estado shared.EstadoJogo, playerID string, seq int) (bool, *shared.Movimento) {
	switch ev.Tipo {
	case "sair":
		return false, nil
	case "interagir":
		// Aqui pode-se implementar a lógica de interação, atualmente só status local
		// Precisaria da referência ao jogo para atualizar StatusMsg, ou poderia enviar comando ao servidor
		return true, nil
	case "mover":
		dx, dy := PersonagemMover(ev.Tecla)
		if dx == 0 && dy == 0 {
			return true, nil
		}
		mov := &shared.Movimento{
			ID:       playerID,
			DeltaX:   dx,
			DeltaY:   dy,
			Sequence: seq + 1,
		}
		return true, mov
	}
	return true, nil
}
