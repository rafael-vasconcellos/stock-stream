package client

import (
	"context"
	"io"
	"log"
	"stock_stream/pb"
	"time"

	"google.golang.org/grpc"
)


func RunClient() { 
	dial, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil { log.Println("Error connecting to server:")
		//panic(err) 
	}
	defer dial.Close()

	client := pb.NewStockServiceClient(dial)
	stream, err := client.StreamStockPrices(context.Background())
	if err != nil { log.Println("Error creating stream:")
		//panic(err) 
	}


	done := make(chan bool)
	go func() { 
		for { 
			response, err := stream.Recv()
			if err == io.EOF { break }
			if err != nil { break }

			log.Printf("Received price for symbol %s: $%.2f", response.GetSymbol(), response.GetPrice())
		}

		close(done)
	}()

	symbols := []string{"AAPL", "GOOGL", "MSFT", "AMZN", "FB"}
	for _, symbol := range symbols { 
		log.Println("Sending request for symbol:", symbol)
		err := stream.Send(&pb.StockRequest{Symbol: symbol})
		if err != nil { log.Println("Error sending request:")
			//panic(err) 
		}
		time.Sleep(time.Second * 3)
	}

	<-done
}