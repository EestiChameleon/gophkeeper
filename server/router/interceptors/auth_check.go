package interceptors

import (
	"context"
	"github.com/EestiChameleon/gophkeeper/server/ctxfunc"
	"github.com/EestiChameleon/gophkeeper/server/service"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

var (
	SkipCheckMethods = map[string]struct{}{
		"/gophkeeper.proto.Keeper/RegisterUser": {},
		"/gophkeeper.proto.Keeper/LoginUser":    {},
	}
)

// AuthCheck is an interceptor to authenticate requests
func AuthCheck(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	userID, err := service.JWTDecodeUserID(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	return ctxfunc.SetUserIDToCTX(ctx, userID), nil
}

// AuthCheckGRPC
func AuthCheckGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Println("--> unary interceptor: ", info.FullMethod)
	// check for method, which doesn't need to be intercepted
	_, ok := SkipCheckMethods[info.FullMethod]
	if ok {
		return handler(ctx, req)
	}
	// check part
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	userID, err := service.JWTDecodeUserID(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	return handler(ctxfunc.SetUserIDToCTX(ctx, userID), req)
}
