package grpcsrv

import (
	context "context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/app"
	memorystorage "github.com/nsmak/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/runtime/protoimpl"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

type mockLogger struct {
}

func (m *mockLogger) Info(msg string, fields ...zap.Field) {
}

func (m *mockLogger) Warn(msg string, fields ...zap.Field) {
}

func (m *mockLogger) Error(msg string, fields ...zap.Field) {
}

func (m *mockLogger) String(key string, val string) zap.Field {
	return zap.Field{}
}

func (m *mockLogger) Int64(key string, val int64) zap.Field {
	return zap.Field{}
}

func (m *mockLogger) Duration(key string, val time.Duration) zap.Field {
	return zap.Field{}
}

func TestCreateEventSuccess(t *testing.T) {
	s := grpcServer()
	defer s.Stop()

	c, conn := grpcClient()
	defer conn.Close()

	ctx := context.Background()

	resp, err := c.CreateEvent(ctx, &Event{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Id:            "unique_event_id_3",
		Title:         "Event_Title_3",
		StartDate:     0,
		EndDate:       0,
		Description:   "Event_Description_3",
		OwnerId:       "",
		RemindIn:      0,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestCreateEventFail(t *testing.T) {
	s := grpcServer()
	defer s.Stop()

	c, conn := grpcClient()
	defer conn.Close()

	ctx := context.Background()

	resp, err := c.CreateEvent(ctx, &Event{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Id:            "unique_event_id_1",
		Title:         "Event_Title_1",
		StartDate:     100500,
		EndDate:       300800,
		Description:   "Event_Description_1",
		OwnerId:       "unique_owner_uid",
		RemindIn:      0,
	})

	require.Error(t, err)
	require.Nil(t, resp)
}

func TestUpdateEventSuccess(t *testing.T) {
	s := grpcServer()
	defer s.Stop()

	c, conn := grpcClient()
	defer conn.Close()

	ctx := context.Background()

	resp, err := c.UpdateEvent(ctx, &Event{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Id:            "unique_event_id_1",
		Title:         "Event_Title_New",
		StartDate:     100500,
		EndDate:       300800,
		Description:   "Event_Description_New",
		OwnerId:       "unique_owner_uid",
		RemindIn:      0,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestUpdateEventFail(t *testing.T) {
	s := grpcServer()
	defer s.Stop()

	c, conn := grpcClient()
	defer conn.Close()

	ctx := context.Background()

	resp, err := c.UpdateEvent(ctx, &Event{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Id:            "unique_event_id_3",
		Title:         "Event_Title_3",
		StartDate:     0,
		EndDate:       0,
		Description:   "Event_Description_3",
		OwnerId:       "",
		RemindIn:      0,
	})

	require.Error(t, err)
	require.Nil(t, resp)
}

func TestRemoveEventSuccess(t *testing.T) {
	s := grpcServer()
	defer s.Stop()

	c, conn := grpcClient()
	defer conn.Close()

	ctx := context.Background()

	resp, err := c.RemoveEvent(ctx, &EventID{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Id:            "unique_event_id_1",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestRemoveEventFail(t *testing.T) {
	s := grpcServer()
	defer s.Stop()

	c, conn := grpcClient()
	defer conn.Close()

	ctx := context.Background()

	resp, err := c.RemoveEvent(ctx, &EventID{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Id:            "NaN",
	})

	require.Error(t, err)
	require.Nil(t, resp)
}

func TestEventsSuccess(t *testing.T) {
	s := grpcServer()
	defer s.Stop()

	c, conn := grpcClient()
	defer conn.Close()

	ctx := context.Background()

	resp, err := c.Events(ctx, &EventsQuery{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		From:          300800,
		To:            500200,
	})

	require.NoError(t, err)
	require.Len(t, resp.Events, 1)
	require.Equal(t, resp.Events[0].Id, "unique_event_id_2")
}

func TestEventsFail(t *testing.T) {
	s := grpcServer()
	defer s.Stop()

	c, conn := grpcClient()
	defer conn.Close()

	ctx := context.Background()

	resp, err := c.Events(ctx, &EventsQuery{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		From:          900800,
		To:            10000200,
	})

	require.Error(t, err)
	require.Nil(t, resp)
}

func grpcServer() *grpc.Server {
	store := memorystorage.New()
	for _, e := range mockEvents() {
		_ = store.NewEvent(context.Background(), e)
	}
	api := NewAPI(app.New(&mockLogger{}, store))
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	RegisterEventServiceServer(s, api)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("server stop with error: %v", err)
		}
	}()
	return s
}

func grpcClient() (EventServiceClient, *grpc.ClientConn) {
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can't dial: %v", err)
	}

	client := NewEventServiceClient(conn)

	return client, conn
}

func mockEvents() []app.Event {
	return []app.Event{
		{
			ID:          "unique_event_id_1",
			Title:       "Event_Title_1",
			StartDate:   100500,
			EndDate:     300800,
			Description: "Event_Description_1",
			OwnerID:     "unique_owner_uid",
			RemindIn:    0,
		},
		{
			ID:          "unique_event_id_2",
			Title:       "Event_Title_2",
			StartDate:   300800,
			EndDate:     300900,
			Description: "Event_Description_2",
			OwnerID:     "unique_owner_uid",
			RemindIn:    0,
		},
	}
}
