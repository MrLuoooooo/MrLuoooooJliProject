import { useState } from 'react'
import { api } from '../api'

export default function Search({ addToast }) {
  const [kw, setKw] = useState('')
  const [type, setType] = useState('post')
  const [results, setResults] = useState(null)

  async function search() {
    if (!kw.trim()) return
    setResults(null)
    try {
      const data = await api(`/search?keyword=${kw}&type=${type}&page=1&page_size=20`)
      setResults(type === 'user' ? (data?.users || []) : (data?.posts || []))
    } catch (err) { addToast(err.message, 'error'); setResults([]) }
  }

  async function followUser(uid) {
    try { await api('/follows', { method: 'POST', body: { follow_id: uid } }); addToast('已关注', 'success') }
    catch (err) { addToast(err.message, 'error') }
  }

  return <>
    <div className="flex gap8 mb16">
      <input value={kw} onChange={e=>setKw(e.target.value)} onKeyDown={e=>e.key==='Enter'&&search()} placeholder="搜索帖子或用户..." className="flex1" />
      <select value={type} onChange={e=>setType(e.target.value)} style={{width:100}}><option value="post">帖子</option><option value="user">用户</option></select>
      <button className="btn" onClick={search}>搜索</button>
    </div>
    {results === null ? null : results.length === 0 ? <div className="empty">无结果</div> :
     type === 'post' ? results.map(p => (
       <div key={p.id} className="card">
         <div className="header"><div className="avatar">{(p.nickname||p.username||'?')[0].toUpperCase()}</div><div className="name">{p.nickname||p.username}</div></div>
         <div className="title">{p.title}</div>
       </div>)) :
     results.map(u => (
       <div key={u.id} className="card flex" style={{justifyContent:'space-between',alignItems:'center'}}>
         <div>
           <div><strong>{u.nickname||u.username}</strong> <span style={{fontSize:11,color:'var(--text3)'}}>#{u.id}</span></div>
           <div style={{fontSize:12,color:'var(--text2)'}}>{u.bio || u.username}</div>
         </div>
         <button className="btn btn-sm" onClick={()=>followUser(u.id)}>关注</button>
       </div>))}
  </>
}
