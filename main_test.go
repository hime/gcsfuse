// Copyright 2024 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPosixFlagsConversion(t *testing.T) {
	data := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Single hyphen converted to double",
			input:    []string{"abc", "-flag"},
			expected: []string{"abc", "--flag"},
		},
		{
			name:     "Double-hyphens remain unchanged",
			input:    []string{"abc", "--flag"},
			expected: []string{"abc", "--flag"},
		},
		{
			name:     "Single hyphen converted to double for flags with values",
			input:    []string{"abc", "-flag=\"test\""},
			expected: []string{"abc", "--flag=\"test\""},
		},
		{
			name:     "Values with hyphen stay unchanged",
			input:    []string{"abc", "--flag=\"-test\""},
			expected: []string{"abc", "--flag=\"-test\""},
		},
		{
			name:     "Help shorthand stays unchanged.",
			input:    []string{"abc", "-h"},
			expected: []string{"abc", "-h"},
		},
		{
			name:     "Version shorthand stays unchanged.",
			input:    []string{"abc", "-v"},
			expected: []string{"abc", "-v"},
		},
		{
			name:     "Help flag becomes help shorthand.",
			input:    []string{"abc", "--h"},
			expected: []string{"abc", "-h"},
		},
		{
			name:     "Version flag becomes version shorthand.",
			input:    []string{"abc", "--v"},
			expected: []string{"abc", "-v"},
		},
	}
	for _, tc := range data {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, convertToPosixArgs(tc.input))
		})
	}
}
