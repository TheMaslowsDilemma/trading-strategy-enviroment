package main

import (
	"fmt"
	"os"
	"tse-p1/candles"
	"tse-p1/simulation"
	"tse-p1/strategy"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var EthereumCSVDescriptor = candles.CandlesCsvDescriptor {
	Filepath: 		"../market-data/ETH_1min.csv",
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
	cs, err := candles.LoadCandlesFromCsv(EthereumCSVDescriptor)
	if err != nil {
		panic(err)
	}

	// SIMULATION STARTUP AND RUN //
	initialBalance := 12.0
	fee := 0.004
	strat := &strategy.SimpleStrategy{ShortInterval: 7, LongInterval: 24 * 60}
	sim := simulation.NewSimulator(initialBalance, strat, fee)
	ns := sim.Run(cs) // networth history of the bot


	// GRAPHING STUFF //
	var ps_ld []opts.LineData
	var nw_ld []opts.LineData

	x := 0
	nsmax := len(ns) - 1
	csmax := len(cs) - 1

	for i := 0; i < csmax; i++ {

		if i % 3000 == 0 {
			ps_ld = append(ps_ld, opts.LineData{Value: []interface{}{x, cs[i].Close}})
			if i < nsmax {
				nw_ld = append(nw_ld, opts.LineData{Value: []interface{}{x, ns[i]}})
			}

			x += 1
		}
	}

	//--- Graphing Setup --- //
	linegraph := charts.NewLine()
	linegraph.AddSeries("eth close price", ps_ld)
	linegraph.AddSeries(fmt.Sprintf("%s net worth", strat.GetName()), nw_ld)
	linegraph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Trading Strategy Performance", Subtitle: fmt.Sprintf("Initial balance of $%v and fee of %v%%", initialBalance, fee * 100)}),
	)

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