package server

import (
	"io"
	"log"
	"math/rand/v2"
	"net"
	"stock_stream/pb"
	"time"

	"google.golang.org/grpc"
);

type StockServiceServer struct { 
	pb.UnimplementedStockServiceServer
}

func (s *StockServiceServer) StreamStockPrices(stream pb.StockService_StreamStockPricesServer) error { 
	for { 
		req, err := stream.Recv()
		if err == io.EOF { return nil }
		if err != nil { return err }
		symbol := req.GetSymbol()
		log.Printf("Received request for symbol: %s", symbol)

		go func(symbol string) { 
			for i:=0; i<50; i++ { 
				price := rand.Float32() * 100
				err := stream.Send(&pb.StockResponse{
					Symbol: symbol, 
					Price: price,
				})
				if err != nil { 
					log.Printf("Error sending price for symbol %s: %v", symbol, err)
					return
				}

				time.Sleep(time.Second * 1)
			}
		}(symbol)
	}
}


func RunServer() { 
	listen, err := net.Listen("tcp", ":50051")
	if err != nil { log.Println("Error starting server:") 
		//panic(err) 
	}

	gRPCServer := grpc.NewServer()
	pb.RegisterStockServiceServer(gRPCServer, &StockServiceServer{})
	if err := gRPCServer.Serve(listen); err != nil { log.Println("Error starting server:")
		//panic(err) 
	}
}