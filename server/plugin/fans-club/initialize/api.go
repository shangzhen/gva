package initialize

import (
	"context"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

// Api 初始化API注册
func Api(ctx context.Context) {
	entities := []model.SysApi{
		{
			Path:        "/fansClub/createFansClub",
			Description: "创建粉丝团",
			ApiGroup:    "粉丝团",
			Method:      "POST",
		},
		{
			Path:        "/fansClub/updateFansClub",
			Description: "更新粉丝团",
			ApiGroup:    "粉丝团",
			Method:      "PUT",
		},
		{
			Path:        "/fansClub/deleteFansClub",
			Description: "删除粉丝团",
			ApiGroup:    "粉丝团",
			Method:      "DELETE",
		},
		{
			Path:        "/fansClub/getFansClub",
			Description: "获取粉丝团详情",
			ApiGroup:    "粉丝团",
			Method:      "GET",
		},
		{
			Path:        "/fansClub/getFansClubList",
			Description: "获取粉丝团列表",
			ApiGroup:    "粉丝团",
			Method:      "GET",
		},
		{
			Path:        "/fansClub/getMyClubs",
			Description: "获取我的粉丝团",
			ApiGroup:    "粉丝团",
			Method:      "GET",
		},
		{
			Path:        "/fansClubMember/joinClub",
			Description: "加入粉丝团",
			ApiGroup:    "粉丝团成员",
			Method:      "POST",
		},
		{
			Path:        "/fansClubMember/quitClub",
			Description: "退出粉丝团",
			ApiGroup:    "粉丝团成员",
			Method:      "POST",
		},
		{
			Path:        "/fansClubMember/getMemberList",
			Description: "获取成员列表",
			ApiGroup:    "粉丝团成员",
			Method:      "GET",
		},
		{
			Path:        "/fansClubMember/updateMemberRole",
			Description: "更新成员角色",
			ApiGroup:    "粉丝团成员",
			Method:      "PUT",
		},
		{
			Path:        "/fansClubMember/removeMember",
			Description: "移除成员",
			ApiGroup:    "粉丝团成员",
			Method:      "DELETE",
		},
		{
			Path:        "/fansClubPost/createPost",
			Description: "创建动态",
			ApiGroup:    "粉丝团动态",
			Method:      "POST",
		},
		{
			Path:        "/fansClubPost/updatePost",
			Description: "更新动态",
			ApiGroup:    "粉丝团动态",
			Method:      "PUT",
		},
		{
			Path:        "/fansClubPost/deletePost",
			Description: "删除动态",
			ApiGroup:    "粉丝团动态",
			Method:      "DELETE",
		},
		{
			Path:        "/fansClubPost/getPostList",
			Description: "获取动态列表",
			ApiGroup:    "粉丝团动态",
			Method:      "GET",
		},
		{
			Path:        "/fansClubPost/getPost",
			Description: "获取动态详情",
			ApiGroup:    "粉丝团动态",
			Method:      "GET",
		},
		{
			Path:        "/fansClubPost/likePost",
			Description: "点赞动态",
			ApiGroup:    "粉丝团动态",
			Method:      "POST",
		},
	}
	utils.RegisterApis(entities...)
}
