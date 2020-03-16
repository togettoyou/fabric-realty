const getters = {
  sidebar: state => state.app.sidebar,
  device: state => state.app.device,
  token: state => state.account.token,
  accountId: state => state.account.accountId,
  userName: state => state.account.userName,
  balance: state => state.account.balance,
  roles: state => state.account.roles,
  permission_routes: state => state.permission.routes
}
export default getters
