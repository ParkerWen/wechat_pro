import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/mj/user/v1/email/login',
    method: 'post',
    data
  })
}

export function signup(data) {
  return request({
    url: '/mj/user/v1/email/register',
    method: 'post',
    data
  })
}

export function getInfo(token) {
  return request({
    url: '/mj/user/v1/fetch',
    method: 'get',
    params: { token }
  })
}

export function logout() {
  return request({
    url: '/vue-admin-template/user/logout',
    method: 'post'
  })
}
