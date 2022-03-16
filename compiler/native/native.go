// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"time"

	"github.com/go-vela/server/compiler"

	"github.com/go-vela/server/compiler/registry"
	"github.com/go-vela/server/compiler/registry/github"

	"github.com/go-vela/types"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type ModificationConfig struct {
	Timeout  time.Duration
	Retries  int
	Endpoint string
	Secret   string
}

type client struct {
	Github              registry.Service
	PrivateGithub       registry.Service
	UsePrivateGithub    bool
	ModificationService ModificationConfig

	build    *library.Build
	comment  string
	files    []string
	local    bool
	metadata *types.Metadata
	repo     *library.Repo
	user     *library.User
}

// New returns a Pipeline implementation that integrates with the supported registries.
//
// nolint: revive // ignore returning unexported client
func New(ctx *cli.Context) (*client, error) {
	logrus.Debug("Creating registry clients from CLI configuration")

	c := new(client)

	if ctx.String("modification-addr") != "" {
		c.ModificationService = ModificationConfig{
			Timeout:  ctx.Duration("modification-timeout"),
			Endpoint: ctx.String("modification-addr"),
			Secret:   ctx.String("modification-secret"),
			Retries:  ctx.Int("modification-retries"),
		}
	}

	// setup github template service
	github, err := setupGithub()
	if err != nil {
		return nil, err
	}

	c.Github = github

	if ctx.Bool("github-driver") {
		logrus.Tracef("setting up Private GitHub Client for %s", ctx.String("github-url"))
		// setup private github service
		privGithub, err := setupPrivateGithub(ctx.String("github-url"), ctx.String("github-token"))
		if err != nil {
			return nil, err
		}

		c.PrivateGithub = privGithub
		c.UsePrivateGithub = true
	}

	return c, nil
}

// setupGithub is a helper function to setup the
// Github registry service from the CLI arguments.
func setupGithub() (registry.Service, error) {
	logrus.Tracef("Creating %s registry client from CLI configuration", "github")
	return github.New("", "")
}

// setupPrivateGithub is a helper function to setup the
// Github registry service from the CLI arguments.
func setupPrivateGithub(addr, token string) (registry.Service, error) {
	logrus.Tracef("Creating private %s registry client from CLI configuration", "github")
	return github.New(addr, token)
}

// Duplicate creates a clone of the Engine.
func (c *client) Duplicate() compiler.Engine {
	cc := new(client)

	// copy the essential fields from the existing client
	cc.Github = c.Github
	cc.PrivateGithub = c.PrivateGithub
	cc.UsePrivateGithub = c.UsePrivateGithub
	cc.ModificationService = c.ModificationService

	return cc
}

// WithBuild sets the library build type in the Engine.
func (c *client) WithBuild(b *library.Build) compiler.Engine {
	if b != nil {
		c.build = b
	}

	return c
}

// WithComment sets the comment in the Engine.
func (c *client) WithComment(cmt string) compiler.Engine {
	if cmt != "" {
		c.comment = cmt
	}

	return c
}

// WithFiles sets the changeset files in the Engine.
func (c *client) WithFiles(f []string) compiler.Engine {
	if f != nil {
		c.files = f
	}

	return c
}

// WithLocal sets the compiler metadata type in the Engine.
func (c *client) WithLocal(local bool) compiler.Engine {
	c.local = local

	return c
}

// WithMetadata sets the compiler metadata type in the Engine.
func (c *client) WithMetadata(m *types.Metadata) compiler.Engine {
	if m != nil {
		c.metadata = m
	}

	return c
}

// WithPrivateGitHub sets the private github client in the Engine.
func (c *client) WithPrivateGitHub(url, token string) compiler.Engine {
	if len(url) != 0 && len(token) != 0 {
		privGithub, _ := setupPrivateGithub(url, token)

		c.PrivateGithub = privGithub
	}

	return c
}

// WithRepo sets the library repo type in the Engine.
func (c *client) WithRepo(r *library.Repo) compiler.Engine {
	if r != nil {
		c.repo = r
	}

	return c
}

// WithUser sets the library user type in the Engine.
func (c *client) WithUser(u *library.User) compiler.Engine {
	if u != nil {
		c.user = u
	}

	return c
}
