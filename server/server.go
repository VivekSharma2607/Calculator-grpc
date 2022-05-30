package main

import (
	pb "calculator/calculatorpb"
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{
	pb.UnimplementedCalculatorServiceServer
}


func (*server) ComputeAverage(stream pb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("Compute Average RPC is invoked ... \n")
	sum := int32(0)
	count := 0
	for {
		req , err := stream.Recv()
		if err == io.EOF {
			average := float64(sum)/float64(count)
			stream.SendAndClose(&pb.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("Error reading the incoming stream %v \n" , err)
		}
		sum += req.GetNumber()
		count++
	}
}
 
func (*server) PrimeNumber(req *pb.PrimeNumberRequest , stream pb.CalculatorService_PrimeNumberServer) error{
	fmt.Printf("Received the Prime Number Server Request: %v\n" , req)
	number := req.GetNumber()

	div := int64(2)
	for number > 1 {
		if number%div == 0{
			stream.Send(&pb.PrimeNumberResponse{
				PrimeFactor: div,
			})
			number /= div
		} else {
			div++
			fmt.Printf("The divisor is increased %v" , div)
		}
	}
	return nil
}

func (*server) Sum(ctx context.Context , req *pb.SumRequest) (*pb.SumResponse, error) {
	fmt.Println("Sum Funcion is implemented")
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber
	sum := firstNumber + secondNumber
	res := &pb.SumResponse{
		SumResult: sum,
	}
	return res , nil
}
func main() {
	fmt.Println("Calculator Server Intiated")
	lis , err := net.Listen("tcp" , "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v" ,err)
	}
	s:= grpc.NewServer()
	pb.RegisterCalculatorServiceServer(s , &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed load ther server %v" , err)
	}
}