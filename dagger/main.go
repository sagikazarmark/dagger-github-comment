package main

import (
	"context"
	"dagger/dagger-github-comment/internal/dagger"
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v63/github"
	"golang.org/x/oauth2"
)

type DaggerGithubComment struct {
	// +private
	GithubToken *dagger.Secret
}

func New(githubToken *dagger.Secret) *DaggerGithubComment {
	return &DaggerGithubComment{
		GithubToken: githubToken,
	}
}

func (m *DaggerGithubComment) Process(
	ctx context.Context,

	eventName string,

	// Event payload.
	payload *dagger.File,
) error {
	payloadContents, err := payload.Contents(ctx)
	if err != nil {
		return err
	}

	githubToken, err := m.GithubToken.Plaintext(ctx)
	if err != nil {
		return err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	handler := issueCommentEventHandler{
		Client: client,
	}

	switch eventName {
	case "issue_comment":
		var event github.IssueCommentEvent

		if err := json.Unmarshal([]byte(payloadContents), &event); err != nil {
			return err
		}

		handler.handle(ctx, event)

	default:
		fmt.Println("unknown event: ", eventName)
	}

	return nil
}

type issueCommentEventHandler struct {
	Client *github.Client
}

func (h issueCommentEventHandler) handle(ctx context.Context, event github.IssueCommentEvent) error {
	_, _, err := h.Client.Reactions.CreateIssueCommentReaction(
		ctx,
		event.GetRepo().GetOwner().GetLogin(),
		event.GetRepo().GetName(),
		event.GetComment().GetID(),
		"+1",
	)
	if err != nil {
		return err
	}

	return nil
}
