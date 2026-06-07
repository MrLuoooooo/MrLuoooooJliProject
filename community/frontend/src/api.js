const API = 'http://localhost:1807/api/v1'
const WS_URL = 'ws://localhost:1807/ws'

let token = localStorage.getItem('token') || ''
let user = null

function headers() {
  const h = {}
  if (token) h['Authorization'] = `Bearer ${token}`
  return h
}

export async function api(path, opts = {}) {
  const h = { ...headers(), ...(opts.headers || {}) }
  if (opts.body && typeof opts.body === 'object') {
    h['Content-Type'] = 'application/json'
    opts.body = JSON.stringify(opts.body)
  }
  const res = await fetch(API + path, { ...opts, headers: h })
  const data = await res.json()
  if (data.code !== 0) throw new Error(data.msg || '请求失败')
  return data.data
}

export function setToken(t) { token = t; localStorage.setItem('token', t || '') }
export function getToken() { return token }
export function setUser(u) { user = u }
export function getUser() { return user }
export function clearAuth() { token = ''; user = null; localStorage.removeItem('token') }
export function getWsUrl() { return token ? `${WS_URL}?token=${token}` : null }
