// Copyright 2022 ASoulDocs. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark/ast"
)

func TestIDs_Generate(t *testing.T) {
	ids := newIDs()

	tests := []struct {
		name  string
		value string
		kind  ast.NodeKind
		want  string
	}{
		{
			name:  "normal",
			value: "Hello, 世界",
			kind:  ast.KindHeading,
			want:  "Hello, 世界",
		},
		{
			name:  "empty heading",
			value: "",
			kind:  ast.KindHeading,
			want:  "heading",
		},
		{
			name:  "empty id",
			value: "",
			kind:  ast.KindImage,
			want:  "id",
		},

		{
			name:  "duplicated heading",
			value: "",
			kind:  ast.KindHeading,
			want:  "heading-1",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ids.Generate([]byte(test.value), test.kind)
			assert.Equal(t, test.want, string(got))
		})
	}
}
