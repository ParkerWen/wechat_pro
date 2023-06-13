import storage from 'store'
import expirePlugin from 'store/plugins/expire'
import { login, getInfo } from '@/api/login'
import { register } from '@/api/register'
import { ACCESS_TOKEN } from '@/store/mutation-types'
import { welcome } from '@/utils/util'

storage.addPlugin(expirePlugin)
const user = {
  state: {
    token: '',
    name: '',
    welcome: '',
    avatar: '',
    roles: [],
    info: {}
  },

  mutations: {
    SET_TOKEN: (state, token) => {
      state.token = token
    },
    SET_NAME: (state, { name, welcome }) => {
      state.name = name
      state.welcome = welcome
    },
    SET_AVATAR: (state, avatar) => {
      state.avatar = avatar
    },
    SET_ROLES: (state, roles) => {
      state.roles = roles
    },
    SET_INFO: (state, info) => {
      state.info = info
    }
  },

  actions: {
    // 注册
    Register ({ commit }, userInfo) {
      return new Promise((resolve, reject) => {
        register(userInfo).then(response => {
          resolve()
        }).catch(error => {
          reject(error)
        })
      })
    },
    // 登录
    Login ({ commit }, userInfo) {
      return new Promise((resolve, reject) => {
        login(userInfo).then(response => {
          const data = response.data
          storage.set(ACCESS_TOKEN, data.access_token, new Date().getTime() + 7 * 24 * 60 * 60 * 1000)
          commit('SET_TOKEN', data.access_token)
          resolve()
        }).catch(error => {
          reject(error)
        })
      })
    },

    // 获取用户信息
    GetInfo ({ commit }) {
      return new Promise((resolve, reject) => {
        // 请求后端获取用户信息 /api/user/info
        getInfo().then(response => {
          const result = response.data
          console.log(result)
          result.roleId = 'admin'
          result.role = roleObj
          if (result.role && result.role.permissions.length > 0) {
            const role = { ...result.role }
            role.permissions = result.role.permissions.map(permission => {
              const per = {
                ...permission,
                actionList: (permission.actionEntitySet || {}).map(item => item.action)
               }
              return per
            })
            role.permissionList = role.permissions.map(permission => { return permission.permissionId })
            // 覆盖响应体的 role, 供下游使用
            result.role = role

            commit('SET_ROLES', role)
            commit('SET_INFO', result)
            commit('SET_NAME', { name: result.name, welcome: welcome() })
            commit('SET_AVATAR', result.avatar)
            // 下游
            resolve(result)
            console.log(result)
          } else {
            reject(new Error('getInfo: roles must be a non-null array !'))
          }
        }).catch(error => {
          reject(error)
        })
      })
    },

    // 登出
    Logout ({ commit, state }) {
      return new Promise((resolve) => {
        commit('SET_TOKEN', '')
        commit('SET_ROLES', [])
        storage.remove(ACCESS_TOKEN)
        resolve()
      })
    }

  }
}

// role
const roleObj = {
  id: 'admin',
  name: '管理员',
  describe: '拥有所有权限',
  status: 1,
  creatorId: 'system',
  createTime: 1497160610259,
  deleted: 0,
  permissions: [
    {
      roleId: 'admin',
      permissionId: 'jewelry',
      permissionName: '珠宝设计',
      actions:
        '[{"action":"add","defaultCheck":false,"describe":"新增"},{"action":"query","defaultCheck":false,"describe":"查询"},{"action":"get","defaultCheck":false,"describe":"详情"},{"action":"update","defaultCheck":false,"describe":"修改"},{"action":"delete","defaultCheck":false,"describe":"删除"}]',
      actionEntitySet: [
        {
          action: 'add',
          describe: '新增',
          defaultCheck: false
        },
        {
          action: 'query',
          describe: '查询',
          defaultCheck: false
        },
        {
          action: 'get',
          describe: '详情',
          defaultCheck: false
        },
        {
          action: 'update',
          describe: '修改',
          defaultCheck: false
        },
        {
          action: 'delete',
          describe: '删除',
          defaultCheck: false
        }
      ],
      actionList: null,
      dataAccess: null
    }
  ]
}

export default user
