// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package access_test

import (
	"testing"

	"code.aegit.io/aegit/models/unittest"

	_ "code.aegit.io/aegit/models"
	_ "code.aegit.io/aegit/models/actions"
	_ "code.aegit.io/aegit/models/activities"
	_ "code.aegit.io/aegit/models/repo"
	_ "code.aegit.io/aegit/models/user"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
