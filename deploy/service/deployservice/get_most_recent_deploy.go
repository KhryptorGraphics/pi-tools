package deployservice

import (
	"context"
	"strings"
	"time"

	"github.com/google/go-github/v33/github"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
)

func (s *Server) GetMostRecentDeploy(ctx context.Context, _ *deploypb.GetMostRecentDeployRequest) (*deploypb.GetMostRecentDeployResponse, error) {
	span := trace.SpanFromContext(ctx)

	repoParts := strings.SplitN(s.Config.GitHubRepo, "/", 2)
	deployments, _, err := s.GitHubClient.Repositories.ListDeployments(ctx, repoParts[0], repoParts[1], &github.DeploymentsListOptions{
		Task:        "deploy",
		Environment: "production",
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	})
	if err != nil {
		return nil, err
	}

	span.SetAttributes(label.Int("deployment.count", len(deployments)))
	if len(deployments) == 0 {
		return nil, status.Errorf(codes.NotFound, "no deployments found for GitHub repository %s", s.Config.GitHubRepo)
	}

	deploy := deployments[0]
	span.SetAttributes(
		label.Int64("deployment.id", deploy.GetID()),
		label.String("deployment.sha", deploy.GetSHA()))

	deployResponse := &deploypb.Deploy{
		Id:        deploy.GetID(),
		CommitSha: deploy.GetSHA(),
		StartedAt: deploy.GetCreatedAt().Format(time.RFC3339),
	}

	commit, _, err := s.GitHubClient.Repositories.GetCommit(ctx, repoParts[0], repoParts[1], deploy.GetSHA())
	if err != nil {
		return nil, err
	}

	deployResponse.CommitMessage = commit.GetCommit().GetMessage()

	statuses, _, err := s.GitHubClient.Repositories.ListDeploymentStatuses(ctx, repoParts[0], repoParts[1], deploy.GetID(), &github.ListOptions{
		PerPage: 1,
	})
	if err != nil {
		return nil, err
	}

	span.SetAttributes(label.Int("deployment.status_count", len(statuses)))
	if len(statuses) > 0 {
		deployStatus := statuses[0]
		span.SetAttributes(
			label.Int64("deployment.status_id", deployStatus.GetID()),
			label.String("deployment.state", deployStatus.GetState()))

		switch deployStatus.GetState() {
		case "in_progress":
			deployResponse.State = deploypb.Deploy_IN_PROGRESS
		case "success":
			deployResponse.State = deploypb.Deploy_SUCCESS
			deployResponse.FinishedAt = deployStatus.GetCreatedAt().Format(time.RFC3339)
		case "failure":
			deployResponse.State = deploypb.Deploy_FAILURE
			deployResponse.FinishedAt = deployStatus.GetCreatedAt().Format(time.RFC3339)
		case "inactive":
			deployResponse.State = deploypb.Deploy_INACTIVE
		default:
			deployResponse.State = deploypb.Deploy_UNKNOWN
		}
	} else {
		deployResponse.State = deploypb.Deploy_PENDING
	}

	return &deploypb.GetMostRecentDeployResponse{
		Deploy: deployResponse,
	}, nil
}
