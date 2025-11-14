import service from '@/utils/request'

// 创建粉丝团
export const createFansClub = (data) => {
  return service({
    url: '/fansClub/createFansClub',
    method: 'post',
    data
  })
}

// 更新粉丝团
export const updateFansClub = (data) => {
  return service({
    url: '/fansClub/updateFansClub',
    method: 'put',
    data
  })
}

// 删除粉丝团
export const deleteFansClub = (data) => {
  return service({
    url: '/fansClub/deleteFansClub',
    method: 'delete',
    data
  })
}

// 获取粉丝团详情
export const getFansClub = (params) => {
  return service({
    url: '/fansClub/getFansClub',
    method: 'get',
    params
  })
}

// 获取粉丝团列表
export const getFansClubList = (params) => {
  return service({
    url: '/fansClub/getFansClubList',
    method: 'get',
    params
  })
}

// 获取我的粉丝团
export const getMyClubs = (params) => {
  return service({
    url: '/fansClub/getMyClubs',
    method: 'get',
    params
  })
}
