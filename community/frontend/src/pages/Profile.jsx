import { useState, useEffect } from 'react'
import { api, getUser } from '../api'
import Follows from './Follows'

export default function Profile({ addToast }) {
  const [profile, setProfile] = useState(null)
  const [editing, setEditing] = useState(false)
  const [form, setForm] = useState({})
  const [showFollows, setShowFollows] = useState(false)
  const [counts, setCounts] = useState({ followers: 0, following: 0 })

  useEffect(() => { load() }, [])

  async function load() {
    try {
      const p = await api('/users/profile'); setProfile(p); setForm({ nickname: p.nickname||'', bio: p.bio||'', avatar: p.avatar||'' })
      const c = await api(`/follows/${p.id}/counts`); setCounts({ followers: c?.followers||0, following: c?.following||0 })
    } catch (e) { addToast(e.message,'error') }
  }

  async function save() {
    try { await api('/users/profile', {method:'PUT', body:form}); addToast('已保存','success'); setEditing(false); load() }
    catch (e) { addToast(e.message,'error') }
  }

  if (!profile) return <div className="spinner" />
  return (
    <>
      <div className="card">
        <div className="flex" style={{justifyContent:'space-between',marginBottom:16}}>
          <div className="avatar" style={{width:80,height:80,fontSize:32}}>{(profile.nickname||profile.username)[0].toUpperCase()}</div>
          <button className="btn-outline" onClick={()=>setEditing(!editing)}>{editing?'取消':'编辑资料'}</button>
        </div>
        {editing ? <>
          <div className="field"><label>昵称</label><input value={form.nickname||''} onChange={e=>setForm(f=>({...f,nickname:e.target.value}))} /></div>
          <div className="field"><label>简介</label><textarea value={form.bio||''} onChange={e=>setForm(f=>({...f,bio:e.target.value}))} rows={3} /></div>
          <button className="btn" onClick={save}>保存</button>
        </> : <>
          <h2>{profile.nickname || profile.username}</h2>
          <p className="text-muted">@{profile.username}</p>
          <p style={{marginTop:8}}>{profile.bio || '这个人很懒，什么都没写'}</p>
        </>}
      </div>
      <div className="card" style={{display:'flex',gap:32,justifyContent:'center'}}>
        <button style={{background:'none',color:'var(--text)',textAlign:'center'}} onClick={()=>setShowFollows(!showFollows)}>
          <div style={{fontSize:20,fontWeight:700}}>{counts.followers}</div>
          <div className="text-muted text-sm">粉丝</div>
        </button>
        <button style={{background:'none',color:'var(--text)',textAlign:'center'}} onClick={()=>setShowFollows(!showFollows)}>
          <div style={{fontSize:20,fontWeight:700}}>{counts.following}</div>
          <div className="text-muted text-sm">关注</div>
        </button>
      </div>
      {showFollows && <Follows addToast={addToast} userId={profile.id} onClose={()=>setShowFollows(false)} />}
    </>
  )
}
