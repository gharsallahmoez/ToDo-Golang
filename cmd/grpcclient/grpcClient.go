package main

import (
	"context"
	"io"
	"log"
	"sync"
	"time"

	pb "github.com/3almadmoon/ameni-assignment/api/proto"
	config "github.com/3almadmoon/ameni-assignment/configs"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var mutex sync.Mutex
var wg sync.WaitGroup

// runClient
// create grpc client
// make a remote call
func runClient() {
	ctx, cancelTimeoutFunc := context.WithTimeout(context.Background(), 3*time.Second)
	conn, err := grpc.DialContext(ctx, viper.GetString("grpcserver.host"), grpc.WithInsecure(), grpc.WithBlock())
	cancelTimeoutFunc()
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	log.Printf("Dialled OK...")
	defer conn.Close()

	client := pb.NewTodoListServiceClient(conn)
	log.Printf("created client")

	log.Println("**********ADD**************")
	ctx1, canc1 := getContextWithTimeout(5 * time.Second)
	defer canc1()
	res1, er1 := client.AddToDo(*ctx1, &pb.ToDoItem{Title: "aaaaaaaa", Description: "bug", Status: pb.Status_TODO})
	log.Printf("ADD1: |RES|:\n %v, \n |ERROR|: \n %v", res1, er1)

	log.Println("**********GETALL**************")
	ctx2, canc2 := getContextWithTimeout(10 * time.Second)
	defer canc2()
	res2, er2 := client.GetAllToDo(*ctx2, &empty.Empty{})
	var elem string
	for {
		item, er := res2.Recv()
		if er == io.EOF {
			break
		}
		if er != nil {
			log.Fatalf("can't load items %v", er)
		}
		elem = item.ToDoItems.GetHash()
		log.Printf("GETALL 1: |RES|:\n %v, \n |ERROR|: \n %v", item, er2)
	}

	varia := elem
	log.Println("hash 1 ; ", varia)

	wg.Add(6)
	for i := 0; i < 5; i++ {
		go func() {
			log.Println("**********UPDATE**************")
			ctx3, canc3 := getContextWithTimeout(5 * time.Second)
			defer canc3()
			res3, er3 := client.UpdateToDo(*ctx3, &pb.UpdateToDoItem{Hash: varia, Status: pb.Status_INPROGRESS})
			log.Printf("UPDATE : |RES|:\n %v, \n |ERROR|: \n %v", res3, er3)
			wg.Done()
		}()
	}
	go func() {
		log.Println("**********DELETE**************")
		ctx4, canc4 := getContextWithTimeout(5 * time.Second)
		defer canc4()
		res4, er4 := client.DeleteToDo(*ctx4, &pb.DeleteToDoItem{Hash: varia})
		log.Printf("DELETE1: |RES|:\n %v, \n |ERROR|: \n %v", res4, er4)
		wg.Done()
	}()
	wg.Wait()
}

func getContextWithTimeout(timeout time.Duration) (*context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &ctx, cancel
}

func main() {
	if err := config.SetViper(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	runClient()
}
