package shipper

import (
	"context"
	"fmt"
	"time"

	"github.com/elastic/elastic-agent-libs/mapstr"
	"github.com/elastic/elastic-agent-shipper-client/pkg/helpers"
	sc "github.com/elastic/elastic-agent-shipper-client/pkg/proto"
	"github.com/elastic/elastic-agent-shipper-client/pkg/proto/messages"
	"github.com/elastic/go-ucfg"
	"github.com/leehinman/spigot/pkg/output"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const Name = "shipper"

type ShipperOutput struct {
	client sc.ProducerClient
	conn   *grpc.ClientConn
	batch  []*messages.Event
	config config
}

func init() {
	output.Register(Name, New)
}

// New is factory for creating a new Shipper output
func New(cfg *ucfg.Config) (s output.Output, err error) {
	c := defaultConfig()
	if err := cfg.Unpack(&c); err != nil {
		return nil, err
	}
	opts := defaultDialOptions()
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, c.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s: %w", c.Address, err)
	}
	client := sc.NewProducerClient(conn)

	so := &ShipperOutput{
		client: client,
		conn:   conn,
		config: c,
	}
	return so, nil
}

func (s *ShipperOutput) Write(b []byte) (n int, err error) {
	source := &messages.Source{
		InputId:  s.config.InputId,
		StreamId: s.config.StreamId,
	}
	datastream := &messages.DataStream{
		Type:      s.config.DataStreamType,
		Dataset:   s.config.DataStreamDataset,
		Namespace: s.config.DataStreamNamespace,
	}
	meta := mapstr.M{
		"input_id":  s.config.InputId,
		"stream_id": s.config.StreamId,
	}

	metaStruct, err := helpers.NewStruct(meta)
	if err != nil {
		return 0, err
	}
	fields := mapstr.M{
		"message": string(b),
		"data_stream": mapstr.M{
			"type":      s.config.DataStreamType,
			"data_set":  s.config.DataStreamDataset,
			"namespace": s.config.DataStreamNamespace,
		},
	}
	fieldsStruct, err := helpers.NewStruct(fields)
	if err != nil {
		return 0, err
	}
	e := &messages.Event{
		Timestamp:  timestamppb.New(time.Now()),
		Source:     source,
		DataStream: datastream,
		Metadata:   metaStruct,
		Fields:     fieldsStruct,
	}
	s.batch = append(s.batch, e)
	return len(b), nil
}

func (s *ShipperOutput) Close() error {
	if len(s.batch) > 0 {
		if err := s.send(); err != nil {
			return err
		}
	}
	return s.conn.Close()
}

func (s *ShipperOutput) NewInterval() error {
	return s.send()
}

func (s *ShipperOutput) send() error {
	if s.conn == nil {
		return fmt.Errorf("connection is not established to: %s", s.config.Address)
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout)
	defer cancel()
	_, err := s.client.PublishEvents(ctx, &messages.PublishRequest{Events: s.batch})
	if err != nil {
		return fmt.Errorf("publish events failed: %w", err)
	}
	s.batch = s.batch[0:0]
	return nil
}
