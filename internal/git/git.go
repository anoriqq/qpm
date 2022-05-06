package git

import (
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type client struct {
	username    string
	accessToken string
	repo        *git.Repository
}

func (c *client) Open(path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	c.repo = repo

	return nil
}

func (c *client) Clone(ctx context.Context, path, url string) error {
	o := &git.CloneOptions{
		URL: url,
		Auth: &http.BasicAuth{
			Username: c.username,
			Password: c.accessToken,
		},
	}
	repo, err := git.PlainCloneContext(ctx, path, false, o)
	if err != nil {
		return err
	}

	c.repo = repo

	return nil
}

func NewClient(username, accessToken string) *client {
	return &client{
		username:    username,
		accessToken: accessToken,
	}
}
