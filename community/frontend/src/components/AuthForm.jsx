import { useState } from 'react'
import { api, setToken, setUser } from '../api'

export default function AuthForm({ onLogin, addToast }) {
  const [mode, setMode] = useState('login')
  const [form, setForm] = useState({ username: '', password: '', email: '', nickname: '', token: '', new_password: '' })

  function update(k, v) { setForm(f => ({ ...f, [k]: v })) }

  async function submit(e) {
    e.preventDefault()
    console.log('[AUTH] submit called, mode=' + mode)
    try {
      if (mode === 'login') {
        console.log('[AUTH] calling login API for user=' + form.username)
        const data = await api('/users/login', { method: 'POST', body: { username: form.username, password: form.password } })
        console.log('[AUTH] login success, data=', data)
        setToken(data.token)
        setUser(data)
        console.log('[AUTH] redirecting to /feed')
        window.location.href = '/feed'
      } else if (mode === 'register') {
        await api('/users/register', { method: 'POST', body: { username: form.username, password: form.password, email: form.email, nickname: form.nickname } })
        const data = await api('/users/login', { method: 'POST', body: { username: form.username, password: form.password } })
        setToken(data.token)
        setUser(data)
        window.location.href = '/feed'
      } else if (mode === 'forgot') {
        await api('/users/forgot-password', { method: 'POST', body: { email: form.email } })
        addToast('令牌已生成（查看日志），请输入令牌重置密码', 'success')
        setMode('reset')
      } else if (mode === 'reset') {
        await api('/users/reset-password', { method: 'POST', body: { token: form.token, new_password: form.new_password } })
        addToast('密码重置成功，请登录', 'success')
        setMode('login')
      }
    } catch (err) {
      console.error('[AUTH] error:', err.message)
      addToast(err.message, 'error')
    }
  }

  return (
    <div className="auth-page">
      <form className="auth-card" onSubmit={submit} noValidate>
        <h2>{ {login:'登录',register:'注册',forgot:'忘记密码',reset:'重置密码'}[mode] }</h2>
        {(mode === 'register' || mode === 'forgot') && <div className="field"><label>邮箱</label><input type="email" value={form.email} onChange={e=>update('email',e.target.value)} placeholder="可选" /></div>}
        {mode === 'register' && <div className="field"><label>昵称</label><input value={form.nickname} onChange={e=>update('nickname',e.target.value)} placeholder="可选" /></div>}
        <div className="field"><label>用户名</label><input value={form.username} onChange={e=>update('username',e.target.value)} required minLength={3} /></div>
        {mode !== 'reset' && <div className="field"><label>密码</label><input type="password" value={form.password} onChange={e=>update('password',e.target.value)} required minLength={6} /></div>}
        {mode === 'reset' && <div className="field"><label>重置令牌</label><input value={form.token} onChange={e=>update('token',e.target.value)} required /></div>}
        {mode === 'reset' && <div className="field"><label>新密码</label><input type="password" value={form.new_password} onChange={e=>update('new_password',e.target.value)} required minLength={6} /></div>}
        <button className="btn" type="submit">确认</button>
        <div className="text-center text-sm" style={{marginTop:16,color:'var(--text3)'}}>
          {mode === 'login' && <><span onClick={()=>setMode('register')} style={{cursor:'pointer',color:'var(--hover)'}}>注册</span> · <span onClick={()=>setMode('forgot')} style={{cursor:'pointer',color:'var(--text3)'}}>忘记密码</span></>}
          {mode !== 'login' && <span onClick={()=>setMode('login')} style={{cursor:'pointer',color:'var(--hover)'}}>返回登录</span>}
        </div>
      </form>
    </div>
  )
}
