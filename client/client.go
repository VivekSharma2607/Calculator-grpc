package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	// "net"
	pb "calculator/calculatorpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client of the calculator is intiatiated")
	cc , err := grpc.Dial("localhost:50051" , grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect : %v" , err)
	}
	defer cc.Close()
	c := pb.NewCalculatorServiceClient(cc)
	doSum(c)
	//doPrime(c)
	//doComputeAverage(c)
	//doFindMaximum(c)
}

func doFindMaximum(c pb.CalculatorServiceClient) {
	fmt.Printf("Starting the Find Maximum in Client \n")
	stream , err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while opeaning the server file in :%v" , err)
	}
	waitc := make(chan struct{})
	//send go routine
	go func()  {
		numbers := []int32{23, 34 ,89 , 56 , 76}
		for _, num := range numbers {
			fmt.Printf("Sending the number : %v \n" , num)
			stream.Send(&pb.FindMaximumRequest{
				Number: num,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	//receving Go routine
	go func() {
		for {
			res , err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while clinet from the server stream %v \n" , err)
			}
			maximum := res.GetMaximum()
			fmt.Printf("The new Maximum : %v \n" , maximum)

		}
		close(waitc)
	}()
	<-waitc
}


func doComputeAverage(c pb.CalculatorServiceClient) {
	fmt.Printf("Streaming the Compute Error RPC.... \n")

	stream , err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opeaning the stream %v \n" , err)
	}
	numbers := []int32{3,45,56,76}
	for _, number := range numbers {
		stream.Send(&pb.ComputeAverageRequest{
			Number: number,
		})
	}
	res , err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("The error while computing the Average : %v \n" , err)
	}
	fmt.Printf("The Average is : %v" , res.GetAverage())
}	


func doSum(c pb.CalculatorServiceClient) {
	fmt.Printf("The client is intiated in the Calcualtor clinet \n")
	req := &pb.SumRequest{
		FirstNumber: 5,
		SecondNumber: 15,
	}
	res , err := c.Sum(context.Background() , req)
	if err != nil{
		log.Fatalf("The Error while calling the Sum RPC: %v" , err)
	}
	log.Printf("Response from the Server: %v\n" , res.SumResult)
}

func doPrime(c pb.CalculatorServiceClient) {
	fmt.Printf("The client for Prime Number is intiated \n")
	req := &pb.PrimeNumberRequest{
		Number: 15,
	}
	result , err := c.PrimeNumber(context.Background() , req)
	if err != nil {
		log.Fatalf("The Error while calling the Prime:  %v",err)
	}
	for {
		res , err := result.Recv()
		if err == io.EOF{
			break
		}
		if err != nil {
			log.Fatalf("The error while reading the response: %v" , err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
}