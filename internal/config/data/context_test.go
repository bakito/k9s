// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package data_test

import (
	"testing"

	"github.com/derailed/k9s/internal/config/data"
	"github.com/derailed/k9s/internal/config/mock"
	"github.com/stretchr/testify/assert"
)

func TestClusterValidate(t *testing.T) {
	c := data.NewContext()
	c.Validate(mock.NewMockConnection(), "ct-1", "cl-1")

	assert.Equal(t, data.DefaultView, c.View.Active)
	assert.Equal(t, "default", c.Namespace.Active)
	assert.Len(t, c.Namespace.Favorites, 1)
	assert.Equal(t, []string{"default"}, c.Namespace.Favorites)
}

func TestClusterValidateEmpty(t *testing.T) {
	c := data.NewContext()
	c.Validate(mock.NewMockConnection(), "ct-1", "cl-1")

	assert.Equal(t, data.DefaultView, c.View.Active)
	assert.Equal(t, "default", c.Namespace.Active)
	assert.Len(t, c.Namespace.Favorites, 1)
	assert.Equal(t, []string{"default"}, c.Namespace.Favorites)
}

func TestContextDecorateContextName(t *testing.T) {
	uu := map[string]struct {
		prefix string
		suffix string
		name   string
		expect string
	}{
		"no-decoration": {
			prefix: "",
			suffix: "",
			name:   "my-context",
			expect: "my-context",
		},
		"prefix-only": {
			prefix: "⚙️ ",
			suffix: "",
			name:   "my-context",
			expect: "⚙️ my-context",
		},
		"suffix-only": {
			prefix: "",
			suffix: " ⚡",
			name:   "prod",
			expect: "prod ⚡",
		},
		"prefix-and-suffix": {
			prefix: "[",
			suffix: "]",
			name:   "dev",
			expect: "[dev]",
		},
	}

	for n := range uu {
		u := uu[n]
		t.Run(n, func(t *testing.T) {
			c := data.NewContext()
			c.ContextPrefix = u.prefix
			c.ContextSuffix = u.suffix
			assert.Equal(t, u.expect, c.DecorateContextName(u.name))
		})
	}
}
