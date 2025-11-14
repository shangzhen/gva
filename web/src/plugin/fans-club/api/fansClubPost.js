import service from '@/utils/request'

// 创建动态
export const createPost = (data) => {
  return service({
    url: '/fansClubPost/createPost',
    method: 'post',
    data
  })
}

// 更新动态
export const updatePost = (data) => {
  return service({
    url: '/fansClubPost/updatePost',
    method: 'put',
    data
  })
}

// 删除动态
export const deletePost = (data) => {
  return service({
    url: '/fansClubPost/deletePost',
    method: 'delete',
    data
  })
}

// 获取动态列表
export const getPostList = (params) => {
  return service({
    url: '/fansClubPost/getPostList',
    method: 'get',
    params
  })
}

// 获取动态详情
export const getPost = (params) => {
  return service({
    url: '/fansClubPost/getPost',
    method: 'get',
    params
  })
}

// 点赞动态
export const likePost = (data) => {
  return service({
    url: '/fansClubPost/likePost',
    method: 'post',
    data
  })
}
