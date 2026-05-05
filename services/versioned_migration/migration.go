// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package versioned_migration

import (
	"context"

	"code.aegit.io/aegit/models/migrations"
	"code.aegit.io/aegit/modules/globallock"

	"xorm.io/xorm"
)

func Migrate(ctx context.Context, x *xorm.Engine) error {
	// only one instance can do the migration at the same time if there are multiple instances
	release, err := globallock.Lock(ctx, "gitea_versioned_migration")
	if err != nil {
		return err
	}
	defer release()

	return migrations.Migrate(ctx, x)
}
