package trader

import (
    "fmt"
    "tse-p2/token"
    "tse-p2/ledger"
    "tse-p2/candles"
    "tse-p2/strategy"
    "tse-p2/wallet"
    "tse-p2/exchange"
)

type Trader struct {
    ExchangeAddr    ledger.LedgerAddr
    ExchangeSymA    string
    ExchangeSymB    string
    WalletAddr      ledger.LedgerAddr
    DecisionStrategy strategy.Strategy
    HasPendingTx    bool
    LogChannel      chan string
}

/** NOTE: Traders are associated with a specific exchange **/
func CreateTrader(
    decisionStrategy strategy.Strategy, 
    logChannelSize int, 
    walletAddr ledger.LedgerAddr, 
    exchangeAddr ledger.LedgerAddr, 
    symbolA, symbolB string, 
    ledgerState ledger.Ledger,
) Trader {
    return Trader {
        ExchangeAddr:    exchangeAddr,
        ExchangeSymA:    symbolA,
        ExchangeSymB:    symbolB,
        WalletAddr:      walletAddr,
        DecisionStrategy: decisionStrategy,
        HasPendingTx:    false,
        LogChannel:      make(chan string, logChannelSize),
    }
}

func (t *Trader) MakeDecision(candles []candles.Candle, ledgerState ledger.Ledger) ledger.Tx {
    if t.HasPendingTx {
        // TODO: Log "Cannot make decision - pending transaction exists"
        return nil
    }

    transaction, err := t.createTransaction(candles, ledgerState)
    if err != nil {
        // TODO: Log error "Failed to create transaction: %v", err
        return nil
    }

    return transaction
}

func (t *Trader) createTransaction(candles []candles.Candle, ledgerState ledger.Ledger) (ledger.Tx, error) {
    tradingAction, confidence := t.DecisionStrategy.Decide(candles, ledgerState)
 
    switch tradingAction {
    case strategy.Hold:
        return nil, nil
    case strategy.Buy:
        return t.CreateSwapTransaction(t.ExchangeSymA, t.ExchangeSymB, confidence, ledgerState)
    case strategy.Sell:
        return t.CreateSwapTransaction(t.ExchangeSymB, t.ExchangeSymA, confidence, ledgerState)
    default:
        return nil, fmt.Errorf("unknown trading action: %v", tradingAction)
    }
}

func (t *Trader) CreateSwapTransaction(
    inputSymbol, outputSymbol string, 
    confidence float64, 
    ledgerState ledger.Ledger,
) (ledger.Tx, error) {
    wallet, err := wallet.WalletFromLedgerItem(ledgerState[t.WalletAddr])
    if err != nil {
        return nil, fmt.Errorf("failed to load wallet: %v", err)
    }

    inputTokenReserveAddr, err := wallet.GetReserveAddr(inputSymbol, ledgerState)
    if err != nil {
        return nil, fmt.Errorf("failed to get reserve address for symbol %s: %v", inputSymbol, err)
    }
    
    inputTokenReserve, err := token.TkrFromLedgerItem(ledgerState[inputTokenReserveAddr])
    if err != nil {
        return nil, fmt.Errorf("failed to load input token reserve for %s: %v", inputSymbol, err)
    }

    if inputTokenReserve.Amount <= 0 {
        return nil, fmt.Errorf("no balance available for %s", inputSymbol)
    }

    inputAmount := inputTokenReserve.Amount * confidence
    if inputAmount <= 0 {
        return nil, fmt.Errorf("calculated input amount is zero or negative")
    }

    if inputTokenReserve.Amount < inputAmount {
        return nil, fmt.Errorf("insufficient balance of %s", inputSymbol)
    }


    // BUG we do no calculation and accept any slippage :/
    minimumOutputAmount := 0.0 // Placeholder - needs real calculation

    swapTx := exchange.SwapExactTokensForTokensTx{
        SymbolIn:     inputSymbol,
        SymbolOut:    outputSymbol,
        AmountIn:     inputAmount,
        AmountMinOut: minimumOutputAmount,
        WalletAddr:   t.WalletAddr,
        ExchangeAddr: t.ExchangeAddr,
    }

    // Set pending transaction flag
    t.HasPendingTx = true
    
    return swapTx, nil
}