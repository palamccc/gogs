// Copyright 2016 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package org

import (
	api "github.com/gogits/go-gogs-client"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/modules/context"
	"github.com/gogits/gogs/routers/api/v1/convert"
	"github.com/gogits/gogs/routers/api/v1/user"
)

func AddOrgMembership(ctx *context.APIContext, form api.AddOrgMembershipOption) {
	org := ctx.Org.Organization
	if !org.IsOwnedBy(ctx.User.ID) {
		ctx.Status(403)
		return
	}

	user := user.GetUserByParams(ctx)
	if ctx.Written() {
		return
	}

	// TODO: put add memeber and add team in one session in case of roll back.
	if err := org.AddMember(user.ID); err != nil {
		ctx.Error(500, "org.AddMember", err)
		return
	}

	team, err := org.GetOwnerTeam()
	if err != nil {
		ctx.Error(500, "GetOwnerTeam", err)
		return
	}
	if form.Role == "admin" {
		if err := team.AddMember(user.ID); err != nil {
			ctx.Error(500, "team.AddMember", err)
			return
		}
	} else {
		if err := team.RemoveMember(user.ID); err != nil {
			if models.IsErrLastOrgOwner(err) {
				ctx.Error(422, "", err)
			} else {
				ctx.Error(500, "team.RemoveMember", err)
			}
			return
		}
	}

	ctx.JSON(200, map[string]interface{}{
		"organization": convert.ToOrganization(org),
		"user":         convert.ToUser(user),
	})
}
