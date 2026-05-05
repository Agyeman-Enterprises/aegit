// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package models

import (
	"testing"

	activities_model "code.aegit.io/aegit/models/activities"
	"code.aegit.io/aegit/models/organization"
	repo_model "code.aegit.io/aegit/models/repo"
	"code.aegit.io/aegit/models/unittest"
	user_model "code.aegit.io/aegit/models/user"

	_ "code.aegit.io/aegit/models/actions"
	_ "code.aegit.io/aegit/models/system"

	"github.com/stretchr/testify/assert"
)

// TestFixturesAreConsistent assert that test fixtures are consistent
func TestFixturesAreConsistent(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	unittest.CheckConsistencyFor(t,
		&user_model.User{},
		&repo_model.Repository{},
		&organization.Team{},
		&activities_model.Action{})
}

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
