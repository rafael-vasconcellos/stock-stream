package main


import ( 
	"stock_stream/server"
	"stock_stream/client"
)

func main() { 
	go server.RunServer()
	go client.RunClient()

	ch := make(chan int)
	<-ch
}