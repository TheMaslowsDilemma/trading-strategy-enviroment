### Part Two: A Simulated Marketplace

#### Background
I want an environment for strategies to directly compete with each other. This can be done by simulating a market place for strategies to buy or sell assets, where their decisions also affect the market. For this part we will temporarily move away from using historical data.

#### System Overview

- **Simulation:** defines start logic, result output, and cleanup logic. basically manages lifecyle.

- **Ledger:** holds state information for the simulation - currently just Traders, Mempool, and Exchanges

- **Miner:** is responsible for updating the ledger through the creation and application of Transaction Blocks

- **Exchange:** currently only *ConstantProductExchange* which is state information on the liquidity pools.
    - defines the *Swap* transaction, currently only `SwapExact_For_` could be updated to include `Swap_ForExact_`


- **Trader:** makes decisions (transactions) based on wallet state, candle history, and strategy.

- **Mempool:** holds pending transactions, and is used by the *Miner* to fill transaction blocks
