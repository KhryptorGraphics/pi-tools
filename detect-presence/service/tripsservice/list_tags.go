package tripsservice

import (
	"context"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
)

func (s *Server) ListTags(ctx context.Context, req *tripspb.ListTagsRequest) (*tripspb.ListTagsResponse, error) {
	span := trace.SpanFromContext(ctx)

	var limit int32 = 100
	if req.GetLimit() > 0 && req.GetLimit() < 100 {
		limit = req.GetLimit()
	}

	span.SetAttributes(label.Int32("limit", limit))

	tags, err := s.q.ListTags(ctx, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "listing tags: %s", err)
	}

	span.SetAttributes(label.Int("tag.count", len(tags)))

	res := &tripspb.ListTagsResponse{}

	for _, tag := range tags {
		res.Tags = append(res.Tags, &tripspb.Tag{
			Name:      tag.Name,
			TripCount: tag.TripCount,
		})
	}

	return res, nil
}
