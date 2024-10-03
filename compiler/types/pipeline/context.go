// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"
)

// contextKey defines the key type for
// storing pipeline types in a context.
type contextKey int

const (
	// buildKey defines the key type for
	// storing a Build type in a context.
	buildKey contextKey = iota

	// secretKey defines the key type for
	// storing a Secret type in a context.
	secretKey

	// stageKey defines the key type for
	// storing a Stage type in a context.
	stageKey

	// containerKey defines the key type for
	// storing a Step type in a context.
	containerKey
)

// BuildFromContext retrieves the Build type from the context.
func BuildFromContext(c context.Context) *Build {
	// get build value from context
	v := c.Value(buildKey)
	if v == nil {
		return nil
	}

	// cast build value to expected Build type
	b, ok := v.(*Build)
	if !ok {
		return nil
	}

	return b
}

// BuildWithContext inserts the Build type to the context.
func BuildWithContext(c context.Context, b *Build) context.Context {
	return context.WithValue(c, buildKey, b)
}

// SecretFromContext retrieves the Secret type from the context.
func SecretFromContext(c context.Context) *Secret {
	// get secret value from context
	v := c.Value(secretKey)
	if v == nil {
		return nil
	}

	// cast secret value to expected Secret type
	s, ok := v.(*Secret)
	if !ok {
		return nil
	}

	return s
}

// SecretWithContext inserts the Secret type to the context.
func SecretWithContext(c context.Context, s *Secret) context.Context {
	return context.WithValue(c, secretKey, s)
}

// StageFromContext retrieves the Stage type from the context.
func StageFromContext(c context.Context) *Stage {
	// get stage value from context
	v := c.Value(stageKey)
	if v == nil {
		return nil
	}

	// cast stage value to expected Stage type
	s, ok := v.(*Stage)
	if !ok {
		return nil
	}

	return s
}

// StageWithContext inserts the Stage type to the context.
func StageWithContext(c context.Context, s *Stage) context.Context {
	return context.WithValue(c, stageKey, s)
}

// ContainerFromContext retrieves the container type from the context.
func ContainerFromContext(c context.Context) *Container {
	// get container value from context
	v := c.Value(containerKey)
	if v == nil {
		return nil
	}

	// cast step value to expected Container type
	s, ok := v.(*Container)
	if !ok {
		return nil
	}

	return s
}

// ContainerWithContext inserts the Container type to the context.
func ContainerWithContext(c context.Context, s *Container) context.Context {
	return context.WithValue(c, containerKey, s)
}
