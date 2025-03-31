// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/compiler/registry"
	"github.com/go-vela/server/compiler/registry/github"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal"
	"github.com/go-vela/server/internal/image"
	"github.com/go-vela/server/scm"
)

type ModificationConfig struct {
	Timeout  time.Duration
	Retries  int
	Endpoint string
	Secret   string
}

type client struct {
	Github           registry.Service
	PrivateGithub    registry.Service
	UsePrivateGithub bool

	ModificationService ModificationConfig
	TemplateCache       map[string][]byte

	settings.Compiler

	build          *api.Build
	comment        string
	commit         string
	files          []string
	local          bool
	localTemplates []string
	metadata       *internal.Metadata
	repo           *api.Repo
	user           *api.User
	labels         []string
	db             database.Interface
	scm            scm.Service
	netrc          *string
}

// FromCLIContext returns a Pipeline implementation that integrates with the supported registries.
//

func FromCLIContext(ctx *cli.Context) (*client, error) {
	logrus.Debug("creating registry clients from CLI configuration")

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
	github, err := setupGithub(ctx.Context)
	if err != nil {
		return nil, err
	}

	c.Github = github

	c.Compiler = settings.Compiler{}

	cloneImage := ctx.String("clone-image")

	// validate clone image
	_, err = image.ParseWithError(cloneImage)
	if err != nil {
		return nil, fmt.Errorf("invalid clone image %s: %w", cloneImage, err)
	}

	// set the clone image to use for the injected clone step
	c.SetCloneImage(cloneImage)

	// set the template depth to use for nested templates
	c.SetTemplateDepth(ctx.Int("max-template-depth"))

	// set the starlark execution step limit for compiling starlark pipelines
	c.SetStarlarkExecLimit(ctx.Int64("compiler-starlark-exec-limit"))

	if ctx.Bool("github-driver") {
		logrus.Tracef("setting up Private GitHub Client for %s", ctx.String("github-url"))
		// setup private github service
		privGithub, err := setupPrivateGithub(ctx.Context, ctx.String("github-url"), ctx.String("github-token"))
		if err != nil {
			return nil, err
		}

		c.PrivateGithub = privGithub
		c.UsePrivateGithub = true
	}

	c.TemplateCache = make(map[string][]byte)

	return c, nil
}

// setupGithub is a helper function to setup the
// Github registry service from the CLI arguments.
func setupGithub(ctx context.Context) (registry.Service, error) {
	logrus.Tracef("creating %s registry client from CLI configuration", "github")
	return github.New(ctx, "", "")
}

// setupPrivateGithub is a helper function to setup the
// Github registry service from the CLI arguments.
func setupPrivateGithub(ctx context.Context, addr, token string) (registry.Service, error) {
	logrus.Tracef("creating private %s registry client from CLI configuration", "github")
	return github.New(ctx, addr, token)
}

// Duplicate creates a clone of the Engine.
func (c *client) Duplicate() compiler.Engine {
	cc := new(client)

	// copy the essential fields from the existing client
	cc.Github = c.Github
	cc.PrivateGithub = c.PrivateGithub
	cc.UsePrivateGithub = c.UsePrivateGithub
	cc.ModificationService = c.ModificationService
	cc.CloneImage = c.CloneImage
	cc.TemplateDepth = c.TemplateDepth
	cc.StarlarkExecLimit = c.StarlarkExecLimit
	cc.TemplateCache = make(map[string][]byte)

	return cc
}

// WithBuild sets the API build type in the Engine.
func (c *client) WithBuild(b *api.Build) compiler.Engine {
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

// WithCommit sets the comment in the Engine.
func (c *client) WithCommit(cmt string) compiler.Engine {
	if cmt != "" {
		c.commit = cmt
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

// WithLocalTemplates sets the compiler local templates in the Engine.
func (c *client) WithLocalTemplates(templates []string) compiler.Engine {
	c.localTemplates = templates

	return c
}

// WithMetadata sets the compiler metadata type in the Engine.
func (c *client) WithMetadata(m *internal.Metadata) compiler.Engine {
	if m != nil {
		c.metadata = m
	}

	return c
}

// WithPrivateGitHub sets the private github client in the Engine.
func (c *client) WithPrivateGitHub(ctx context.Context, url, token string) compiler.Engine {
	if len(url) != 0 && len(token) != 0 {
		privGithub, _ := setupPrivateGithub(ctx, url, token)

		c.PrivateGithub = privGithub
	}

	return c
}

// WithRepo sets the API repo type in the Engine.
func (c *client) WithRepo(r *api.Repo) compiler.Engine {
	if r != nil {
		c.repo = r
	}

	return c
}

// WithUser sets the API user type in the Engine.
func (c *client) WithUser(u *api.User) compiler.Engine {
	if u != nil {
		c.user = u
	}

	return c
}

// WithLabels sets the label(s) in the Engine.
func (c *client) WithLabels(labels []string) compiler.Engine {
	if len(labels) != 0 {
		c.labels = labels
	}

	return c
}

// WithNetrc sets the netrc in the Engine.
func (c *client) WithNetrc(n string) compiler.Engine {
	c.netrc = &n

	return c
}

// WithSCM sets the scm in the Engine.
func (c *client) WithSCM(_scm scm.Service) compiler.Engine {
	c.scm = _scm

	return c
}

// WithDatabase sets the database in the Engine.
func (c *client) WithDatabase(db database.Interface) compiler.Engine {
	c.db = db

	return c
}
