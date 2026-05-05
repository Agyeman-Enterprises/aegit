// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo_test

import (
	"testing"

	"code.aegit.io/aegit/models/unittest"

	_ "code.aegit.io/aegit/models" // register table model
	_ "code.aegit.io/aegit/models/actions"
	_ "code.aegit.io/aegit/models/activities"
	_ "code.aegit.io/aegit/models/perm/access" // register table model
	_ "code.aegit.io/aegit/models/repo"        // register table model
	_ "code.aegit.io/aegit/models/user"        // register table model
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
