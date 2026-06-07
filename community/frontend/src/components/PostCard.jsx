import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../api'

export default function PostCard({ post, addToast }) {
  const [liked, setLiked] = useState(post.liked || false)
  const [favorited, setFavorited] = useState(post.favorited || false)
  const [likeCount, setLikeCount] = useState(post.like_count || 0)
  const navigate = useNavigate()

  async function toggleLike(e) { e.stopPropagation()
    try {
      if (liked) { await api(`/posts/${post.id}/like`, { method: 'DELETE' }); setLiked(false); setLikeCount(c => c - 1) }
      else { await api(`/posts/${post.id}/like`, { method: 'POST' }); setLiked(true); setLikeCount(c => c + 1) }
    } catch (err) { addToast(err.message, 'error') }
  }

  async function toggleFav(e) { e.stopPropagation()
    try {
      if (favorited) { await api(`/posts/${post.id}/favorite`, { method: 'DELETE' }); setFavorited(false) }
      else { await api(`/posts/${post.id}/favorite`, { method: 'POST' }); setFavorited(true) }
      addToast(favorited ? '已取消收藏' : '已收藏', 'success')
    } catch (err) { addToast(err.message, 'error') }
  }

  function goDetail() { navigate(`/post?id=${post.id}`) }

  return (
    <div className="card">
      <div className="header">
        <div className="avatar">{(post.nickname || post.username || '?')[0].toUpperCase()}</div>
        <div>
          <div className="name">{post.nickname || post.username}</div>
          <div className="time">{post.created_at}</div>
        </div>
      </div>
      <div className="title" onClick={goDetail}>{post.title}</div>
      <div className="body">{post.summary || post.content?.slice(0, 200) || ''}</div>
      {post.tags?.length > 0 && <div className="feed-meta">{post.tags.map(t=><span key={t.id||t} className="feed-tag">{t.name||t}</span>)}</div>}
      <div className="actions">
        <button onClick={toggleLike} className={liked ? 'liked' : ''}>❤️ {likeCount}</button>
        <button onClick={goDetail}>💬 {post.comment_count || 0}</button>
        <button onClick={toggleFav} className={favorited ? 'on' : ''}>⭐ {favorited ? '已收藏' : '收藏'}</button>
        <button>👁️ {post.view_count || 0}</button>
      </div>
    </div>
  )
}
