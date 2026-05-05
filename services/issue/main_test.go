// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package issue

import (
	"testing"

	"code.aegit.io/aegit/models/unittest"

	_ "code.aegit.io/aegit/models"
	_ "code.aegit.io/aegit/models/actions"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
