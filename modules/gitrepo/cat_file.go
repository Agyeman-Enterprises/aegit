// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package gitrepo

import (
	"context"

	"code.aegit.io/aegit/modules/git"
)

func NewBatch(ctx context.Context, repo Repository) (git.CatFileBatchCloser, error) {
	return git.NewBatch(ctx, repoPath(repo))
}
