// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package label

import (
	"code.aegit.io/aegit/modules/label"
	"code.aegit.io/aegit/modules/web"
	"code.aegit.io/aegit/services/context"
	"code.aegit.io/aegit/services/forms"
)

func GetLabelEditForm(ctx *context.Context) *forms.CreateLabelForm {
	form := web.GetForm(ctx).(*forms.CreateLabelForm)
	if ctx.HasError() {
		ctx.JSONError(ctx.GetErrMsg())
		return nil
	}
	var err error
	form.Color, err = label.NormalizeColor(form.Color)
	if err != nil {
		ctx.JSONError(ctx.Tr("repo.issues.label_color_invalid"))
		return nil
	}
	return form
}
