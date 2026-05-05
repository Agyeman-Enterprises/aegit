// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package system_test

import (
	"testing"

	"code.aegit.io/aegit/models/unittest"

	_ "code.aegit.io/aegit/models" // register models
	_ "code.aegit.io/aegit/models/actions"
	_ "code.aegit.io/aegit/models/activities"
	_ "code.aegit.io/aegit/models/system" // register models of system
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
