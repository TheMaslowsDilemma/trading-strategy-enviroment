package candles

type CandleBuffer struct {
	Size uint
	Head uint
	Tail uint
	Cnds []Candle
}

func CreateCandleBuffer(size uint) CandleBuffer {
	return CandleBuffer {
		Size: size,
		Head: 0,
		Tail: 0,
		Cnds: make([]Candle, size),
	}
}

func (cb *CandleBuffer) Push(cndl Candle) {
	cb.Cnds[cb.Head] = cndl
	if cb.Head == (cb.Tail + cb.Size - 1) % cb.Size {
		cb.Tail = (cb.Tail + 1) % cb.Size
	}
	cb.Head = (cb.Head + 1) % cb.Size
}

func (cb CandleBuffer) GetCandles() []Candle {
	var (
		count 	uint
		i		uint
		cs		[]Candle
	)
	
	if cb.Head == cb.Tail {
		return []Candle{}
	}

	count = (cb.Head + cb.Size - cb.Tail) % cb.Size
	cs = make([]Candle, count)

	for i = 0; i < count; i++ {
		cs[i] = cb.Cnds[(cb.Tail + i) % cb.Size]
	}
	
	return cs
}

func (cb CandleBuffer) Clone() CandleBuffer {
	
	return CandleBuffer {
		Cnds: cb.Cnds,
		Head: cb.Head,
		Tail: cb.Tail,
		Size: cb.Size,
	}
}

func min(a uint, b uint) uint {
	if a > b {
		return b
	}
	return a
}
