package main

import (
	"os"
	"tse-p1/candles"
	"tse-p1/simulation"
	"tse-p1/strategy"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

const EthereumCSVPath = "../market-data/ETH_1min.csv"
var EthereumCSVDescriptor = candles.CandlesCsvDescriptor {
	TimestampIndex: 1,
	OpenIndex:      3,
	HighIndex:      4,
	LowIndex:       5,
	CloseIndex:     6,
	VolumeIndex:    7,
	ColumnCount:    8,
	Order:          candles.OldestLast,
}

func main() {
	candles, err := candles.LoadCandlesFromCsv(EthereumCSVPath, EthereumCSVDescriptor)
	if err != nil {
		panic(err)
	}

	strategy := &strategy.SimpleMAStrategy{ShortPeriod: 500, LongPeriod: 5000}
	sim := simulation.NewSimulator(1000.0, strategy, 0.015)
	networth_history := sim.Run(candles)

	var price_points []opts.LineData
	var netw_points []opts.LineData
	x := 0

	netwlen := len(networth_history) - 1
	for i := 0; i < len(candles) - 1; i++ {
		if i % 1000 == 0 {
			price_points = append(price_points, opts.LineData{Value: []interface{}{x, candles[i].Close}})

			if i < netwlen {
				netw_points = append(netw_points, opts.LineData{Value: []interface{}{x, networth_history[i]}})
			}

			x += 1

		}
	}

	linegraph := charts.NewLine()
	linegraph.AddSeries("Close Price", price_points)
	linegraph.AddSeries("Bot Networth", netw_points)

	f, err := os.Create("linegraph.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = linegraph.Render(f)
	if err != nil {
		panic(err)
	}
}