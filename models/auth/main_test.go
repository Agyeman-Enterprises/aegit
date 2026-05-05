// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package auth_test

import (
	"testing"

	"code.aegit.io/aegit/models/unittest"

	_ "code.aegit.io/aegit/models"
	_ "code.aegit.io/aegit/models/actions"
	_ "code.aegit.io/aegit/models/activities"
	_ "code.aegit.io/aegit/models/auth"
	_ "code.aegit.io/aegit/models/perm/access"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
