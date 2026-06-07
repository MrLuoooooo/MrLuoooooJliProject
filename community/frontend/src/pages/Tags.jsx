import { useState, useEffect } from 'react'
import { api } from '../api'

export default function Tags({ addToast }) {
  const [tags, setTags] = useState(null)
  const [name, setName] = useState('')
  const [editing, setEditing] = useState(null)
  const [editName, setEditName] = useState('')

  useEffect(() => { load() }, [])

  async function load() {
    try { const d = await api('/tags?page=1&page_size=100'); setTags(d?.items||[]) }
    catch (e) { addToast(e.message,'error'); setTags([]) }
  }

  async function create() {
    if (!name.trim()) return
    try { await api('/tags', {method:'POST', body:{name:name.trim()}}); setName(''); load(); addToast('创建成功','success') }
    catch (e) { addToast(e.message,'error') }
  }

  async function updateTag(id) {
    if (!editName.trim()) return
    try { await api(`/tags/${id}`, {method:'PUT', body:{name:editName.trim()}}); setEditing(null); load() }
    catch (e) { addToast(e.message,'error') }
  }

  async function deleteTag(id) {
    try { await api(`/tags/${id}`, {method:'DELETE'}); load(); addToast('已删除','success') }
    catch (e) { addToast(e.message,'error') }
  }

  return <div className="card">
    <h3 style={{marginBottom:16}}>标签管理</h3>
    <div className="flex gap8 mb16">
      <input value={name} onChange={e=>setName(e.target.value)} onKeyDown={e=>e.key==='Enter'&&create()} placeholder="标签名" className="flex1" />
      <button className="btn btn-sm" onClick={create}>创建</button>
    </div>
    {tags === null ? <div className="spinner"/> :
     tags.length === 0 ? <div className="empty">暂无标签</div> :
     tags.map(t => (
       <div key={t.id} className="flex" style={{justifyContent:'space-between',alignItems:'center',padding:'8px 0',borderBottom:'1px solid var(--border)'}}>
         {editing===t.id ? (
           <div className="flex gap8 flex1"><input value={editName} onChange={e=>setEditName(e.target.value)} onKeyDown={e=>e.key==='Enter'&&updateTag(t.id)}/><button className="btn btn-sm" onClick={()=>updateTag(t.id)}>保存</button></div>
         ) : (
           <span>{t.name} <span className="text-muted text-sm">{t.post_count||0} 帖</span></span>
         )}
         <div className="flex gap8">
           <button className="btn-outline" onClick={()=>{setEditing(t.id);setEditName(t.name)}}>编辑</button>
           <button className="btn-danger" onClick={()=>deleteTag(t.id)}>删除</button>
         </div>
       </div>
     ))}
  </div>
}
