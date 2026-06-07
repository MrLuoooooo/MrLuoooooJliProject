import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../api'

export default function CreatePost({ addToast }) {
  const navigate = useNavigate()
  const [form, setForm] = useState({ title: '', content: '', summary: '', cover_image: '', category_id: '' })

  function update(k, v) { setForm(f => ({ ...f, [k]: v })) }

  async function submit(e) {
    e.preventDefault()
    if (!form.title.trim() || !form.content.trim()) { addToast('标题和内容不能为空', 'error'); return }
    try {
      await api('/posts', { method: 'POST', body: { ...form, category_id: parseInt(form.category_id) || 0 } })
      addToast('发布成功', 'success')
      navigate('/feed')
    } catch (err) { addToast(err.message, 'error') }
  }

  return (
    <form className="card" onSubmit={submit}>
      <div className="field"><label>标题</label><input value={form.title} onChange={e=>update('title',e.target.value)} maxLength={200} placeholder="给帖子起个标题" /></div>
      <div className="field"><label>摘要</label><input value={form.summary} onChange={e=>update('summary',e.target.value)} maxLength={500} placeholder="一句话概括（可选）" /></div>
      <div className="field"><label>内容</label><textarea value={form.content} onChange={e=>update('content',e.target.value)} rows={8} placeholder="写点什么..." /></div>
      <div style={{display:'grid',gridTemplateColumns:'1fr 1fr',gap:16,marginBottom:16}}>
        <div className="field" style={{margin:0}}><label>分类ID</label><input value={form.category_id} onChange={e=>update('category_id',e.target.value)} type="number" placeholder="可选" /></div>
        <div className="field" style={{margin:0}}><label>封面图URL</label><input value={form.cover_image} onChange={e=>update('cover_image',e.target.value)} placeholder="可选" /></div>
      </div>
      <button className="btn" type="submit">发布</button>
    </form>
  )
}
