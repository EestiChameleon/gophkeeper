package grpcserver

import (
	"context"
	"github.com/EestiChameleon/gophkeeper/models"
	pb "github.com/EestiChameleon/gophkeeper/proto"
	"github.com/EestiChameleon/gophkeeper/server/service"
	"github.com/EestiChameleon/gophkeeper/server/storage"
	"github.com/EestiChameleon/gophkeeper/server/storage/testdb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

const bufSize = 1024 * 1024

var (
	lis *bufconn.Listener
)

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterKeeperServer(s, &GRPCServer{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Test server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

// keeperTestClient creates a client side for test server.
func keeperTestConn(t *testing.T) (context.Context, *grpc.ClientConn) {
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	return ctx, conn
}

// TestRegisterUser verifies, that:
// 1) null values are not accepted
// 2) in case of success - we receive a valid JWT
func TestRegisterUser(t *testing.T) {
	// init connect
	ctx, conn := keeperTestConn(t)
	defer conn.Close()

	// create client
	client := pb.NewKeeperClient(conn)
	// init test storage
	storage.InitTest()
	type want struct {
		errStatusCode codes.Code
		errStatusMsg  string
		jwtUserID     int
	}

	tests := []struct {
		name     string
		number   uint8
		login    string
		password string
		want     want
	}{
		{
			name:     "Test #1: empty data",
			number:   1,
			login:    "",
			password: "",
			want: want{
				errStatusCode: codes.InvalidArgument,
				errStatusMsg:  "invalid argument",
			},
		},
		{
			name:     "Test #2: correct data",
			number:   2,
			login:    "testUser",
			password: "testPassword",
			want: want{
				jwtUserID: 7,
			},
		},
	}
	for _, tt := range tests {
		// make request
		resp, err := client.RegisterUser(ctx,
			&pb.RegisterUserRequest{
				ServiceLogin: tt.login,
				ServicePass:  tt.password,
			},
		)
		switch tt.number {
		case 1:
			// empty values should return err
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 2:
			// check successfully created jwt
			id, err2 := service.JWTDecodeUserID(resp.GetJwt())
			assert.NoError(t, err2)
			assert.Equal(t, tt.want.jwtUserID, id)
		}
	}
}

// TestLoginUser verifies, that:
// 1) null values are not accepted
// 2) passed data is correct - auth check is successful
// 2) in case of success - we receive a valid JWT
func TestLoginUser(t *testing.T) {
	// init connect
	ctx, conn := keeperTestConn(t)
	defer conn.Close()

	// create client
	client := pb.NewKeeperClient(conn)
	// init test storage
	storage.InitTest()
	type want struct {
		errStatusCode codes.Code
		errStatusMsg  string
		jwtUserID     int
	}

	tests := []struct {
		name     string
		number   uint8
		login    string
		password string
		want     want
	}{
		{
			name:     "Test #1: empty data",
			number:   1,
			login:    "",
			password: "",
			want: want{
				errStatusCode: codes.InvalidArgument,
				errStatusMsg:  "invalid argument",
			},
		},
		{
			name:     "Test #2: incorrect data: unknown user",
			number:   2,
			login:    "user1",
			password: "pass1",
			want: want{
				errStatusCode: codes.Unauthenticated,
				errStatusMsg:  "access denied",
			},
		},
		{
			name:     "Test #3: incorrect data: wrong pass",
			number:   3,
			login:    "user7",
			password: "pass1",
			want: want{
				errStatusCode: codes.Unauthenticated,
				errStatusMsg:  "access denied",
			},
		},
		{
			name:     "Test #4: correct data",
			number:   4,
			login:    "user7",
			password: "pass7",
			want: want{
				jwtUserID: 7,
			},
		},
	}
	for _, tt := range tests {
		// make request
		resp, err := client.LoginUser(ctx,
			&pb.LoginUserRequest{
				ServiceLogin: tt.login,
				ServicePass:  tt.password,
			},
		)
		switch tt.number {
		case 1:
			// empty values should return err InvalidArgument
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 2:
			// unknown login return err Unauthenticated
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 3:
			// wrong password return err Unauthenticated
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 4:
			// check successfully created jwt
			id, err2 := service.JWTDecodeUserID(resp.GetJwt())
			assert.NoError(t, err2)
			assert.Equal(t, tt.want.jwtUserID, id)
		}
	}
}

// TestGetPair verifies, that:
// 1) null title is not accepted
// 2) not found is correctly processed
// 2) in case of success - we receive test pair data
func TestGetPair(t *testing.T) {
	// init connect
	ctx, conn := keeperTestConn(t)
	defer conn.Close()

	// create client
	client := pb.NewKeeperClient(conn)
	// init test storage
	storage.InitTest()
	type want struct {
		errStatusCode codes.Code
		errStatusMsg  string
		data          *pb.Pair
		status        string
	}

	tests := []struct {
		name   string
		number uint8
		title  string
		want   want
	}{
		{
			name:   "Test #1: empty data",
			number: 1,
			title:  "",
			want: want{
				errStatusCode: codes.InvalidArgument,
				errStatusMsg:  "invalid argument",
			},
		},
		{
			name:   "Test #2: unknown data",
			number: 2,
			title:  "unknown",
			want: want{
				errStatusCode: codes.NotFound,
				errStatusMsg:  "not found",
			},
		},
		{
			name:   "Test #3: correct data",
			number: 3,
			title:  testdb.TestPair.Title,
			want: want{
				data:   models.ModelsToProtoPair(testdb.TestPair),
				status: "success",
			},
		},
	}
	for _, tt := range tests {
		// make request
		resp, err := client.GetPair(ctx, &pb.GetPairRequest{Title: tt.title})
		switch tt.number {
		case 1:
			// empty values should return err InvalidArgument
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 2:
			// unknown title return err Not Found
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 3:
			// check successful response
			d := resp.GetPairs()
			assert.Equal(t, tt.want.data, d)
			assert.Equal(t, tt.want.status, resp.GetStatus())
		}
	}
}

// TestPostPair verifies, that:
// 1) null data is not accepted
// 2) already exists is correctly processed
// 2) in case of success - we receive success status
func TestPostPair(t *testing.T) {
	// init connect
	ctx, conn := keeperTestConn(t)
	defer conn.Close()

	// create client
	client := pb.NewKeeperClient(conn)
	// init test storage
	storage.InitTest()
	type want struct {
		errStatusCode codes.Code
		errStatusMsg  string
		status        string
	}

	tests := []struct {
		name    string
		number  uint8
		title   string
		login   string
		pass    string
		comment string
		version uint32
		want    want
	}{
		{
			name:   "Test #1: empty data",
			number: 1,
			want: want{
				errStatusCode: codes.InvalidArgument,
				errStatusMsg:  "invalid argument",
			},
		},
		{
			name:    "Test #2: data already exists",
			number:  2,
			title:   testdb.TestPair.Title,
			login:   testdb.TestPair.Login,
			pass:    testdb.TestPair.Pass,
			comment: testdb.TestPair.Comment,
			version: testdb.TestPair.Version,
			want: want{
				errStatusCode: codes.AlreadyExists,
				errStatusMsg:  "Current / newer version found in database. Please synchronize you app to get the most actual data.",
			},
		},
		{
			name:    "Test #3: correct new data",
			number:  3,
			title:   "new pair",
			login:   "new login",
			pass:    "new pass",
			comment: "new comment",
			version: 1,
			want:    want{status: "success"},
		},
	}
	for _, tt := range tests {
		// make request
		resp, err := client.PostPair(ctx,
			&pb.PostPairRequest{Pair: &pb.Pair{
				Title:   tt.title,
				Login:   tt.login,
				Pass:    tt.pass,
				Comment: tt.comment,
				Version: tt.version,
			}})
		switch tt.number {
		case 1:
			// empty values should return err InvalidArgument
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 2:
			// unknown title return err Not Found
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 3:
			// check successful response
			assert.Equal(t, tt.want.status, resp.GetStatus())
		}
	}
}

// TestDelPair verifies, that:
// 1) null title is not accepted
// 2) not found received after delete
// 3) in case of success - we receive success message
func TestDelPair(t *testing.T) {
	// init connect
	ctx, conn := keeperTestConn(t)
	defer conn.Close()

	// create client
	client := pb.NewKeeperClient(conn)
	// init test storage
	storage.InitTest()
	type want struct {
		errStatusCode codes.Code
		errStatusMsg  string
		status        string
	}

	tests := []struct {
		name   string
		number uint8
		title  string
		want   want
	}{
		{
			name:   "Test #1: empty data",
			number: 1,
			title:  "",
			want: want{
				errStatusCode: codes.InvalidArgument,
				errStatusMsg:  "invalid argument",
			},
		},
		{
			name:   "Test #3: correct data",
			number: 2,
			title:  testdb.TestPair.Title,
			want: want{
				status:        "success",
				errStatusCode: codes.NotFound,
				errStatusMsg:  "not found",
			},
		},
	}
	for _, tt := range tests {
		// make request
		resp, err := client.DelPair(ctx, &pb.DelPairRequest{Title: tt.title})
		switch tt.number {
		case 1:
			// empty values should return err InvalidArgument
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 2:
			// check successful response of delete
			assert.Equal(t, tt.want.status, resp.GetStatus())
			// make get request to receive not found
			_, err2 := client.GetPair(ctx, &pb.GetPairRequest{Title: tt.title})
			// empty values should return err InvalidArgument
			st, ok := status.FromError(err2)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())

		}
	}
}

// TestGetText verifies, that:
// 1) null title is not accepted
// 2) not found is correctly processed
// 2) in case of success - we receive test text data
func TestGetText(t *testing.T) {
	// init connect
	ctx, conn := keeperTestConn(t)
	defer conn.Close()

	// create client
	client := pb.NewKeeperClient(conn)
	// init test storage
	storage.InitTest()
	type want struct {
		errStatusCode codes.Code
		errStatusMsg  string
		data          *pb.Text
		status        string
	}

	tests := []struct {
		name   string
		number uint8
		title  string
		want   want
	}{
		{
			name:   "Test #1: empty data",
			number: 1,
			title:  "",
			want: want{
				errStatusCode: codes.InvalidArgument,
				errStatusMsg:  "invalid argument",
			},
		},
		{
			name:   "Test #2: unknown data",
			number: 2,
			title:  "unknown",
			want: want{
				errStatusCode: codes.NotFound,
				errStatusMsg:  "not found",
			},
		},
		{
			name:   "Test #3: correct data",
			number: 3,
			title:  testdb.TestText.Title,
			want: want{
				data:   models.ModelsToProtoText(testdb.TestText),
				status: "success",
			},
		},
	}
	for _, tt := range tests {
		// make request
		resp, err := client.GetText(ctx, &pb.GetTextRequest{Title: tt.title})
		switch tt.number {
		case 1:
			// empty values should return err InvalidArgument
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 2:
			// unknown title return err Not Found
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 3:
			// check successful response
			d := resp.GetText()
			assert.Equal(t, tt.want.data, d)
			assert.Equal(t, tt.want.status, resp.GetStatus())
		}
	}
}

// TestPostText verifies, that:
// 1) null data is not accepted
// 2) already exists is correctly processed
// 2) in case of success - we receive success status
func TestPostText(t *testing.T) {
	// init connect
	ctx, conn := keeperTestConn(t)
	defer conn.Close()

	// create client
	client := pb.NewKeeperClient(conn)
	// init test storage
	storage.InitTest()
	type want struct {
		errStatusCode codes.Code
		errStatusMsg  string
		status        string
	}

	tests := []struct {
		name    string
		number  uint8
		title   string
		body    string
		comment string
		version uint32
		want    want
	}{
		{
			name:   "Test #1: empty data",
			number: 1,
			want: want{
				errStatusCode: codes.InvalidArgument,
				errStatusMsg:  "invalid argument",
			},
		},
		{
			name:    "Test #2: data already exists",
			number:  2,
			title:   testdb.TestText.Title,
			body:    testdb.TestText.Body,
			comment: testdb.TestText.Comment,
			version: testdb.TestText.Version,
			want: want{
				errStatusCode: codes.AlreadyExists,
				errStatusMsg:  "Current / newer version found in database. Please synchronize you app to get the most actual data.",
			},
		},
		{
			name:    "Test #3: correct new data",
			number:  3,
			title:   "new text",
			body:    "new body",
			comment: "new comment",
			version: 1,
			want:    want{status: "success"},
		},
	}
	for _, tt := range tests {
		// make request
		resp, err := client.PostText(ctx,
			&pb.PostTextRequest{Text: &pb.Text{
				Title:   tt.title,
				Body:    tt.body,
				Comment: tt.comment,
				Version: tt.version,
			}})
		switch tt.number {
		case 1:
			// empty values should return err InvalidArgument
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 2:
			// unknown title return err Not Found
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 3:
			// check successful response
			assert.Equal(t, tt.want.status, resp.GetStatus())
		}
	}
}

// TestDelText verifies, that:
// 1) null title is not accepted
// 2) not found received after delete
// 3) in case of success - we receive success message
func TestDelText(t *testing.T) {
	// init connect
	ctx, conn := keeperTestConn(t)
	defer conn.Close()

	// create client
	client := pb.NewKeeperClient(conn)
	// init test storage
	storage.InitTest()
	type want struct {
		errStatusCode codes.Code
		errStatusMsg  string
		status        string
	}

	tests := []struct {
		name   string
		number uint8
		title  string
		want   want
	}{
		{
			name:   "Test #1: empty data",
			number: 1,
			title:  "",
			want: want{
				errStatusCode: codes.InvalidArgument,
				errStatusMsg:  "invalid argument",
			},
		},
		{
			name:   "Test #2: correct data",
			number: 2,
			title:  testdb.TestText.Title,
			want: want{
				status:        "success",
				errStatusCode: codes.NotFound,
				errStatusMsg:  "not found",
			},
		},
	}
	for _, tt := range tests {
		// make request
		resp, err := client.DelText(ctx, &pb.DelTextRequest{Title: tt.title})
		switch tt.number {
		case 1:
			// empty values should return err InvalidArgument
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 2:
			// check successful response of delete
			assert.Equal(t, tt.want.status, resp.GetStatus())
			// make get request to receive not found
			_, err2 := client.GetText(ctx, &pb.GetTextRequest{Title: tt.title})
			// empty values should return err InvalidArgument
			st, ok := status.FromError(err2)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())

		}
	}
}

// TestGetBin verifies, that:
// 1) null title is not accepted
// 2) not found is correctly processed
// 2) in case of success - we receive test text data
func TestGetBin(t *testing.T) {
	// init connect
	ctx, conn := keeperTestConn(t)
	defer conn.Close()

	// create client
	client := pb.NewKeeperClient(conn)
	// init test storage
	storage.InitTest()
	type want struct {
		errStatusCode codes.Code
		errStatusMsg  string
		data          *pb.Bin
		status        string
	}

	tests := []struct {
		name   string
		number uint8
		title  string
		want   want
	}{
		{
			name:   "Test #1: empty data",
			number: 1,
			title:  "",
			want: want{
				errStatusCode: codes.InvalidArgument,
				errStatusMsg:  "invalid argument",
			},
		},
		{
			name:   "Test #2: unknown data",
			number: 2,
			title:  "unknown",
			want: want{
				errStatusCode: codes.NotFound,
				errStatusMsg:  "not found",
			},
		},
		{
			name:   "Test #3: correct data",
			number: 3,
			title:  testdb.TestBin.Title,
			want: want{
				data:   models.ModelsToProtoBin(testdb.TestBin),
				status: "success",
			},
		},
	}
	for _, tt := range tests {
		// make request
		resp, err := client.GetBin(ctx, &pb.GetBinRequest{Title: tt.title})
		switch tt.number {
		case 1:
			// empty values should return err InvalidArgument
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 2:
			// unknown title return err Not Found
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 3:
			// check successful response
			d := resp.GetBinData()
			assert.Equal(t, tt.want.data, d)
			assert.Equal(t, tt.want.status, resp.GetStatus())
		}
	}
}

// TestPostBin verifies, that:
// 1) null data is not accepted
// 2) already exists is correctly processed
// 2) in case of success - we receive success status
func TestPostBin(t *testing.T) {
	// init connect
	ctx, conn := keeperTestConn(t)
	defer conn.Close()

	// create client
	client := pb.NewKeeperClient(conn)
	// init test storage
	storage.InitTest()
	type want struct {
		errStatusCode codes.Code
		errStatusMsg  string
		status        string
	}

	tests := []struct {
		name    string
		number  uint8
		title   string
		body    []byte
		comment string
		version uint32
		want    want
	}{
		{
			name:   "Test #1: empty data",
			number: 1,
			want: want{
				errStatusCode: codes.InvalidArgument,
				errStatusMsg:  "invalid argument",
			},
		},
		{
			name:    "Test #2: data already exists",
			number:  2,
			title:   testdb.TestBin.Title,
			body:    testdb.TestBin.Body,
			comment: testdb.TestBin.Comment,
			version: testdb.TestBin.Version,
			want: want{
				errStatusCode: codes.AlreadyExists,
				errStatusMsg:  "Current / newer version found in database. Please synchronize you app to get the most actual data.",
			},
		},
		{
			name:    "Test #3: correct new data",
			number:  3,
			title:   "new text",
			body:    []byte("new body"),
			comment: "new comment",
			version: 1,
			want:    want{status: "success"},
		},
	}
	for _, tt := range tests {
		// make request
		resp, err := client.PostBin(ctx,
			&pb.PostBinRequest{BinData: &pb.Bin{
				Title:   tt.title,
				Body:    tt.body,
				Comment: tt.comment,
				Version: tt.version,
			}})
		switch tt.number {
		case 1:
			// empty values should return err InvalidArgument
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 2:
			// unknown title return err Not Found
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 3:
			// check successful response
			assert.Equal(t, tt.want.status, resp.GetStatus())
		}
	}
}

// TestDelBin verifies, that:
// 1) null title is not accepted
// 2) not found received after delete
// 3) in case of success - we receive success message
func TestDelBin(t *testing.T) {
	// init connect
	ctx, conn := keeperTestConn(t)
	defer conn.Close()

	// create client
	client := pb.NewKeeperClient(conn)
	// init test storage
	storage.InitTest()
	type want struct {
		errStatusCode codes.Code
		errStatusMsg  string
		status        string
	}

	tests := []struct {
		name   string
		number uint8
		title  string
		want   want
	}{
		{
			name:   "Test #1: empty data",
			number: 1,
			title:  "",
			want: want{
				errStatusCode: codes.InvalidArgument,
				errStatusMsg:  "invalid argument",
			},
		},
		{
			name:   "Test #2: correct data",
			number: 2,
			title:  testdb.TestBin.Title,
			want: want{
				status:        "success",
				errStatusCode: codes.NotFound,
				errStatusMsg:  "not found",
			},
		},
	}
	for _, tt := range tests {
		// make request
		resp, err := client.DelBin(ctx, &pb.DelBinRequest{Title: tt.title})
		switch tt.number {
		case 1:
			// empty values should return err InvalidArgument
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 2:
			// check successful response of delete
			assert.Equal(t, tt.want.status, resp.GetStatus())
			// make get request to receive not found
			_, err2 := client.GetBin(ctx, &pb.GetBinRequest{Title: tt.title})
			// empty values should return err InvalidArgument
			st, ok := status.FromError(err2)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())

		}
	}
}

// TestGetCard verifies, that:
// 1) null title is not accepted
// 2) not found is correctly processed
// 2) in case of success - we receive test text data
func TestGetCard(t *testing.T) {
	// init connect
	ctx, conn := keeperTestConn(t)
	defer conn.Close()

	// create client
	client := pb.NewKeeperClient(conn)
	// init test storage
	storage.InitTest()
	type want struct {
		errStatusCode codes.Code
		errStatusMsg  string
		data          *pb.Card
		status        string
	}

	tests := []struct {
		name   string
		number uint8
		title  string
		want   want
	}{
		{
			name:   "Test #1: empty data",
			number: 1,
			title:  "",
			want: want{
				errStatusCode: codes.InvalidArgument,
				errStatusMsg:  "invalid argument",
			},
		},
		{
			name:   "Test #2: unknown data",
			number: 2,
			title:  "unknown",
			want: want{
				errStatusCode: codes.NotFound,
				errStatusMsg:  "not found",
			},
		},
		{
			name:   "Test #3: correct data",
			number: 3,
			title:  testdb.TestCard.Title,
			want: want{
				data:   models.ModelsToProtoCard(testdb.TestCard),
				status: "success",
			},
		},
	}
	for _, tt := range tests {
		// make request
		resp, err := client.GetCard(ctx, &pb.GetCardRequest{Title: tt.title})
		switch tt.number {
		case 1:
			// empty values should return err InvalidArgument
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 2:
			// unknown title return err Not Found
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 3:
			// check successful response
			d := resp.GetCard()
			assert.Equal(t, tt.want.data, d)
			assert.Equal(t, tt.want.status, resp.GetStatus())
		}
	}
}

// TestPostCard verifies, that:
// 1) null data is not accepted
// 2) already exists is correctly processed
// 2) in case of success - we receive success status
func TestPostCard(t *testing.T) {
	// init connect
	ctx, conn := keeperTestConn(t)
	defer conn.Close()

	// create client
	client := pb.NewKeeperClient(conn)
	// init test storage
	storage.InitTest()
	type want struct {
		errStatusCode codes.Code
		errStatusMsg  string
		status        string
	}

	tests := []struct {
		name       string
		number     uint8
		title      string
		cardNumber string
		expdate    string
		comment    string
		version    uint32
		want       want
	}{
		{
			name:   "Test #1: empty data",
			number: 1,
			want: want{
				errStatusCode: codes.InvalidArgument,
				errStatusMsg:  "invalid argument",
			},
		},
		{
			name:       "Test #2: data already exists",
			number:     2,
			title:      testdb.TestCard.Title,
			cardNumber: testdb.TestCard.Number,
			expdate:    testdb.TestCard.ExpirationDate,
			comment:    testdb.TestCard.Comment,
			version:    testdb.TestCard.Version,
			want: want{
				errStatusCode: codes.AlreadyExists,
				errStatusMsg:  "Current / newer version found in database. Please synchronize you app to get the most actual data.",
			},
		},
		{
			name:       "Test #3: correct new data",
			number:     3,
			title:      "new card",
			cardNumber: "new number",
			expdate:    "new exp date",
			comment:    "new comment",
			version:    1,
			want:       want{status: "success"},
		},
	}
	for _, tt := range tests {
		// make request
		resp, err := client.PostCard(ctx,
			&pb.PostCardRequest{Card: &pb.Card{
				Title:   tt.title,
				Number:  tt.cardNumber,
				Expdate: tt.expdate,
				Comment: tt.comment,
				Version: tt.version,
			}})
		switch tt.number {
		case 1:
			// empty values should returns err InvalidArgument
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 2:
			// data already exists returns AlreadyExists
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 3:
			// check successful response
			assert.Equal(t, tt.want.status, resp.GetStatus())
		}
	}
}

// TestDelCard verifies, that:
// 1) null title is not accepted
// 2) not found received after delete
// 3) in case of success - we receive success message
func TestDelCard(t *testing.T) {
	// init connect
	ctx, conn := keeperTestConn(t)
	defer conn.Close()

	// create client
	client := pb.NewKeeperClient(conn)
	// init test storage
	storage.InitTest()
	type want struct {
		errStatusCode codes.Code
		errStatusMsg  string
		status        string
	}

	tests := []struct {
		name   string
		number uint8
		title  string
		want   want
	}{
		{
			name:   "Test #1: empty data",
			number: 1,
			title:  "",
			want: want{
				errStatusCode: codes.InvalidArgument,
				errStatusMsg:  "invalid argument",
			},
		},
		{
			name:   "Test #2: correct data",
			number: 2,
			title:  testdb.TestCard.Title,
			want: want{
				status:        "success",
				errStatusCode: codes.NotFound,
				errStatusMsg:  "not found",
			},
		},
	}
	for _, tt := range tests {
		// make request
		resp, err := client.DelCard(ctx, &pb.DelCardRequest{Title: tt.title})
		switch tt.number {
		case 1:
			// empty values should return err InvalidArgument
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())
		case 2:
			// check successful response of delete
			assert.Equal(t, tt.want.status, resp.GetStatus())
			// make get request to receive not found
			_, err2 := client.GetCard(ctx, &pb.GetCardRequest{Title: tt.title})
			// empty values should return err InvalidArgument
			st, ok := status.FromError(err2)
			assert.True(t, ok)
			assert.Equal(t, tt.want.errStatusCode, st.Code())
			assert.Equal(t, tt.want.errStatusMsg, st.Message())

		}
	}
}
