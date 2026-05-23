package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"amlakcrm/backend/internal/domain"
	"amlakcrm/backend/internal/searchrpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type RequestMatcher interface {
	RequestMatches(ctx context.Context, userID, businessID, contactID, requestID string, limit, offset int) ([]domain.PropertyMatchResult, int, error)
}

type SearchMatchingClient struct {
	addr   string
	token  string
	mu     sync.Mutex
	conn   *grpc.ClientConn
	client searchrpc.SearchServiceClient
}

func NewSearchMatchingClient(addr, token string) *SearchMatchingClient {
	return &SearchMatchingClient{
		addr:  strings.TrimSpace(addr),
		token: token,
	}
}

func (s *SearchMatchingClient) RequestMatches(ctx context.Context, userID, businessID, contactID, requestID string, limit, offset int) ([]domain.PropertyMatchResult, int, error) {
	if s == nil || s.addr == "" {
		return nil, 0, errors.New("search service is not configured")
	}
	client, err := s.grpcClient()
	if err != nil {
		return nil, 0, err
	}
	callCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	if s.token != "" {
		callCtx = metadata.AppendToOutgoingContext(callCtx, "authorization", "Bearer "+s.token)
	}
	resp, err := client.RequestMatches(callCtx, &searchrpc.RequestMatchesRequest{
		UserID:     userID,
		BusinessID: businessID,
		ContactID:  contactID,
		RequestID:  requestID,
		Limit:      int32(limit),
		Offset:     int32(offset),
	}, searchrpc.ForceJSONCodec())
	if err != nil {
		return nil, 0, fmt.Errorf("search service unavailable: %w", err)
	}
	return resp.Matches, int(resp.Total), nil
}

func (s *SearchMatchingClient) grpcClient() (searchrpc.SearchServiceClient, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.client != nil {
		return s.client, nil
	}
	conn, err := grpc.NewClient(s.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	s.conn = conn
	s.client = searchrpc.NewSearchServiceClient(conn)
	return s.client, nil
}
