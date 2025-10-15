package website

import (
	"fmt"
	"net/http"

)

type homeData struct {
	Wallet string
	ExAddr string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := homeData{
		Wallet: fmt.Sprintf("%v",Sim.CliWallet),
		ExAddr: fmt.Sprintf("%v", Sim.ExAddr),
	}
	homeTemplate.Execute(w, data)
}