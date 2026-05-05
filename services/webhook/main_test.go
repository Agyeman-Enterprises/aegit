// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webhook

import (
	"testing"

	"code.aegit.io/aegit/models/unittest"
	"code.aegit.io/aegit/modules/hostmatcher"
	"code.aegit.io/aegit/modules/setting"

	_ "code.aegit.io/aegit/models"
	_ "code.aegit.io/aegit/models/actions"
)

func TestMain(m *testing.M) {
	// for tests, allow only loopback IPs
	setting.Webhook.AllowedHostList = hostmatcher.MatchBuiltinLoopback
	unittest.MainTest(m, &unittest.TestOptions{
		SetUp: func() error {
			setting.LoadQueueSettings()
			return Init()
		},
	})
}
