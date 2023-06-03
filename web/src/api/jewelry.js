import request from "@/utils/request";

export function create(data) {
  return request({
    url: "/mj/task/v1/imagine",
    method: "post",
    data,
  });
}

export function upscale(data) {
  return request({
    url: "/mj/task/v1/upscale",
    method: "post",
    data,
  });
}

export function variation(data) {
  return request({
    url: "/mj/task/v1/variation",
    method: "post",
    data,
  });
}

export function getInfo(data) {
  return request({
    url: "/mj/task/v1/fetch",
    method: "get",
    params: data,
  });
}

export function getImageList(data) {
  return request({
    url: "/mj/task/v1/list",
    method: "get",
    params: data,
  });
}
