package main

import (
	"context"
	"dagger/dagger-github-comment/internal/dagger"
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v63/github"
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

	// Event payload.
	payload *dagger.File,
) error {
	payloadContents, err := payload.Contents(ctx)
	if err != nil {
		return err
	}

	var event github.IssueCommentEvent

	if err := json.Unmarshal([]byte(payloadContents), &event); err != nil {
		return err
	}

	fmt.Println("Event: ", event.Issue.GetID())

	return nil
}
