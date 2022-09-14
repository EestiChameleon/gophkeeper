package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/EestiChameleon/gophkeeper/server/cfg"
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

// InitGRPCServer initializes a new gRPC server.
func InitGRPCServer() (*GRPCServer, error) {
	// creates a gRPC server
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.AuthCheckGRPC))
	// register the service
	pb.RegisterKeeperServer(s, &GRPCServer{})

	return &GRPCServer{serv: s}, nil
}

// Start launch the server.
func (g *GRPCServer) Start() error {
	// determines the server port
	listen, err := net.Listen("tcp", cfg.ServerAddress)
	if err != nil {
		return err
	}

	fmt.Println("Сервер gRPC начал работу")
	// listen for gRPC requests
	return g.serv.Serve(listen)
}

// ShutDown graceful stops the server.
func (g *GRPCServer) ShutDown() error {
	g.serv.GracefulStop()
	return nil
}

// RegisterUser handler creates new user. Returns JWT with userID encoded.
func (g *GRPCServer) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	if in.ServiceLogin == `` || in.ServicePass == `` {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	usrID, err := service.UserAdd(in.ServiceLogin, in.ServicePass)
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

// LoginUser authenticates the user. Provides the latest data from database and JWT.
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

	return &pb.LoginUserResponse{
		Status: "login successful",
		Jwt:    token,
	}, nil
}

// GetPair handler returns the found by title pair data.
func (g *GRPCServer) GetPair(ctx context.Context, in *pb.GetPairRequest) (*pb.GetPairResponse, error) {
	if in.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	data, err := service.PairByTitle(in.Title, ctxfunc.GetUserIDFromCTX(ctx))
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
		Pairs: &pb.Pair{
			Title:   data.Title,
			Login:   data.Login,
			Pass:    data.Pass,
			Comment: data.Comment,
			Version: data.Version,
		},
		Status: "success",
	}, nil
}

// PostPair handler checks and saves new pair data to database.
// Handler verifies, if the provided pair data is the latest version and allows to proceed further.
func (g *GRPCServer) PostPair(ctx context.Context, in *pb.PostPairRequest) (*pb.PostPairResponse, error) {
	if in.Pair.Version < 0 || in.Pair.Login == `` || in.Pair.Title == `` || in.Pair.Pass == `` {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	// first - compare version in DB

	dbPair, err := service.PairByTitle(in.Pair.Title, ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		log.Println(err)
		return nil, status.Error(codes.Internal, failedDBQuery)
	}
	// in case not found - we can save passed version. No overwriting data.
	if errors.Is(err, storage.ErrNotFound) {
		err = service.PairAdd(ctxfunc.GetUserIDFromCTX(ctx), in.Pair.Title, in.Pair.Login, in.Pair.Pass, in.Pair.Comment, in.Pair.Version)
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
	err = service.PairAdd(ctxfunc.GetUserIDFromCTX(ctx), in.Pair.Title, in.Pair.Login, in.Pair.Pass, in.Pair.Comment, in.Pair.Version)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, failedToSaveNewVersion)
	}
	// record saved (the latest version)
	return &pb.PostPairResponse{Status: "success"}, nil
}

// DelPair handler deletes the provided pair data by title.
func (g *GRPCServer) DelPair(ctx context.Context, in *pb.DelPairRequest) (*pb.DelPairResponse, error) {
	if in.Title == `` {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	if err := service.PairDelete(in.Title, ctxfunc.GetUserIDFromCTX(ctx)); err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Delete failed. Please try again")
	}

	return &pb.DelPairResponse{Status: "success"}, nil
}

// GetText handler returns the found by title text data.
func (g *GRPCServer) GetText(ctx context.Context, in *pb.GetTextRequest) (*pb.GetTextResponse, error) {
	if in.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	data, err := service.TextByTitle(in.Title, ctxfunc.GetUserIDFromCTX(ctx))
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
		Text: &pb.Text{
			Title:   data.Title,
			Body:    data.Body,
			Comment: data.Comment,
			Version: data.Version,
		},
		Status: "success",
	}, nil
}

// PostText handler checks and saves new text data to database.
// Handler verifies, if the provided text data is the latest version and allows to proceed further.
func (g *GRPCServer) PostText(ctx context.Context, in *pb.PostTextRequest) (*pb.PostTextResponse, error) {
	if in.Text.Version < 0 || in.Text.Title == `` || in.Text.Body == `` {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	// first - compare version in DB
	dbText, err := service.TextByTitle(in.Text.Title, ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		log.Println(err)
		return nil, status.Error(codes.Internal, failedDBQuery)
	}
	// in case not found - we can save version 1? or passed version?
	if errors.Is(err, storage.ErrNotFound) {
		err = service.TextAdd(ctxfunc.GetUserIDFromCTX(ctx), in.Text.Title, in.Text.Body, in.Text.Comment, in.Text.Version)
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
	err = service.TextAdd(ctxfunc.GetUserIDFromCTX(ctx), in.Text.Title, in.Text.Body, in.Text.Comment, in.Text.Version)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, failedToSaveNewVersion)
	}
	// record saved (the latest version)
	return &pb.PostTextResponse{Status: "success"}, nil
}

// DelText handler deletes the provided text data by title.
func (g *GRPCServer) DelText(ctx context.Context, in *pb.DelTextRequest) (*pb.DelTextResponse, error) {
	if in.Title == `` {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	err := service.TextDelete(in.Title, ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Delete failed. Please try again")
	}

	return &pb.DelTextResponse{Status: "success"}, nil
}

// GetBin handler returns the found by title binary data.
func (g *GRPCServer) GetBin(ctx context.Context, in *pb.GetBinRequest) (*pb.GetBinResponse, error) {
	if in.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	data, err := service.BinByTitle(in.Title, ctxfunc.GetUserIDFromCTX(ctx))
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
		BinData: &pb.Bin{
			Title:   data.Title,
			Body:    data.Body,
			Comment: data.Comment,
			Version: data.Version,
		},
		Status: "success",
	}, nil
}

// PostBin handler checks and saves new binary data to database.
// Handler verifies, if the provided binary data is the latest version and allows to proceed further.
func (g *GRPCServer) PostBin(ctx context.Context, in *pb.PostBinRequest) (*pb.PostBinResponse, error) {
	if in.BinData.Version < 0 || in.BinData.Title == `` || len(in.BinData.Body) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	// first - compare version in DB
	dbBin, err := service.BinByTitle(in.BinData.Title, ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		log.Println(err)
		return nil, status.Error(codes.Internal, failedDBQuery)
	}
	// in case not found - we save passed version = 1
	if errors.Is(err, storage.ErrNotFound) {
		err = service.BinAdd(ctxfunc.GetUserIDFromCTX(ctx), in.BinData.Title, in.BinData.Body, in.BinData.Comment, in.BinData.Version)
		if err != nil {
			log.Println(err)
			return nil, status.Error(codes.Internal, failedToSaveNewVersion)
		}
		// new record saved
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
	err = service.BinAdd(ctxfunc.GetUserIDFromCTX(ctx), in.BinData.Title, in.BinData.Body, in.BinData.Comment, in.BinData.Version)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, failedToSaveNewVersion)
	}
	// record saved (the latest version)
	return &pb.PostBinResponse{Status: "success"}, nil
}

// DelBin handler deletes the provided binary data by title.
func (g *GRPCServer) DelBin(ctx context.Context, in *pb.DelBinRequest) (*pb.DelBinResponse, error) {
	if in.Title == `` {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	err := service.BinDelete(in.Title, ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Delete failed. Please try again")
	}

	return &pb.DelBinResponse{Status: "success"}, nil
}

// GetCard handler returns the found by title card data.
func (g *GRPCServer) GetCard(ctx context.Context, in *pb.GetCardRequest) (*pb.GetCardResponse, error) {
	if in.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	data, err := service.CardByTitle(in.Title, ctxfunc.GetUserIDFromCTX(ctx))
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
		Card: &pb.Card{
			Title:   data.Title,
			Number:  data.Number,
			Expdate: data.ExpirationDate,
			Comment: data.Comment,
			Version: data.Version,
		},
		Status: "success",
	}, nil
}

// PostCard handler checks and saves new card data to database.
// Handler verifies, if the provided card data is the latest version and allows to proceed further.
func (g *GRPCServer) PostCard(ctx context.Context, in *pb.PostCardRequest) (*pb.PostCardResponse, error) {
	if in.Card.Version < 0 || in.Card.Title == `` || in.Card.Number == `` || in.Card.Expdate == `` {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	// first - compare version in DB
	dbCard, err := service.CardByTitle(in.Card.Title, ctxfunc.GetUserIDFromCTX(ctx))
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		log.Println(err)
		return nil, status.Error(codes.Internal, failedDBQuery)
	}

	// in case not found - we can save version 1? or passed version?
	if errors.Is(err, storage.ErrNotFound) {
		err = service.CardAdd(ctxfunc.GetUserIDFromCTX(ctx), in.Card.Title, in.Card.Number, in.Card.Expdate, in.Card.Comment, in.Card.Version)
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
	err = service.CardAdd(ctxfunc.GetUserIDFromCTX(ctx), in.Card.Title, in.Card.Number, in.Card.Expdate, in.Card.Comment, in.Card.Version)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, failedToSaveNewVersion)
	}
	// record saved (the latest version)
	return &pb.PostCardResponse{Status: "success"}, nil
}

// DelCard handler deletes the provided card data by title.
func (g *GRPCServer) DelCard(ctx context.Context, in *pb.DelCardRequest) (*pb.DelCardResponse, error) {
	if in.Title == `` {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}

	err := service.CardDelete(in.Title, ctxfunc.GetUserIDFromCTX(ctx))
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
