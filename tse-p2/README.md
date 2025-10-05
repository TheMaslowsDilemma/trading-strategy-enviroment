### Part Two: A Simulated Marketplace

#### Background
I want an environment for strategies to directly compete with each other. This can be done by simulating a market place for strategies to buy or sell assets, where their decisions also affect the market. For this part we will temporarily move away from using historical data.

#### Structural Overview
- **Simulation** tracks of ticks and closes the program once complete. interfaces with SDL2 to display the current state of candle history. decides on a random order for the traders to trade in after each round and their initial balance.

- **Market** tracks its pools and trades. it implements UNI swap AMM

- **Pool** is a liquidity pool with an initial amount of tokens and some value per token.

- **Candle** holds open, close, high, low, and volume data

- **Trader** tracks balance and uses a strategy to make decisions

- **Strategy**: same idea as `tse-p1`

#### System Overview (will remove above shortly)

- **Simulation:** defines simulation fields, and start logic, result output, and cleanup logic

- **Ledger:** holds state information for the simulation. this includes traders, and exchanges.

- **Miner:** is responsible for updating the ledger through the creation and application of Transaction Blocks

- **Trader:** makes decisions (transactions) based on wallet state, candle history, and strategy

- **Mempool:** holds pending transactions, and is used by the *Miner* to fill transaction blocks
