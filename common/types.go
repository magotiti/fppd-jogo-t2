package shared

type EstadoPlayer struct {
	ID       string
	X, Y     int
	Sequence int
}

type EstadoJogo struct {
	Players map[string]EstadoPlayer // armazenar posicao e sequencia dos jogadores
	Mapa    [][]rune                // refresh do mapa
}

type Movimento struct {
	ID       string
	DeltaX   int
	DeltaY   int
	Sequence int
}