package simulation

func (s *Simulation) PlaceUserTrade(from, to string, confidence float64) {
    (&s.Mempool).PushTx(
        s.CliTrader.SwapTx(
            from,
            to,
            confidence,
            s.Ledger,
        ),
    )
}
