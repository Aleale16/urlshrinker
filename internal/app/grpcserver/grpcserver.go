package grpcserver

import (
	// импортируем пакет со сгенерированными protobuf-файлами
	"context"
	pb "urlshrinker/internal/app/grpcapp/proto"
	"urlshrinker/internal/app/storage"
)

// ActionsServer поддерживает все необходимые методы сервера.
type ActionsServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedActionsServer
}

// GetRecord реализует интерфейс получения полного урла по его короткому айди.
func (s *ActionsServer) GetRecord(ctx context.Context, in *pb.GetRecordRequest) (*pb.GetRecordResponse, error) {
	var response pb.GetRecordResponse

	response.fullURL, response.status = storage.Getrecord(in.shortURL)

	return &response, nil
}

// CheckPGdbConn реализует интерфейс проверки соединения с БД.
func (s *ActionsServer) CheckPGdbConn(ctx context.Context, in *pb.CheckPGdbConnRequest) (*pb.CheckPGdbConnResponse, error) {
	var response pb.CheckPGdbConnResponse

	response.connected = storage.CheckPGdbConn

	return &response, nil
}
