package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club/api"
	"github.com/gin-gonic/gin"
)

type FansClubRouter struct{}

func (r *FansClubRouter) InitFansClubRouter(group *gin.RouterGroup) {
	fansClubApi := api.ApiGroupApp.FansClubApi
	fansClubMemberApi := api.ApiGroupApp.FansClubMemberApi
	fansClubPostApi := api.ApiGroupApp.FansClubPostApi

	// 使用中间件
	fansClubRouter := group.Group("fansClub").Use(middleware.OperationRecord())
	fansClubMemberRouter := group.Group("fansClubMember").Use(middleware.OperationRecord())
	fansClubPostRouter := group.Group("fansClubPost").Use(middleware.OperationRecord())

	// 粉丝团路由
	{
		fansClubRouter.POST("createFansClub", fansClubApi.CreateFansClub)
		fansClubRouter.PUT("updateFansClub", fansClubApi.UpdateFansClub)
		fansClubRouter.DELETE("deleteFansClub", fansClubApi.DeleteFansClub)
		fansClubRouter.GET("getFansClub", fansClubApi.GetFansClub)
		fansClubRouter.GET("getFansClubList", fansClubApi.GetFansClubList)
		fansClubRouter.GET("getMyClubs", fansClubApi.GetMyClubs)
	}

	// 成员路由
	{
		fansClubMemberRouter.POST("joinClub", fansClubMemberApi.JoinClub)
		fansClubMemberRouter.POST("quitClub", fansClubMemberApi.QuitClub)
		fansClubMemberRouter.GET("getMemberList", fansClubMemberApi.GetMemberList)
		fansClubMemberRouter.PUT("updateMemberRole", fansClubMemberApi.UpdateMemberRole)
		fansClubMemberRouter.DELETE("removeMember", fansClubMemberApi.RemoveMember)
	}

	// 动态路由
	{
		fansClubPostRouter.POST("createPost", fansClubPostApi.CreatePost)
		fansClubPostRouter.PUT("updatePost", fansClubPostApi.UpdatePost)
		fansClubPostRouter.DELETE("deletePost", fansClubPostApi.DeletePost)
		fansClubPostRouter.GET("getPostList", fansClubPostApi.GetPostList)
		fansClubPostRouter.GET("getPost", fansClubPostApi.GetPost)
		fansClubPostRouter.POST("likePost", fansClubPostApi.LikePost)
	}
}
