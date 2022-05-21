package git

import (
	"errors"
	"net/http"

	"github.com/go-git/go-git/v5"
	ghttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

type client struct {
	username    string
	accessToken string
	repo        *git.Repository
}

func (c *client) Clone(path, url string) error {
	o := &git.CloneOptions{
		URL: url,
		Auth: &ghttp.BasicAuth{
			Username: c.username,
			Password: c.accessToken,
		},
	}
	repo, err := git.PlainClone(path, false, o)
	if err != nil {
		return err
	}

	c.repo = repo

	return nil
}

func (c *client) authTest() error {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	req.SetBasicAuth(c.username, c.accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("token is invalid")
	}

	return nil
}

func NewClient(username, accessToken string) (*client, error) {
	c := &client{
		username:    username,
		accessToken: accessToken,
	}

	err := c.authTest()
	if err != nil {
		return nil, err
	}

	return c, nil
}
