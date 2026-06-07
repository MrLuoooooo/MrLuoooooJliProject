import { useState, useEffect, useRef } from 'react'
import { useSearchParams } from 'react-router-dom'
import { api, getUser } from '../api'

export default function PostDetail({ addToast }) {
  const [params] = useSearchParams()
  const id = params.get('id')
  const [post, setPost] = useState(null)
  const [comments, setComments] = useState(null)
  const [commentText, setCommentText] = useState('')
  const [liked, setLiked] = useState(false)
  const [favorited, setFavorited] = useState(false)
  const inputRef = useRef(null)

  useEffect(() => { loadPost(); loadComments() }, [id])

  async function loadPost() {
    try {
      const p = await api(`/posts/${id}`)
      setPost(p)
      setLiked(p.liked || false)
      setFavorited(p.favorited || false)
    } catch (err) { addToast(err.message, 'error') }
  }

  async function loadComments() {
    try {
      const data = await api(`/comments?post_id=${id}&page=1&page_size=100`)
      setComments(data?.items || [])
    } catch (_) { setComments([]) }
  }

  async function toggleLike() {
    try {
      if (liked) { await api(`/posts/${id}/like`, { method: 'DELETE' }); setLiked(false) }
      else { await api(`/posts/${id}/like`, { method: 'POST' }); setLiked(true) }
    } catch (err) { addToast(err.message, 'error') }
  }

  async function toggleFav() {
    try {
      if (favorited) { await api(`/posts/${id}/favorite`, { method: 'DELETE' }); setFavorited(false) }
      else { await api(`/posts/${id}/favorite`, { method: 'POST' }); setFavorited(true) }
      addToast(favorited ? '已取消收藏' : '已收藏', 'success')
    } catch (err) { addToast(err.message, 'error') }
  }

  async function submitComment() {
    if (!commentText.trim()) return
    try {
      await api('/comments', { method: 'POST', body: { post_id: parseInt(id), content: commentText.trim() } })
      setCommentText('')
      loadComments()
      addToast('评论成功', 'success')
    } catch (err) { addToast(err.message, 'error') }
  }

  async function likeComment(cid, liked) {
    try {
      if (liked) await api(`/comments/${cid}/like`, { method: 'DELETE' })
      else await api(`/comments/${cid}/like`, { method: 'POST' })
      loadComments()
    } catch (_) {}
  }

  if (!post) return <div className="spinner" />
  const u = getUser()
  return (
    <div>
      <div className="card">
        <div className="header">
          <div className="avatar" style={{width:48,height:48,fontSize:20}}>{(post.nickname||post.username||'?')[0].toUpperCase()}</div>
          <div>
            <div className="name">{post.nickname || post.username}</div>
            <div className="time">{post.created_at} · {post.view_count || 0} 次浏览</div>
          </div>
        </div>
        <h2 style={{marginBottom:12}}>{post.title}</h2>
        <div style={{lineHeight:1.8,whiteSpace:'pre-wrap',color:'var(--text2)',marginBottom:16}}>{post.content}</div>
        <div className="actions">
          <button onClick={toggleLike} className={liked ? 'liked' : ''}>❤️ {post.like_count}</button>
          <button onClick={toggleFav} className={favorited ? 'on' : ''}>⭐ {favorited ? '已收藏' : '收藏'}</button>
          <button>💬 {post.comment_count}</button>
        </div>
      </div>

      <div className="card" style={{marginTop:0}}>
        <h3 style={{marginBottom:16}}>评论 ({comments?.length || 0})</h3>
        <div className="flex gap8" style={{marginBottom:16}}>
          <input ref={inputRef} value={commentText} onChange={e=>setCommentText(e.target.value)} onKeyDown={e=>e.key==='Enter'&&submitComment()} placeholder="写评论..." />
          <button className="btn btn-sm" onClick={submitComment}>发送</button>
        </div>
        {comments === null ? <div className="spinner" /> :
         comments.length === 0 ? <div className="empty" style={{padding:0}}>暂无评论</div> :
         comments.map(c => (
           <div key={c.id} style={{borderTop:'1px solid var(--border)',padding:'12px 0'}}>
             <div className="flex" style={{justifyContent:'space-between',marginBottom:4}}>
               <span style={{fontWeight:600,fontSize:13}}>{c.username || '用户'}</span>
               <span className="text-muted text-sm">{c.created_at}</span>
             </div>
             <div style={{fontSize:13,color:'var(--text2)'}}>{c.content}</div>
             <button onClick={()=>likeComment(c.id,c.liked)} style={{background:'none',color:c.liked?'var(--red)':'var(--text3)',fontSize:11,padding:'2px 0',marginTop:4}}>
               {c.liked ? '❤️' : '🤍'} {c.like_count||0}
             </button>
           </div>
         ))}
      </div>
    </div>
  )
}
