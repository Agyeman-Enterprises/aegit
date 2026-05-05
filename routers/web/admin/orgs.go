// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2020 The Gitea Authors.
// SPDX-License-Identifier: MIT

package admin

import (
	"code.aegit.io/aegit/models/db"
	user_model "code.aegit.io/aegit/models/user"
	"code.aegit.io/aegit/modules/setting"
	"code.aegit.io/aegit/modules/structs"
	"code.aegit.io/aegit/modules/templates"
	"code.aegit.io/aegit/routers/web/explore"
	"code.aegit.io/aegit/services/context"
)

const (
	tplOrgs templates.TplName = "admin/org/list"
)

// Organizations show all the organizations
func Organizations(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("admin.organizations")
	ctx.Data["PageIsAdminOrganizations"] = true

	if ctx.FormString("sort") == "" {
		ctx.SetFormString("sort", UserSearchDefaultAdminSort)
	}

	explore.RenderUserSearch(ctx, user_model.SearchUserOptions{
		Actor:           ctx.Doer,
		Types:           []user_model.UserType{user_model.UserTypeOrganization},
		IncludeReserved: true, // administrator needs to list all accounts include reserved
		ListOptions: db.ListOptions{
			PageSize: setting.UI.Admin.OrgPagingNum,
		},
		Visible: []structs.VisibleType{structs.VisibleTypePublic, structs.VisibleTypeLimited, structs.VisibleTypePrivate},
	}, tplOrgs)
}
