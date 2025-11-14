import service from '@/utils/request'

// 加入粉丝团
export const joinClub = (data) => {
  return service({
    url: '/fansClubMember/joinClub',
    method: 'post',
    data
  })
}

// 退出粉丝团
export const quitClub = (data) => {
  return service({
    url: '/fansClubMember/quitClub',
    method: 'post',
    data
  })
}

// 获取成员列表
export const getMemberList = (params) => {
  return service({
    url: '/fansClubMember/getMemberList',
    method: 'get',
    params
  })
}

// 更新成员角色
export const updateMemberRole = (data) => {
  return service({
    url: '/fansClubMember/updateMemberRole',
    method: 'put',
    data
  })
}

// 移除成员
export const removeMember = (data) => {
  return service({
    url: '/fansClubMember/removeMember',
    method: 'delete',
    data
  })
}
