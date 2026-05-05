// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package misc

import (
	api "code.aegit.io/aegit/modules/structs"
	"code.aegit.io/aegit/modules/util"
	"code.aegit.io/aegit/modules/web"
	"code.aegit.io/aegit/routers/common"
	"code.aegit.io/aegit/services/context"
)

// Markup render markup document to HTML
func Markup(ctx *context.Context) {
	form := web.GetForm(ctx).(*api.MarkupOption)
	mode := util.Iif(form.Wiki, "wiki", form.Mode) //nolint:staticcheck // form.Wiki is deprecated
	common.RenderMarkup(ctx.Base, ctx.Repo, mode, form.Text, form.Context, form.FilePath)
}
