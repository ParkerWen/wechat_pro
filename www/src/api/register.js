import request from '@/utils/request'

const userApi = {
    Register: '/mj/user/v1/email/register'
}

export function register (parameter) {
    return request({
      url: userApi.Register,
      method: 'post',
      data: parameter
    })
}
