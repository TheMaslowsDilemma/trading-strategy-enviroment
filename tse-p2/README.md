
## Trading Strategy Environment Part Two
This repository contains a Golang-based trading strategy simulation that allows users to test and develop trading strategies in a dynamic, risk-free environment. Powered by an automated market maker inspired by Uniswap, the simulation enables competing strategies to generate synthetic candlestick data through trades.

#### Run the Simulation:

```
cd trading-strategy-simulation/tse-p2
go mod tidy
go run main.go
```

Configure simulation parameters (e.g., number of traders, initial reserves) in config.go.
Add custom strategies by implementing the Strategy interface.
Monitor synthetic candlestick data and trader performance via console output.

#### Contributing
Contributions are welcome! Feel free to submit pull requests or open issues for new features, bug fixes, or improvements.

##### Written with Grok
used it for these parts, but this was primarily a learning project so kept away from it for the vast majority.
- this README
- the `static/js/graph-ws.js` file for graphing candles
- the random strategies under `momentum.go`


##### License
This project is licensed under the MIT License. See the LICENSE file for details.
