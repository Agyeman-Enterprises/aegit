// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package pipeline

import (
	"testing"

	"code.aegit.io/aegit/modules/git"
)

func TestMain(m *testing.M) {
	git.RunGitTests(m)
}
