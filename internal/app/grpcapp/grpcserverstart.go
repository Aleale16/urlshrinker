package grpcapp

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "urlshrinker/internal/app/grpcapp/proto"

	"google.golang.org/grpc"
)

func Grpcserverstart() error {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}

	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()
	// регистрируем сервис
	pb.RegisterActionsServer(s, &ActionsServer{})

	//	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	//	if err := s.Serve(listen); err != nil {
	//		log.Fatal(err)
	//	}

	errChan := make(chan error)
	stopChan := make(chan os.Signal)

	// Ожидаем события от ОС
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// сообщаем об ошибках в канал
	go func() {
		if err := s.Serve(listen); err != nil {
			errChan <- err
		}
	}()

	// у GRPC оказывается есть GracefulStop?
	defer func() {
		s.GracefulStop()
	}()

	select {
	case err := <-errChan:
		log.Printf("Fatal error: %v\n", err)
	case <-stopChan:
	}
	return nil
}
