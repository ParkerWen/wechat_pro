import request from '@/utils/request'

const jewelryApi = {
  List: '/mj/task/v1/list',
  Imagine: '/mj/task/v1/imagine',
  Upscale: '/mj/task/v1/upscale',
  Variation: '/mj/task/v1/variation'
}

export function list () {
  return request({
    url: jewelryApi.List,
    method: 'get',
    headers: {
      'Content-Type': 'application/json;charset=UTF-8'
    }
  })
}

export function imagine (parameter) {
  return request({
    url: jewelryApi.Imagine,
    method: 'post',
    data: parameter
  })
}

export function upscale (parameter) {
  return request({
    url: jewelryApi.Upscale,
    method: 'post',
    data: parameter
  })
}

export function variation (parameter) {
  return request({
    url: jewelryApi.Variation,
    method: 'post',
    data: parameter
  })
}
