package searchrpc

import (
	"context"
	"encoding/json"

	"amlakcrm/backend/internal/domain"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

const ServiceName = "amlak.search.v1.SearchService"

func init() {
	encoding.RegisterCodec(jsonCodec{})
}

type RequestMatchesRequest struct {
	UserID     string `json:"userId"`
	BusinessID string `json:"businessId"`
	ContactID  string `json:"contactId"`
	RequestID  string `json:"requestId"`
	Limit      int32  `json:"limit,omitempty"`
	Offset     int32  `json:"offset,omitempty"`
}

type RequestMatchesResponse struct {
	Matches []domain.PropertyMatchResult `json:"matches"`
	Total   int32                        `json:"total"`
	Limit   int32                        `json:"limit"`
	Offset  int32                        `json:"offset"`
}

type SearchServiceServer interface {
	RequestMatches(context.Context, *RequestMatchesRequest) (*RequestMatchesResponse, error)
}

type SearchServiceClient interface {
	RequestMatches(ctx context.Context, in *RequestMatchesRequest, opts ...grpc.CallOption) (*RequestMatchesResponse, error)
}

type searchServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSearchServiceClient(cc grpc.ClientConnInterface) SearchServiceClient {
	return &searchServiceClient{cc: cc}
}

func (c *searchServiceClient) RequestMatches(ctx context.Context, in *RequestMatchesRequest, opts ...grpc.CallOption) (*RequestMatchesResponse, error) {
	out := new(RequestMatchesResponse)
	err := c.cc.Invoke(ctx, "/"+ServiceName+"/RequestMatches", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func RegisterSearchServiceServer(server grpc.ServiceRegistrar, service SearchServiceServer) {
	server.RegisterService(&grpc.ServiceDesc{
		ServiceName: ServiceName,
		HandlerType: (*SearchServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "RequestMatches",
				Handler:    requestMatchesHandler,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "search.proto",
	}, service)
}

func requestMatchesHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestMatchesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).RequestMatches(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/" + ServiceName + "/RequestMatches",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).RequestMatches(ctx, req.(*RequestMatchesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

type jsonCodec struct{}

func (jsonCodec) Name() string {
	return "json"
}

func (jsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (jsonCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func ForceJSONCodec() grpc.CallOption {
	return grpc.ForceCodec(jsonCodec{})
}
