package trader

import (
    "fmt"
    "tse-p2/token"
    "tse-p2/ledger"
    "tse-p2/strategy"
    "tse-p2/wallet"
    "tse-p2/exchange"
    "tse-p2/candles"
)

type candleFetchFunc func() ([]candles.Candle, error)
type txPlacerFunc    func(ledger.Tx)

type Trader struct {
    ExchangeAddr    ledger.LedgerAddr
    ExchangeSymA    string
    ExchangeSymB    string
    WalletAddr      ledger.LedgerAddr
    DecisionStrategy strategy.Strategy
    HasPendingTx    bool
    LogChannel      chan string
    candleFetcher   candleFetchFunc
    ledgerLookup    ledger.LedgerFetcher
    txPlacer        txPlacerFunc
}

/** NOTE: Traders are associated with a specific exchange **/
func CreateTrader(
    decisionStrategy strategy.Strategy, 
    logChannelSize int, 
    walletAddr ledger.LedgerAddr, 
    exchangeAddr ledger.LedgerAddr, 
    symbolA, symbolB string, 
    cndlfetch candleFetchFunc,
    ldgrfetch ledger.LedgerFetcher,
    trxplacer txPlacerFunc,
) Trader {
    return Trader {
        ExchangeAddr:    exchangeAddr,
        ExchangeSymA:    symbolA,
        ExchangeSymB:    symbolB,
        WalletAddr:      walletAddr,
        DecisionStrategy: decisionStrategy,
        HasPendingTx:    false,
        LogChannel:      make(chan string, logChannelSize),
        candleFetcher:   cndlfetch,
        ledgerLookup:   ldgrfetch,
        txPlacer:        trxplacer,
    }
}

func (t *Trader) MakeDecision(cs []candles.Candle) ledger.Tx {
    if t.HasPendingTx {
        // TODO: Log "Cannot make decision - pending transaction exists"
        return nil
    }

    transaction, err := t.createTransaction(cs)
    if err != nil {
        // TODO: Log error "Failed to create transaction: %v", err
        return nil
    }

    return transaction
}

func (t *Trader) createTransaction(cs []candles.Candle) (ledger.Tx, error) {
    tradingAction, confidence := t.DecisionStrategy.Decide(cs)
 
    switch tradingAction {
    case strategy.Hold:
        return nil, nil
    case strategy.Buy:
        return t.CreateSwapTransaction(t.ExchangeSymA, t.ExchangeSymB, confidence)
    case strategy.Sell:
        return t.CreateSwapTransaction(t.ExchangeSymB, t.ExchangeSymA, confidence)
    default:
        return nil, fmt.Errorf("unknown trading action: %v", tradingAction)
    }
}

func (t *Trader) CreateSwapTransaction(
    inputSymbol, outputSymbol string, 
    confidence float64, 
) (ledger.Tx, error) {
    var (
        wlt                 *wallet.Wallet
        inputTokenReserve   *token.TokenReserve
        inputTokenReserveAddr ledger.LedgerAddr
        minimumOutputAmount float64
        swapTx              ledger.Tx
        err                 error
    )

    if t.HasPendingTx { 
        return nil, fmt.Errorf("an existing tx is waiting to be processed.")
    }

    wlt, err = wallet.WalletFromLedgerItem(t.ledgerLookup(t.WalletAddr))
    if err != nil {
        return nil, fmt.Errorf("failed to load wallet: %v", err)
    }

    inputTokenReserveAddr, err = wlt.GetReserveAddr(inputSymbol, t.ledgerLookup)
    if err != nil {
        return nil, fmt.Errorf("failed to get reserve address for symbol %s: %v", inputSymbol, err)
    }
    
    inputTokenReserve, err = token.TkrFromLedgerItem(t.ledgerLookup(inputTokenReserveAddr))
    if err != nil {
        return nil, fmt.Errorf("failed to load input token reserve for %s: %v", inputSymbol, err)
    }

    if inputTokenReserve.Amount <= 0 {
        return nil, fmt.Errorf("no balance available for %s", inputSymbol)
    }

    inputAmount := inputTokenReserve.Amount * 0.07 * confidence
    if inputAmount <= 0 {
        return nil, fmt.Errorf("calculated input amount is zero or negative")
    }

    if inputTokenReserve.Amount < inputAmount {
        return nil, fmt.Errorf("insufficient balance of %s", inputSymbol)
    }


    // TODO use 95% spot-price * inputAmount = minimumOutput -- 5% slippage? //
    minimumOutputAmount = 0.0

    swapTx = exchange.SwapExactTokensForTokensTx{
        SymbolIn:     inputSymbol,
        SymbolOut:    outputSymbol,
        AmountIn:     inputAmount,
        AmountMinOut: minimumOutputAmount,
        WalletAddr:   t.WalletAddr,
        ExchangeAddr: t.ExchangeAddr,
        Callback: t.notifySwap,
    }

    t.HasPendingTx = true
    
    return swapTx, nil
}


func (t *Trader) notifySwap(_ ledger.TxResult) {
    t.HasPendingTx = false
}