// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package markdown

import (
	"os"
	"testing"

	"code.aegit.io/aegit/modules/markup"
	"code.aegit.io/aegit/modules/setting"
)

func TestMain(m *testing.M) {
	setting.IsInTesting = true
	markup.RenderBehaviorForTesting.DisableAdditionalAttributes = true
	os.Exit(m.Run())
}
