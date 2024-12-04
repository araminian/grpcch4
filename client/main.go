package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/araminian/grpcch4/proto/todo/v1"
)

func addTask(c pb.TodoServiceClient, description string, dueDate time.Time) uint64 {
	req := &pb.AddTaskRequest{
		Description: description,
		DueDate:     timestamppb.New(dueDate),
	}

	res, err := c.AddTask(context.Background(), req)
	if err != nil {
		panic(err)
	}
	log.Printf("Client: added task with id: %d", res.Id)
	return res.Id
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatalln("usage: client [IP_ADDR]")
	}

	addr := args[0]
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(addr, opts...)

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// add task
	c := pb.NewTodoServiceClient(conn)
	fmt.Println("------- ADD TASK -------")
	dueDate := time.Now().Add(time.Hour * 24)
	addTask(c, "Buy milk", dueDate)
	fmt.Println("------- END TASK -------")

	// list tasks
	fmt.Println("------- LIST TASKS -------")
	printTasks(c)
	fmt.Println("------- END LIST TASKS -------")

	// update task
	fmt.Println("------- UPDATE TASK -------")
	id1 := addTask(c, "Buy bread", dueDate)
	id2 := addTask(c, "Buy eggs", dueDate)

	taskToUpdate := []*pb.UpdateTaskRequest{
		{Task: &pb.Task{Id: id1, Done: false, Description: "Buy 2 bread", DueDate: timestamppb.New(dueDate)}},
		{Task: &pb.Task{Id: id2, Done: true, Description: "Buy 3 eggs", DueDate: timestamppb.New(dueDate)}},
	}

	updateTask(c, taskToUpdate...)

	printTasks(c)
	fmt.Println("------- END UPDATE TASK -------")

	// delete task
	fmt.Println("------- DELETE TASK -------")
	deleteTasks := []*pb.DeleteTaskRequest{
		{Id: id1},
		{Id: id2},
	}
	deleteTask(c, deleteTasks...)
	printTasks(c)
	fmt.Println("------- END DELETE TASK -------")

	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.Fatalf("unexpected error: %v", err)
		}
	}(conn)
}

func printTasks(c pb.TodoServiceClient) {

	req := &pb.ListTasksRequest{}

	stream, err := c.ListTasks(context.Background(), req)
	if err != nil {
		log.Fatalf("error listing tasks: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error receiving tasks: %v", err)
		}

		fmt.Println(res.Task.String(), "overdue:", res.Overdue)
	}
}

func updateTask(
	c pb.TodoServiceClient,
	reqs ...*pb.UpdateTaskRequest,
) {
	stream, err := c.UpdateTask(context.Background())

	if err != nil {
		log.Fatalf("error updating task: %v", err)
	}

	for _, req := range reqs {
		if err := stream.Send(req); err != nil {
			log.Fatalf("error sending task: %v", err)
		}
		if req.Task != nil {
			log.Printf("Client: updated task with id: %d", req.Task.Id)
		}
	}

	if _, err := stream.CloseAndRecv(); err != nil {
		log.Fatalf("error closing stream: %v", err)
	}
}

func deleteTask(c pb.TodoServiceClient, reqs ...*pb.DeleteTaskRequest) {
	stream, err := c.DeleteTask(context.Background())

	if err != nil {
		log.Fatalf("error deleting task: %v", err)
	}

	// wait for response
	waitc := make(chan struct{})

	// Handle responses
	go func() {
		for {
			_, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("error receiving task: %v", err)
			}

			log.Println("Client: deleted task")
		}
	}()

	// Send requests
	for _, req := range reqs {
		if err := stream.Send(req); err != nil {
			log.Fatalf("error sending task: %v", err)
		}
	}

	// Close stream
	if err := stream.CloseSend(); err != nil {
		log.Fatalf("error closing stream: %v", err)
	}

	// Wait for response
	<-waitc
}
