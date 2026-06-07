import { NavLink, useNavigate } from 'react-router-dom'
import { useState, useEffect, useRef } from 'react'
import { api, getWsUrl } from '../api'

export default function Layout({ user, onLogout, children, addToast }) {
  const navigate = useNavigate()
  const wsRef = useRef(null)
  const [badges, setBadges] = useState({ msg: 0, notif: 0 })

  useEffect(() => {
    api('/notifications/unread').then(d => setBadges(b => ({ ...b, notif: d?.count || 0 }))).catch(()=>{})
    api('/messages/unread').then(d => setBadges(b => ({ ...b, msg: d?.count || 0 }))).catch(()=>{})
  }, [])

  useEffect(() => {
    const url = getWsUrl()
    if (!url) return
    function connect() {
      const ws = new WebSocket(url)
      wsRef.current = ws
      ws.onmessage = e => {
        const { type, data } = JSON.parse(e.data)
        if (type === 'notification') { setBadges(b => ({ ...b, notif: b.notif + 1 })); addToast(data.content?.slice(0, 20) + '...', 'info') }
        if (type === 'message') { setBadges(b => ({ ...b, msg: b.msg + 1 })); addToast('新消息', 'info') }
        if (type === 'broadcast') addToast('📢 ' + data.content, 'info')
      }
      ws.onclose = () => setTimeout(connect, 3000)
    }
    connect()
    return () => wsRef.current?.close()
  }, [])

  return (
    <div className="app">
      <div className="sidebar">
        <div className="logo">Community</div>
        <nav>
          <NavLink to="/feed">🏠 首页</NavLink>
          <NavLink to="/create">✏️ 发布</NavLink>
          <NavLink to="/messages">💬 私信 {badges.msg > 0 && <span className="badge">{badges.msg}</span>}</NavLink>
          <NavLink to="/notifications">🔔 通知 {badges.notif > 0 && <span className="badge">{badges.notif}</span>}</NavLink>
          <NavLink to="/search">🔍 搜索</NavLink>
          <NavLink to="/ai">🤖 AI</NavLink>
          <NavLink to="/tags">🏷️ 标签</NavLink>
          {user?.admin_type > 0 && <NavLink to="/admin">⚙️ 管理</NavLink>}
        </nav>
        <div className="user" onClick={() => navigate('/profile')}>
          <div className="avatar">{(user?.nickname || user?.username || 'U')[0].toUpperCase()}</div>
          <div>
            <div className="name">{user?.nickname || user?.username}</div>
            <div className="role" style={{fontSize:11,color:'var(--text3)'}}>#{user?.id}</div>
          </div>
        </div>
      </div>
      <div className="main">
        <div className="topbar">
          <h1>Community</h1>
          <button className="btn-outline" onClick={onLogout}>退出</button>
        </div>
        <div className="content">{children}</div>
      </div>
    </div>
  )
}
