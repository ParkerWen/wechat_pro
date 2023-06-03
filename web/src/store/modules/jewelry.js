import {
  create,
  getInfo,
  upscale,
  variation,
  getImageList,
} from "@/api/jewelry";

const actions = {
  create({ commit }, jewelryInfo) {
    const { prompt } = jewelryInfo;
    return new Promise((resolve, reject) => {
      create({ prompt: prompt.trim() })
        .then((response) => {
          const { data } = response;
          resolve(data);
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  upscale({ commit }, upscaleReq) {
    const { task_id, index } = upscaleReq;
    return new Promise((resolve, reject) => {
      upscale({ task_id: task_id.trim(), index: index })
        .then((response) => {
          const { data } = response;
          resolve(data);
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  variation({ commit }, variationReq) {
    const { task_id, index } = variationReq;
    return new Promise((resolve, reject) => {
      variation({ task_id: task_id.trim(), index: index })
        .then((response) => {
          const { data } = response;
          resolve(data);
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  getInfo({ commit }, id) {
    return new Promise((resolve, reject) => {
      getInfo({ id: id })
        .then((response) => {
          const { data } = response;
          if (!data) {
            return reject("Invalid Task.");
          }
          resolve(data);
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
  getImageList({ commit }) {
    return new Promise((resolve, reject) => {
      getImageList()
        .then((response) => {
          const { data } = response;
          if (!data) {
            return reject("Invalid Image List.");
          }
          resolve(data);
        })
        .catch((error) => {
          reject(error);
        });
    });
  },
};

export default {
  namespaced: true,
  actions,
};
