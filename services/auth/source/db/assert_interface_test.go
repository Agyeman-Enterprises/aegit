// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package db_test

import (
	auth_model "code.aegit.io/aegit/models/auth"
	"code.aegit.io/aegit/services/auth"
	"code.aegit.io/aegit/services/auth/source/db"
)

// This test file exists to assert that our Source exposes the interfaces that we expect
// It tightly binds the interfaces and implementation without breaking go import cycles

type sourceInterface interface {
	auth.PasswordAuthenticator
	auth_model.Config
}

var _ (sourceInterface) = &db.Source{}
