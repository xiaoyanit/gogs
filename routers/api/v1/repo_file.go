// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package v1

import (
	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/modules/git"
	"github.com/gogits/gogs/modules/middleware"
	"github.com/gogits/gogs/routers/repo"
)

func GetRepoRawFile(ctx *middleware.Context) {
	if !ctx.Repo.HasAccess() {
		ctx.Error(404)
		return
	}

	blob, err := ctx.Repo.Commit.GetBlobByPath(ctx.Repo.TreeName)
	if err != nil {
		if err == git.ErrNotExist {
			ctx.Error(404)
		} else {
			ctx.APIError(500, "GetBlobByPath", err)
		}
		return
	}
	if err = repo.ServeBlob(ctx, blob); err != nil {
		ctx.APIError(500, "ServeBlob", err)
	}
}

func GetRepoArchive(ctx *middleware.Context) {
	repoPath := models.RepoPath(ctx.Params(":username"), ctx.Params(":reponame"))
	gitRepo, err := git.OpenRepository(repoPath)
	if err != nil {
		ctx.APIError(500, "OpenRepository", err)
		return
	}
	ctx.Repo.GitRepo = gitRepo

	repo.Download(ctx)
}
