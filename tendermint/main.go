package main

import (
	"fmt"
	corerpc "github.com/implement-pbft-pos/tendermint/corerpc"
	"time"
)

func main(){
	fmt.Println("Starting...")

	go func(){
		fmt.Println("Client Call...")
		time.Sleep(5 * time.Second)
		gClientCfg := corerpc.GClientConfig{Host: "0.0.0.0", Port: 9000}
		client := corerpc.NewGClient(gClientCfg)
		msg := corerpc.GProposeMessage{Round: 1, Height: 100, ValidRound: -1, Data: &corerpc.GData{Data: "hehe"}}
		client.SendProposeMessage(&msg)
	}()

	gServerCfg := corerpc.GServerConfig{Port : 9000}
	gServer := corerpc.GServer{ServerCfg: gServerCfg}
	gServer.StartServer()




	//result := new corerpc.GResult{Error: 0, Data: "hehe"}
	//fmt.Println("result: ", result)
	//hello()
}