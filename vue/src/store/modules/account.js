import {
  login
} from '@/api/account'
import {
  getToken,
  setToken,
  removeToken
} from '@/utils/auth'
import {
  resetRouter
} from '@/router'

const getDefaultState = () => {
  return {
    token: getToken(),
    accountId: '',
    userName: '',
    balance: 0,
    roles: []
  }
}

const state = getDefaultState()

const mutations = {
  RESET_STATE: (state) => {
    Object.assign(state, getDefaultState())
  },
  SET_TOKEN: (state, token) => {
    state.token = token
  },
  SET_ACCOUNTID: (state, accountId) => {
    state.accountId = accountId
  },
  SET_USERNAME: (state, userName) => {
    state.userName = userName
  },
  SET_BALANCE: (state, balance) => {
    state.balance = balance
  },
  SET_ROLES: (state, roles) => {
    state.roles = roles
  }
}

const actions = {
  login({
    commit
  }, accountId) {
    return new Promise((resolve, reject) => {
      login({
        args: [{
          accountId: accountId
        }]
      }).then(response => {
        commit('SET_TOKEN', response[0].accountId)
        setToken(response[0].accountId)
        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },
  // get user info
  getInfo({
    commit,
    state
  }) {
    return new Promise((resolve, reject) => {
      login({
        args: [{
          accountId: state.token
        }]
      }).then(response => {
        var roles
        if (response[0].userName === '管理员') {
          roles = ['admin']
        } else {
          roles = ['editor']
        }
        commit('SET_ROLES', roles)
        commit('SET_ACCOUNTID', response[0].accountId)
        commit('SET_USERNAME', response[0].userName)
        commit('SET_BALANCE', response[0].balance)
        resolve(roles)
      }).catch(error => {
        reject(error)
      })
    })
  },
  logout({
    commit
  }) {
    return new Promise(resolve => {
      removeToken()
      resetRouter()
      commit('RESET_STATE')
      resolve()
    })
  },

  resetToken({
    commit
  }) {
    return new Promise(resolve => {
      removeToken()
      commit('RESET_STATE')
      resolve()
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
