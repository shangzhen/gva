package api

type ApiGroup struct {
	FansClubApi
	FansClubMemberApi
	FansClubPostApi
}

var ApiGroupApp = new(ApiGroup)
