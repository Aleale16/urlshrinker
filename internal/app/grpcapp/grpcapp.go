package grpcapp

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

	response.FullURL, response.Status = storage.Getrecord(in.ShortURL)

	return &response, nil
}

// PostRecord реализует интерфейс сохранения полного урла в базу и получение его короткого айди.
func (s *ActionsServer) PostRecord(ctx context.Context, in *pb.PostRecordRequest) (*pb.PostRecordResponse, error) {
	var response pb.PostRecordResponse

	response.ShortURL, response.Status = storage.Storerecord(in.FullURL)

	return &response, nil
}

// Getuserrecords реализует интерфейс получения полных урлов для указанного айди пользователя.
func (s *ActionsServer) Getuserrecords(ctx context.Context, in *pb.GetuserrecordsRequest) (*pb.GetuserrecordsResponse, error) {
	var response pb.GetuserrecordsResponse

	response.FullURLs, response.NoURLs, response.ArrayUserURLs = storage.GetuserURLS(in.Userid)

	return &response, nil
}

// PostShortURLtouser реализует интерфейс создания связки короткого айди урла с айди пользователя его создавшим.
func (s *ActionsServer) PostShortURLtouser(ctx context.Context, in *pb.PostShortURLtouserRequest) error {

	storage.AssignShortURLtouser(in.Userid, in.ShortURLid)

	return nil
}

// CheckPGdbConn реализует интерфейс проверки соединения с БД.
func (s *ActionsServer) CheckPGdbConn(ctx context.Context, in *pb.CheckPGdbConnRequest) (*pb.CheckPGdbConnResponse, error) {
	var response pb.CheckPGdbConnResponse

	response.Isconnected = storage.CheckPGdbConn

	return &response, nil
}
