// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package routers

import (
	"fmt"

	"github.com/Unknwon/paginater"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/modules/base"
	"github.com/gogits/gogs/modules/middleware"
	"github.com/gogits/gogs/modules/setting"
	"github.com/gogits/gogs/routers/user"
)

const (
	HOME          base.TplName = "home"
	EXPLORE_REPOS base.TplName = "explore/repos"
)

func Home(ctx *middleware.Context) {
	if ctx.IsSigned {
		if !ctx.User.IsActive && setting.Service.RegisterEmailConfirm {
			ctx.Data["Title"] = ctx.Tr("auth.active_your_account")
			ctx.HTML(200, user.ACTIVATE)
		} else {
			user.Dashboard(ctx)
		}
		return
	}

	// Check auto-login.
	uname := ctx.GetCookie(setting.CookieUserName)
	if len(uname) != 0 {
		ctx.Redirect(setting.AppSubUrl + "/user/login")
		return
	}

	ctx.Data["PageIsHome"] = true
	ctx.HTML(200, HOME)
}

func Explore(ctx *middleware.Context) {
	ctx.Data["Title"] = ctx.Tr("explore")
	ctx.Data["PageIsExplore"] = true
	ctx.Data["PageIsExploreRepositories"] = true

	page := ctx.QueryInt("page")
	if page <= 1 {
		page = 1
	}

	ctx.Data["Page"] = paginater.New(int(models.CountPublicRepositories()), setting.ExplorePagingNum, page, 5)

	repos, err := models.GetRecentUpdatedRepositories(page)
	if err != nil {
		ctx.Handle(500, "GetRecentUpdatedRepositories", err)
		return
	}
	for _, repo := range repos {
		if err = repo.GetOwner(); err != nil {
			ctx.Handle(500, "GetOwner", fmt.Errorf("%d: %v", repo.ID, err))
			return
		}
	}
	ctx.Data["Repos"] = repos

	ctx.HTML(200, EXPLORE_REPOS)
}

func NotFound(ctx *middleware.Context) {
	ctx.Data["Title"] = "Page Not Found"
	ctx.Handle(404, "home.NotFound", nil)
}
