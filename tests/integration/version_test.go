// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"net/http"
	"testing"

	"code.aegit.io/aegit/modules/setting"
	"code.aegit.io/aegit/modules/structs"
	"code.aegit.io/aegit/tests"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	setting.AppVer = "test-version-1"
	req := NewRequest(t, "GET", "/api/v1/version")
	resp := MakeRequest(t, req, http.StatusOK)

	version := DecodeJSON(t, resp, &structs.ServerVersion{})
	assert.Equal(t, setting.AppVer, version.Version)
}
