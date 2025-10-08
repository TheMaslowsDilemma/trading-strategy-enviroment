package exchange

import (
    "fmt"
    "tse-p2/ledger"
    "tse-p2/wallet"
    "tse-p2/token"
)

type SwapExactTokensForTokensTx struct {
    SymbolIn            string
    SymbolOut           string
    AmountIn            uint64
    AmountMinOut        uint64
    WalletAddr          ledger.LedgerAddr
    ExchangeAddr        ledger.LedgerAddr
}

// Returns a Partial Ledger -- modifications only -- to be merged by the miner.
func (tx SwapExactTokensForTokensTx) Apply(l ledger.Ledger) (ledger.Ledger, error) {
    var (
        exg             ConstantProductExchange
        w               *wallet.Wallet
        waddrO          ledger.LedgerAddr   // wallet tk "out" reserveaddr
        waddrI          ledger.LedgerAddr   // wallet tk "in"  reserveaddr
        wtkrO           *token.TokenReserve
        wtkrI           *token.TokenReserve
        etkrA           *token.TokenReserve // exchange A reserve
        etkrB           *token.TokenReserve // exchange B reserve
        amtO            uint64
        ldgp            ledger.Ledger
        err             error
    )

    w, err = WalletFromLedgerItem(l[tx.WalletAddr])
    if err != nil {
        return nil, fmt.Errorf("failed to cast wallet: %v", err)
    }
    
    waddrI, err = w.GetReserveAddr(tx.SymbolIn, l)
    if err != nil {
        return nil, fmt.Errorf("failed to find wallet's source reserve: %v", err)
    }
    wtkrI, err = token.TkrFromLedgerItem(l[waddrI])
    if err != nil {
        return nil, fmt.Errorf("failed to cast wallet's source reserve %v", err)
    }

    waddrO, err = w.GetReserveAddr(tx.SymbolOut, l)
    if err != nil {
        return nil, fmt.Errorf("failed to find wallet's destination reserve: %v", err)
    }
    wtkrO, err = token.TkrFromLedgerItem(l[waddrO])
    if err != nil {
        return nil, fmt.Errorf("failed to cast wallet's destination reserve: %v", err)
    }

    exg, err = CpeFromLedgerItem(l[tx.ExchangeAddr])
    if err != nil {
        return nil, fmt.Errorf("failed to cast exchange: ", err)
    }

    etkrA, err = token.TkrFromLedgerItem(l[exg.TokenReserveA])
    if err != nil {
        return fmt.Errorf("failed to cast exchange's tkr-a: %v", err)
    }

    etkrB, err = token.TkrFromLedgerItem(l[exg.TokenReserveB])
    if err != nil {
        return fmt.Errorf("failed to cast exchange's tkr-b: %v", err)
    }
    
    if wtkrI.Amount < tx.AmountIn {
        return fmt.Errorf("insufficient funds for swap")
    }
    
    // Find out how we want to swap & if we even can
    if etkrA.Symbol == tx.SymbolIn && etkrB.Symbol == tx.SymbolOut {
        amtO, err = exg.SwapAForB(l, tx.AmountIn)
        if err != nil {
            return fmt.Errorf("swap a for b failed: %v", err)
        }

        if amtO < tx.AmountMinOut {
            return fmt.Errorf("swap slippage too high")
        }
        
        // Make the Diff Ledger
        ldgp = make(ledger.Ledge)

        // Move In Funds to the Exchange
        wtkrIn.Amount -= tx.AmountIn
        etkrA.Amount += tx.AmountIn
        
        // Move Out Funds to the Wallet
        etkrB.Amount -= amtO
        wtkrOut.Amount += amtO

        // Create Diff Ledger & Return
        ldgp[waddrI] = wtkrI
        ldgp[waddrO] = wtkrO
        ldgp[w.TokenReserveA] = etkrA
        ldgp[w.TokenReserveB] = etkrB

        return ldgp, nil
        
    } else if etkrB.Symbol == tx.SymbolIn && etkrA.Symbol == tx.SymbolOut {
        amtO, err = exg.SwapBForA(l, tx.AmountIn)
        if err != nil {
            return fmt.Errorf("swap b for a failed: %v", err)
        }

        if amtO < tx.AmountMinOut {
            return fmt.Errorf("swap slippage too high")
        }
        
        // Make the Diff Ledger
        ldgp = make(ledger.Ledge)

        // Move In Funds to the Exchange
        wtkrIn.Amount -= tx.AmountIn
        etkrA.Amount += tx.AmountIn
        
        // Move Out Funds to the Wallet
        etkrB.Amount -= amtO
        wtkrOut.Amount += amtO

        // Create Diff Ledger & Return
        ldgp[waddrI] = wtkrI
        ldgp[waddrO] = wtkrO
        ldgp[w.TokenReserveA] = etkrA
        ldgp[w.TokenReserveB] = etkrB

        return ldgp, nil
        
        // TODO this side of the swap !
    }
    return fmt.Errorf("failed to match symbols tx{ %v -> %v } != ex{ %v <-> %v }",
        tx.SymbolIn,
        tx.SymbolOut,
        etkrA.Symbol,
        etkrB.Symbol,
    )
}
