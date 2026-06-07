import { useState, useEffect } from 'react'
import { api } from '../api'

export default function Admin({ addToast }) {
  const [tab, setTab] = useState('dashboard')
  const [users, setUsers] = useState(null)
  const [posts, setPosts] = useState(null)
  const [stats, setStats] = useState(null)

  useEffect(() => { loadTab() }, [tab])

  async function loadTab() {
    try {
      if (tab === 'dashboard') {
        const d = await api('/admin/stats'); setStats(d)
      } else if (tab === 'users') {
        const d = await api('/admin/users?page=1&page_size=50'); setUsers(d?.items||[])
      } else {
        const d = await api('/admin/posts?page=1&page_size=50'); setPosts(d?.items||[])
      }
    } catch (err) { addToast(err.message, 'error') }
  }

  return <>
    <div className="tabs">
      <button className={`tab ${tab==='dashboard'?'active':''}`} onClick={()=>setTab('dashboard')}>仪表盘</button>
      <button className={`tab ${tab==='users'?'active':''}`} onClick={()=>setTab('users')}>用户管理</button>
      <button className={`tab ${tab==='posts'?'active':''}`} onClick={()=>setTab('posts')}>帖子管理</button>
      <button className={`tab ${tab==='broadcast'?'active':''}`} onClick={()=>setTab('broadcast')}>广播</button>
    </div>
    {tab === 'dashboard' ? <DashboardPanel stats={stats} /> :
     tab === 'broadcast' ? <BroadcastPanel addToast={addToast} /> :
     tab === 'users' ? <UserPanel users={users} addToast={addToast} onRefresh={loadTab} /> :
     <PostPanel posts={posts} addToast={addToast} onRefresh={loadTab} />}
  </>
}

function DashboardPanel({ stats }) {
  if (!stats) return <div className="spinner"/>
  const cards = [
    { label: '总用户数', value: stats.total_users, color: '#6366f1' },
    { label: '今日新增', value: stats.new_users_today, color: '#22c55e' },
    { label: '总帖子数', value: stats.total_posts, color: '#f59e0b' },
    { label: '今日发帖', value: stats.posts_today, color: '#3b82f6' },
    { label: '总评论数', value: stats.total_comments, color: '#ec4899' },
    { label: '在线用户', value: stats.online_count, color: '#14b8a6' },
  ]
  return <div className="stats-grid">
    {cards.map((c, i) => (
      <div key={i} className="stats-card" style={{borderLeft: `3px solid ${c.color}`}}>
        <div className="stats-label">{c.label}</div>
        <div className="stats-value" style={{color:c.color}}>{c.value}</div>
      </div>
    ))}
  </div>
}

function UserPanel({ users, addToast, onRefresh }) {
  async function delUser(id) { try { await api(`/admin/users/${id}`, { method: 'DELETE' }); addToast('已删除','success'); onRefresh() } catch(e) { addToast(e.message,'error') } }
  async function banUser(id, status) { try { await api(`/admin/users/${id}/status`, { method:'PUT', body:{status} }); addToast(status===1?'已解封':'已封禁','success'); onRefresh() } catch(e) { addToast(e.message,'error') } }
  async function setAdmin(id, adminType) { try { await api(`/admin/users/${id}/admin_type`, { method:'PUT', body:{admin_type:adminType} }); addToast('已设置','success'); onRefresh() } catch(e) { addToast(e.message,'error') } }
  if (!users) return <div className="spinner"/>
  if (users.length === 0) return <div className="empty">暂无用户</div>
  return users.map(u => (
    <div key={u.id} className="card flex" style={{justifyContent:'space-between',alignItems:'center',flexWrap:'wrap',gap:8}}>
      <span><b>{u.username}</b> <span className="text-muted text-sm">{u.email||''}</span> · {u.status===1?'正常':'<span style="color:var(--red)">禁用</span>'} {u.admin_type>0?<span style="color:var(--accent)">· 管理员</span>:''}</span>
      <div className="flex gap8">
        {u.admin_type>0 ? <button className="btn-outline" onClick={()=>setAdmin(u.id,0)}>取消管理</button> : <button className="btn-outline" onClick={()=>setAdmin(u.id,1)}>设为管理</button>}
        {u.status===1 ? <button className="btn-outline" style={{color:'var(--red)'}} onClick={()=>banUser(u.id,0)}>封禁</button> : <button className="btn-outline" style={{color:'var(--green)'}} onClick={()=>banUser(u.id,1)}>解封</button>}
        <button className="btn-danger" onClick={()=>delUser(u.id)}>删除</button>
      </div>
    </div>
  ))
}

function PostPanel({ posts, addToast, onRefresh }) {
  async function delPost(id) { try { await api(`/admin/posts/${id}`, { method: 'DELETE' }); addToast('已删除','success'); onRefresh() } catch(e) { addToast(e.message,'error') } }
  async function setTop(id, v) { try { await api(`/admin/posts/${id}/top`, { method:'PUT', body:{is_top:v} }); onRefresh() } catch(e) { addToast(e.message,'error') } }
  async function setEssence(id, v) { try { await api(`/admin/posts/${id}/essence`, { method:'PUT', body:{is_essence:v} }); onRefresh() } catch(e) { addToast(e.message,'error') } }
  if (!posts) return <div className="spinner"/>
  if (posts.length === 0) return <div className="empty">暂无帖子</div>
  return posts.map(p => (
    <div key={p.id} className="card flex" style={{justifyContent:'space-between',alignItems:'center'}}>
      <span>{p.title}</span>
      <div className="flex gap8">
        <button className="btn-outline" onClick={()=>setTop(p.id,!p.is_top)}>{p.is_top?'取消置顶':'置顶'}</button>
        <button className="btn-outline" onClick={()=>setEssence(p.id,!p.is_essence)}>{p.is_essence?'取消精华':'精华'}</button>
        <button className="btn-danger" onClick={()=>delPost(p.id)}>删除</button>
      </div>
    </div>
  ))
}

function BroadcastPanel({ addToast }) {
  const [msg, setMsg] = useState('')
  async function send() {
    if (!msg.trim()) return
    try { await api('/admin/broadcast', {method:'POST', body:{content:msg.trim()}}); addToast('广播已发送','success'); setMsg('') }
    catch(e) { addToast(e.message,'error') }
  }
  return <div className="card"><textarea value={msg} onChange={e=>setMsg(e.target.value)} placeholder="广播内容..." rows={3} className="mb16"/><button className="btn" onClick={send}>发送全站广播</button></div>
}
