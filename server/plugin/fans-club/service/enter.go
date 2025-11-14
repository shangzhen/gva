package service

type ServiceGroup struct {
	FansClubService
	FansClubMemberService
	FansClubPostService
}

var ServiceGroupApp = new(ServiceGroup)
