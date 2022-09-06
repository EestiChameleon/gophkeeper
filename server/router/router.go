package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/EestiChameleon/gophkeeper/models"
	"github.com/EestiChameleon/gophkeeper/server/ctxfunc"
	"github.com/EestiChameleon/gophkeeper/server/router/interceptors"
	"github.com/EestiChameleon/gophkeeper/server/service"
	"github.com/EestiChameleon/gophkeeper/server/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	// импортируем пакет со сгенерированными protobuf-файлами
	pb "github.com/EestiChameleon/gophkeeper/proto"
)

var (
	newerVersionDetected   = "Current / newer version found in database. Please synchronize you app to get the most actual data."
	failedToSaveNewVersion = "Failed to save new record. Please try again."
	failedDBQuery          = "failed to obtain database data"
)

type GRPCServer struct {
	serv *grpc.Server
	pb.UnimplementedKeeperServer
}

func InitGRPCServer() (*GRPCServer, error) {
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer(
		//grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(interceptors.AuthCheck)))
		grpc.UnaryInterceptor(interceptors.AuthCheckGRPC))
	// регистрируем сервис
	pb.RegisterKeeperServer(s, &GRPCServer{})

	return &GRPCServer{serv: s}, nil
}

func (g *GRPCServer) Start() error {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		return err
	}

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	return g.serv.Serve(listen)
}

func (g *GRPCServer) ShutDown() error {
	g.serv.GracefulStop()
	return nil
}

// RegisterUser .
func (g *GRPCServer) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	if in.ServiceLogin == `` || in.ServicePass == `` {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	encryptPass := service.EncryptPass(in.ServicePass)
	var usrID int
	err := storage.GetSingleValue("users_add", &usrID, in.ServiceLogin, encryptPass)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "failed to register new user")
	}

	userJWT, err := service.JWTEncodeUserID(usrID)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "failed to create jwt")
	}

	return &pb.RegisterUserResponse{
		Status: "registered",
		Jwt:    userJWT,
	}, nil
}

func (g *GRPCServer) LoginUser(ctx context.Context, in *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	if in.ServiceLogin == "" || in.ServicePass == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	token, err := service.CheckAuthData(service.LoginData{
		Login:    in.ServiceLogin,
		Password: in.ServicePass,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrWrongAuthData) || errors.Is(err, storage.ErrNotFound):
			return nil, status.Error(codes.Unauthenticated, "access denied")
		default:
			log.Println(err)
			return nil, status.Error(codes.Internal, "failed to process login/password")
		}
	}

	data, err := storage.GetAllUserDataLastVersion(ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "failed to obtain latest data")
	}

	return &pb.LoginUserResponse{
		AllData: &pb.SyncVaultResponse{
			Pairs:   data.Pairs,
			Texts:   data.Texts,
			BinData: data.Bins,
			Cards:   data.Cards,
			Status:  "success",
		},
		Status: "login successful",
		Jwt:    token,
	}, nil
}

func (g *GRPCServer) GetPair(ctx context.Context, in *pb.GetPairRequest) (*pb.GetPairResponse, error) {
	if in.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	data := new(pb.Pair)
	err := storage.GetOneRow("pair_by_title", data, in.Title, ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &pb.GetPairResponse{
				Pairs:  nil,
				Status: "not found",
			}, status.Error(codes.NotFound, "not found")
		}
		log.Println(err)
		return nil, status.Error(codes.Internal, "failed to obtain latest data")
	}

	return &pb.GetPairResponse{
		Pairs:  data,
		Status: "success",
	}, nil
}

func (g *GRPCServer) PostPair(ctx context.Context, in *pb.PostPairRequest) (*pb.PostPairResponse, error) {
	if in.Pair.Version < 0 || in.Pair.Login == `` || in.Pair.Title == `` {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	// first - compare version in DB
	dbPair := new(models.Pair)
	err := storage.GetOneRow("pair_by_title", dbPair, in.Pair.Title, ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		log.Println(err)
		return nil, status.Error(codes.Internal, failedDBQuery)
	}
	// in case not found - we can save version 1? or passed version?
	if errors.Is(err, storage.ErrNotFound) {
		var resultID int
		err = storage.GetSingleValue("pair_add", &resultID,
			ctxfunc.GetUserIDFromCTX(ctx), in.Pair.Title, in.Pair.Login, in.Pair.Pass, in.Pair.Comment, in.Pair.Version)
		if err != nil {
			log.Println(err)
			return nil, status.Error(codes.Internal, failedToSaveNewVersion)
		}
		// new record saved (first version?)
		return &pb.PostPairResponse{Status: "success"}, nil
	}

	// case when we found a version in DB
	// check for deleted in DB
	if !dbPair.DeletedAt.Valid { // not deleted
		// compare versions:
		if in.Pair.Version <= dbPair.Version {
			// DB has actual or newer version - error and ask customer to sync
			return nil, status.Error(codes.AlreadyExists, newerVersionDetected)
		}
	}

	// db version is deleted OR received version is the latest => save
	var resultID int
	err = storage.GetSingleValue("pair_add", &resultID,
		ctxfunc.GetUserIDFromCTX(ctx), in.Pair.Title, in.Pair.Login, in.Pair.Pass, in.Pair.Comment, in.Pair.Version)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, failedToSaveNewVersion)
	}
	// record saved (the latest version)
	return &pb.PostPairResponse{Status: "success"}, nil
}

func (g *GRPCServer) DelPair(ctx context.Context, in *pb.DelPairRequest) (*pb.DelPairResponse, error) {
	if in.Title == `` {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	var affectedRows int // just to put something in return param in function
	err := storage.GetSingleValue("pair_del_by_title", &affectedRows, in.Title, ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Delete failed. Please try again")
	}

	return &pb.DelPairResponse{Status: "success"}, nil
}

func (g *GRPCServer) GetText(ctx context.Context, in *pb.GetTextRequest) (*pb.GetTextResponse, error) {
	if in.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	data := new(pb.Text)
	err := storage.GetOneRow("text_by_title", data, in.Title, ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &pb.GetTextResponse{
				Text:   nil,
				Status: "not found",
			}, status.Error(codes.NotFound, "not found")
		}
		log.Println(err)
		return nil, status.Error(codes.Internal, "failed to obtain latest data")
	}

	return &pb.GetTextResponse{
		Text:   data,
		Status: "success",
	}, nil
}

func (g *GRPCServer) PostText(ctx context.Context, in *pb.PostTextRequest) (*pb.PostTextResponse, error) {
	if in.Text.Version < 0 || in.Text.Title == `` || in.Text.Body == `` {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	// first - compare version in DB
	dbText := new(models.Text)
	err := storage.GetOneRow("text_by_title", dbText, in.Text.Title, ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		log.Println(err)
		return nil, status.Error(codes.Internal, failedDBQuery)
	}
	// in case not found - we can save version 1? or passed version?
	if errors.Is(err, storage.ErrNotFound) {
		var resultID int
		err = storage.GetSingleValue("text_add", &resultID,
			ctxfunc.GetUserIDFromCTX(ctx), in.Text.Title, in.Text.Body, in.Text.Comment, in.Text.Version)
		if err != nil {
			log.Println(err)
			return nil, status.Error(codes.Internal, failedToSaveNewVersion)
		}
		// new record saved (first version?)
		return &pb.PostTextResponse{Status: "success"}, nil
	}

	// case when we found a version in DB
	// check for deleted in DB
	if !dbText.DeletedAt.Valid { // not deleted
		// compare versions:
		if in.Text.Version <= dbText.Version {
			// DB has actual or newer version - error and ask customer to sync
			return nil, status.Error(codes.AlreadyExists, newerVersionDetected)
		}
	}

	// db version is deleted OR received version is the latest => save
	var resultID int
	err = storage.GetSingleValue("pair_add", &resultID,
		ctxfunc.GetUserIDFromCTX(ctx), in.Text.Title, in.Text.Body, in.Text.Comment, in.Text.Version)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, failedToSaveNewVersion)
	}
	// record saved (the latest version)
	return &pb.PostTextResponse{Status: "success"}, nil
}

func (g *GRPCServer) DelText(ctx context.Context, in *pb.DelTextRequest) (*pb.DelTextResponse, error) {
	if in.Title == `` {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	var affectedRows int // just to put something in return param in function
	err := storage.GetSingleValue("text_del_by_title", &affectedRows, in.Title, ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Delete failed. Please try again")
	}

	return &pb.DelTextResponse{Status: "success"}, nil
}

func (g *GRPCServer) GetBin(ctx context.Context, in *pb.GetBinRequest) (*pb.GetBinResponse, error) {
	if in.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	data := new(pb.Bin)
	err := storage.GetOneRow("bin_by_title", data, in.Title, ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &pb.GetBinResponse{
				BinData: nil,
				Status:  "not found",
			}, status.Error(codes.NotFound, "not found")
		}
		log.Println(err)
		return nil, status.Error(codes.Internal, "failed to obtain latest data")
	}

	return &pb.GetBinResponse{
		BinData: data,
		Status:  "success",
	}, nil
}

func (g *GRPCServer) PostBin(ctx context.Context, in *pb.PostBinRequest) (*pb.PostBinResponse, error) {
	if in.BinData.Version < 0 || in.BinData.Title == `` || len(in.BinData.Body) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	// first - compare version in DB
	dbBin := new(models.Bin)
	err := storage.GetOneRow("bin_by_title", dbBin, in.BinData.Title, ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		log.Println(err)
		return nil, status.Error(codes.Internal, failedDBQuery)
	}
	// in case not found - we can save version 1? or passed version?
	if errors.Is(err, storage.ErrNotFound) {
		var resultID int
		err = storage.GetSingleValue("bin_add", &resultID,
			ctxfunc.GetUserIDFromCTX(ctx), in.BinData.Title, in.BinData.Body, in.BinData.Comment, in.BinData.Version)
		if err != nil {
			log.Println(err)
			return nil, status.Error(codes.Internal, failedToSaveNewVersion)
		}
		// new record saved (first version?)
		return &pb.PostBinResponse{Status: "success"}, nil
	}

	// case when we found a version in DB
	// check for deleted in DB
	if !dbBin.DeletedAt.Valid { // not deleted
		// compare versions:
		if in.BinData.Version <= dbBin.Version {
			// DB has actual or newer version - error and ask customer to sync
			return nil, status.Error(codes.AlreadyExists, newerVersionDetected)
		}
	}

	// db version is deleted OR received version is the latest => save
	var resultID int
	err = storage.GetSingleValue("bin_add", &resultID,
		ctxfunc.GetUserIDFromCTX(ctx), in.BinData.Title, in.BinData.Body, in.BinData.Comment, in.BinData.Version)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, failedToSaveNewVersion)
	}
	// record saved (the latest version)
	return &pb.PostBinResponse{Status: "success"}, nil
}

func (g *GRPCServer) DelBin(ctx context.Context, in *pb.DelBinRequest) (*pb.DelBinResponse, error) {
	if in.Title == `` {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	var affectedRows int // just to put something in return param in function
	err := storage.GetSingleValue("bin_del_by_title", &affectedRows, in.Title, ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Delete failed. Please try again")
	}

	return &pb.DelBinResponse{Status: "success"}, nil
}

func (g *GRPCServer) GetCard(ctx context.Context, in *pb.GetCardRequest) (*pb.GetCardResponse, error) {
	if in.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	data := new(pb.Card)
	err := storage.GetOneRow("card_by_title", data, in.Title, ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return &pb.GetCardResponse{
				Card:   nil,
				Status: "not found",
			}, status.Error(codes.NotFound, "not found")
		}
		log.Println(err)
		return nil, status.Error(codes.Internal, "failed to obtain latest data")
	}

	return &pb.GetCardResponse{
		Card:   data,
		Status: "success",
	}, nil
}

func (g *GRPCServer) PostCard(ctx context.Context, in *pb.PostCardRequest) (*pb.PostCardResponse, error) {
	if in.Card.Version < 0 || in.Card.Title == `` || in.Card.Number == `` || in.Card.Expdate == `` {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	// first - compare version in DB
	dbCard := new(models.Card)
	err := storage.GetOneRow("card_by_title", dbCard, in.Card.Title, ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		log.Println(err)
		return nil, status.Error(codes.Internal, failedDBQuery)
	}
	// in case not found - we can save version 1? or passed version?
	if errors.Is(err, storage.ErrNotFound) {
		var resultID int
		err = storage.GetSingleValue("card_add", &resultID,
			ctxfunc.GetUserIDFromCTX(ctx), in.Card.Title, in.Card.Number, in.Card.Expdate, in.Card.Comment, in.Card.Version)
		if err != nil {
			log.Println(err)
			return nil, status.Error(codes.Internal, failedToSaveNewVersion)
		}
		// new record saved (first version?)
		return &pb.PostCardResponse{Status: "success"}, nil
	}

	// case when we found a version in DB
	// check for deleted in DB
	if !dbCard.DeletedAt.Valid { // not deleted
		// compare versions:
		if in.Card.Version <= dbCard.Version {
			// DB has actual or newer version - error and ask customer to sync
			return nil, status.Error(codes.AlreadyExists, newerVersionDetected)
		}
	}

	// db version is deleted OR received version is the latest => save
	var resultID int
	err = storage.GetSingleValue("card_add", &resultID,
		ctxfunc.GetUserIDFromCTX(ctx), in.Card.Title, in.Card.Number, in.Card.Expdate, in.Card.Comment, in.Card.Version)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, failedToSaveNewVersion)
	}
	// record saved (the latest version)
	return &pb.PostCardResponse{Status: "success"}, nil
}

func (g *GRPCServer) DelCard(ctx context.Context, in *pb.DelCardRequest) (*pb.DelCardResponse, error) {
	if in.Title == `` {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	var affectedRows int // just to put something in return param in function
	err := storage.GetSingleValue("card_del_by_title", &affectedRows, in.Title, ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Delete failed. Please try again")
	}

	return &pb.DelCardResponse{Status: "success"}, nil
}

func (g *GRPCServer) SyncVault(ctx context.Context, in *pb.SyncVaultRequest) (*pb.SyncVaultResponse, error) {
	data, err := storage.GetAllUserDataLastVersion(ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "failed to obtain latest data")
	}

	return &pb.SyncVaultResponse{
		Pairs:   data.Pairs,
		Texts:   data.Texts,
		BinData: data.Bins,
		Cards:   data.Cards,
		Status:  "success",
	}, nil
}
