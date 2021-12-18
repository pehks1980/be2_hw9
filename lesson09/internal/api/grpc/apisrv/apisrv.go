package apisrv

import (
	"context"
	"fmt"
	"io"
	"lesson09/internal/api/grpc/myservice"
	"log"
	"time"
)

var _ myservice.MessageServiceServer = &MyServer{}

type MyServer struct {
	// mainCtx context.Context

	myservice.UnimplementedMessageServiceServer
}

// SendMessage - обработчик входящего сообщения сервиса SendMessage
func (*MyServer) SendMessage(ctx context.Context, msg *myservice.Message) (*myservice.Reply, error) {
	log.Printf("receive message: %s", msg)
	return &myservice.Reply{Id: msg.Id, Status: "OK"}, nil
}
// ServeMessage - обработчик входящего сообщения сервиса ServeMessage класс stream
func (*MyServer) ServeMessage(msg *myservice.Message, stream myservice.MessageService_ServeMessageServer) error {
	log.Printf("receive message: %s", msg)

	// send periodic message to client
	n := 1
	for {
		log.Printf("send message to client")
		time.Sleep(1 * time.Second)

		err := stream.Send(
			&myservice.Message{
				Id:   fmt.Sprint(n),
				Body: "my labs",
			},
		)
		if err != nil {
			log.Printf("error sending message to client, %s", err)
			break
		}
		n++
	}

	return nil
}

func (s *MyServer) StreamMessage(stream myservice.MessageService_StreamMessageServer) error {
	// read client message
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				log.Printf("client write closed, %s", err)
				break
			}

			if err != nil && err != io.EOF {
				log.Printf("error reading from client, %s", err)
				break
			}

			log.Printf("message from client, %s", msg)
		}
	}()

	// send periodic message to client
	n := 1
	for {
		// select {
		// case <-s.mainCtx.Done():
		// 	return fmt.Errorf("ctx break")
		// default:
		// }

		log.Printf("send message to client")
		time.Sleep(1 * time.Second)

		err := stream.Send(
			&myservice.Message{
				Id:   fmt.Sprint(n),
				Body: "my labs",
			},
		)
		if err != nil {
			log.Printf("error sending message to client, %s", err)
			break
		}
		n++
	}

	return nil
}
