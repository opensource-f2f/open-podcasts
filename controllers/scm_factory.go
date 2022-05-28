package controllers

import (
	"context"
	"fmt"
	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
	"strings"
)

// GetGitProvider returns the GitReleaser implement by kind
func GetGitProvider(kind, server, token string) (client *github.Client, err error) {
	ctx := context.Background()

	switch strings.ToLower(kind) {
	case "github":
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)

		client = github.NewClient(tc)
	default:
		err = fmt.Errorf("unknown scm provider: %s", kind)
	}
	return
}
