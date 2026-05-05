// Copyright 2018 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package utils

import (
	"testing"

	"code.aegit.io/aegit/models/unittest"
	"code.aegit.io/aegit/modules/setting"
	webhook_service "code.aegit.io/aegit/services/webhook"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m, &unittest.TestOptions{
		SetUp: func() error {
			setting.LoadQueueSettings()
			return webhook_service.Init()
		},
	})
}
