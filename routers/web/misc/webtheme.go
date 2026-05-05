// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package misc

import (
	"net/http"

	"code.aegit.io/aegit/modules/optional"
	"code.aegit.io/aegit/modules/templates"
	"code.aegit.io/aegit/modules/util"
	"code.aegit.io/aegit/modules/web/middleware"
	"code.aegit.io/aegit/services/context"
	user_service "code.aegit.io/aegit/services/user"
	"code.aegit.io/aegit/services/webtheme"
)

func WebThemeList(ctx *context.Context) {
	curWebTheme := ctx.TemplateContext.CurrentWebTheme()
	renderUtils := templates.NewRenderUtils(ctx)
	allThemes := webtheme.GetAvailableThemes()

	var results []map[string]any
	for _, theme := range allThemes {
		results = append(results, map[string]any{
			"name":  renderUtils.RenderThemeItem(theme, 14),
			"value": theme.InternalName,
			"class": "item js-aria-clickable" + util.Iif(theme.InternalName == curWebTheme.InternalName, " selected", ""),
		})
	}
	ctx.JSON(http.StatusOK, map[string]any{"results": results})
}

func WebThemeApply(ctx *context.Context) {
	themeName := ctx.FormString("theme")
	if ctx.Doer != nil {
		opts := &user_service.UpdateOptions{Theme: optional.Some(themeName)}
		_ = user_service.UpdateUser(ctx, ctx.Doer, opts)
	} else {
		middleware.SetSiteCookie(ctx.Resp, middleware.CookieTheme, themeName, 0)
	}
}
