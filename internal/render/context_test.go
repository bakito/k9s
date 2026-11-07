// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package render_test

import (
	"testing"

	"github.com/derailed/k9s/internal/model1"
	"github.com/derailed/k9s/internal/render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/tools/clientcmd/api"
)

func TestContextHeader(t *testing.T) {
	var c render.Context

	assert.Len(t, c.Header(""), 4)
}

func TestContextRender(t *testing.T) {
	uu := map[string]struct {
		ctx *render.NamedContext
		e   model1.Row
	}{
		"active": {
			ctx: &render.NamedContext{
				Name: "c1",
				Context: &api.Context{
					LocationOfOrigin: "fred",
					Cluster:          "c1",
					AuthInfo:         "u1",
					Namespace:        "ns1",
				},
				Config: &config{},
			},
			e: model1.Row{
				ID:     "c1",
				Fields: model1.Fields{"c1", "c1", "u1", "ns1"},
			},
		},
	}

	var r render.Context
	for k := range uu {
		uc := uu[k]
		t.Run(k, func(t *testing.T) {
			row := model1.NewRow(4)
			err := r.Render(uc.ctx, "", &row)

			require.NoError(t, err)
			assert.Equal(t, uc.e, row)
		})
	}
}

func TestContextRenderWithDecorateFn(t *testing.T) {
	uu := map[string]struct {
		decorateFn func(string) string
		ctx        *render.NamedContext
		e          model1.Row
	}{
		"no-decorator": {
			decorateFn: nil,
			ctx: &render.NamedContext{
				Name: "c1",
				Context: &api.Context{
					LocationOfOrigin: "fred",
					Cluster:          "c1",
					AuthInfo:         "u1",
					Namespace:        "ns1",
				},
				Config: &config{},
			},
			e: model1.Row{
				ID:     "c1",
				Fields: model1.Fields{"c1", "c1", "u1", "ns1"},
			},
		},
		"with-decorator": {
			decorateFn: func(name string) string {
				return "⚙️ " + name
			},
			ctx: &render.NamedContext{
				Name: "c1",
				Context: &api.Context{
					LocationOfOrigin: "fred",
					Cluster:          "c1",
					AuthInfo:         "u1",
					Namespace:        "ns1",
				},
				Config: &config{},
			},
			e: model1.Row{
				ID:     "c1",
				Fields: model1.Fields{"⚙️ c1", "c1", "u1", "ns1"},
			},
		},
		"decorator-surround": {
			decorateFn: func(name string) string {
				return "[" + name + "]"
			},
			ctx: &render.NamedContext{
				Name: "prod",
				Context: &api.Context{
					LocationOfOrigin: "fred",
					Cluster:          "c1",
					AuthInfo:         "u1",
					Namespace:        "ns1",
				},
				Config: &config{},
			},
			e: model1.Row{
				ID:     "prod",
				Fields: model1.Fields{"[prod]", "c1", "u1", "ns1"},
			},
		},
	}

	for n := range uu {
		uc := uu[n]
		t.Run(n, func(t *testing.T) {
			var r render.Context
			if uc.decorateFn != nil {
				r.SetDecorateFn(uc.decorateFn)
			}
			row := model1.NewRow(4)
			err := r.Render(uc.ctx, "", &row)

			require.NoError(t, err)
			assert.Equal(t, uc.e, row)
		})
	}
}

// ----------------------------------------------------------------------------
// Helpers...

type config struct{}

func (config) CurrentContextName() (string, error) {
	return "fred", nil
}
