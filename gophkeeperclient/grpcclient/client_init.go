package grpcclient

import (
	"github.com/EestiChameleon/gophkeeper/gophkeeperclient/cfg"
	pb "github.com/EestiChameleon/gophkeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var clientConn *grpc.ClientConn

// DialUp initiates a connection between the client and the server. Address taken from cfg.GrpcServerPath.
func DialUp() (pb.KeeperClient, error) {
	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(cfg.GrpcServerPath, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	clientConn = conn

	// получаем переменную интерфейсного типа UsersClient,
	// через которую будем отправлять сообщения
	return pb.NewKeeperClient(conn), nil
}

// ConnDown closes the server connection.
func ConnDown() error {
	return clientConn.Close()
}

// ActiveConnection verifies if there is an active connection. Returns true in case of any active connection found.
func ActiveConnection() bool {
	return clientConn != nil
}
