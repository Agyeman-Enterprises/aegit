// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package explore

import (
	"net/http"

	"code.aegit.io/aegit/models/db"
	repo_model "code.aegit.io/aegit/models/repo"
	api "code.aegit.io/aegit/modules/structs"
	"code.aegit.io/aegit/services/context"
	"code.aegit.io/aegit/services/convert"
)

// TopicSearch search for creating topic
func TopicSearch(ctx *context.Context) {
	opts := &repo_model.FindTopicOptions{
		Keyword: ctx.FormString("q"),
		ListOptions: db.ListOptions{
			Page:     ctx.FormInt("page"),
			PageSize: convert.ToCorrectPageSize(ctx.FormInt("limit")),
		},
	}

	topics, total, err := db.FindAndCount[repo_model.Topic](ctx, opts)
	if err != nil {
		ctx.HTTPError(http.StatusInternalServerError)
		return
	}

	topicResponses := make([]*api.TopicResponse, len(topics))
	for i, topic := range topics {
		topicResponses[i] = convert.ToTopicResponse(topic)
	}

	ctx.SetTotalCountHeader(total)
	ctx.JSON(http.StatusOK, map[string]any{
		"topics": topicResponses,
	})
}
