import { useState, useEffect } from 'react'
import { api } from '../api'

export default function Follows({ addToast, userId, onClose }) {
  const [tab, setTab] = useState('followers')
  const [items, setItems] = useState(null)

  useEffect(() => { load() }, [tab])

  async function load() {
    setItems(null)
    try {
      const endpoint = tab === 'followers' ? '/follows/followers' : '/follows/following'
      const data = await api(`${endpoint}?user_id=${userId}&page=1&page_size=100`)
      setItems(data?.items || [])
    } catch (err) { addToast(err.message, 'error'); setItems([]) }
  }

  async function unfollow(uid) {
    try { await api(`/follows/${uid}`, { method: 'DELETE' }); load(); addToast('已取消关注','success') }
    catch (err) { addToast(err.message, 'error') }
  }

  return (
    <div className="card" style={{marginTop:0}}>
      <div className="flex" style={{justifyContent:'space-between',alignItems:'center',marginBottom:16}}>
        <div className="tabs" style={{margin:0}}>
          <button className={`tab ${tab==='followers'?'active':''}`} onClick={()=>setTab('followers')}>粉丝</button>
          <button className={`tab ${tab==='following'?'active':''}`} onClick={()=>setTab('following')}>关注</button>
        </div>
        <button className="btn-outline" onClick={onClose}>关闭</button>
      </div>
      {items === null ? <div className="spinner" /> :
       items.length === 0 ? <div className="empty">暂无</div> :
       items.map(u => (
         <div key={u.user_id||u.id} className="flex" style={{justifyContent:'space-between',alignItems:'center',padding:'8px 0',borderBottom:'1px solid var(--border)'}}>
           <span style={{fontWeight:600}}>{u.nickname || u.username}</span>
           <button className="btn-outline" onClick={()=>unfollow(u.user_id||u.id)}>取消关注</button>
         </div>
       ))}
    </div>
  )
}
