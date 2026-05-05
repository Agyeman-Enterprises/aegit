// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package organization_test

import (
	"testing"

	"code.aegit.io/aegit/models/organization"
	"code.aegit.io/aegit/models/perm"
	"code.aegit.io/aegit/models/repo"
	"code.aegit.io/aegit/models/unit"
	"code.aegit.io/aegit/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestGetTeamsWithAccessToRepoUnit(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	org41 := unittest.AssertExistsAndLoadBean(t, &organization.Organization{ID: 41})
	repo61 := unittest.AssertExistsAndLoadBean(t, &repo.Repository{ID: 61})

	teams, err := organization.GetTeamsWithAccessToAnyRepoUnit(t.Context(), org41.ID, repo61.ID, perm.AccessModeRead, unit.TypePullRequests)
	assert.NoError(t, err)
	if assert.Len(t, teams, 2) {
		assert.EqualValues(t, 21, teams[0].ID)
		assert.EqualValues(t, 22, teams[1].ID)
	}
}
