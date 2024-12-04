package main

import (
	"context"
	"io"
	"slices"
	"time"

	pb "github.com/araminian/grpcch4/proto/todo/v2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func (s *server) AddTask(
	_ context.Context,
	in *pb.AddTaskRequest,
) (*pb.AddTaskResponse, error) {
	id, _ := s.d.addTask(in.Description, in.DueDate.AsTime())
	return &pb.AddTaskResponse{Id: id}, nil
}

func (s *server) ListTasks(
	req *pb.ListTasksRequest,
	stream pb.TodoService_ListTasksServer,
) error {
	return s.d.getTasks(func(i interface{}) error {
		task := i.(*pb.Task)

		// filter
		Filter(task, req.Mask)

		overdue := task.DueDate != nil && !task.Done && task.DueDate.AsTime().Before(time.Now().UTC())
		err := stream.Send(&pb.ListTasksResponse{
			Task:    task,
			Overdue: overdue,
		})
		return err
	})
}

func (s *server) UpdateTask(stream pb.TodoService_UpdateTaskServer) error {

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.UpdateTaskResponse{})
		}
		if err != nil {
			return err
		}

		s.d.updateTask(req.Id, req.Description, req.DueDate.AsTime(), req.Done)
	}
}

func (s *server) DeleteTask(stream pb.TodoService_DeleteTaskServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		s.d.deleteTask(req.Id)
		stream.Send(&pb.DeleteTaskResponse{})
	}
}

func Filter(msg proto.Message, mask *fieldmaskpb.FieldMask) {

	if mask == nil || len(mask.Paths) == 0 {
		return
	}

	rft := msg.ProtoReflect()

	rft.Range(func(fd protoreflect.FieldDescriptor, _ protoreflect.Value) bool {
		if !slices.Contains(mask.Paths, string(fd.Name())) {
			rft.Clear(fd)
		}
		return true
	})

}
