import { useState, useEffect } from 'react'
import { api } from '../api'
import PostCard from '../components/PostCard'

export default function Feed({ addToast }) {
  const [posts, setPosts] = useState(null)
  const [sort, setSort] = useState('newest')
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)

  useEffect(() => { load() }, [sort, page])

  async function load() {
    setPosts(null)
    try {
      const data = await api(`/posts?page=${page}&page_size=15&sort=${sort}`)
      setPosts(data?.items || [])
      setTotal(data?.total || 0)
    } catch (err) { addToast(err.message, 'error'); setPosts([]) }
  }

  return (
    <>
      <div className="tabs">
        {['newest','hot','essence'].map(s => (
          <button key={s} className={`tab ${sort===s?'active':''}`} onClick={()=>{setSort(s);setPage(1)}}>
            { {newest:'最新',hot:'热门',essence:'精华'}[s] }
          </button>
        ))}
      </div>
      {posts === null ? <div className="spinner" /> :
       posts.length === 0 ? <div className="empty">还没有帖子</div> :
       posts.map(p => <PostCard key={p.id} post={p} addToast={addToast} />)}
      <div className="flex text-center" style={{justifyContent:'center',gap:12,marginTop:16}}>
        <button className="btn-outline" disabled={page<=1} onClick={()=>setPage(p=>p-1)}>上一页</button>
        <span className="text-sm text-muted">第 {page} 页</span>
        <button className="btn-outline" onClick={()=>setPage(p=>p+1)}>下一页</button>
      </div>
    </>
  )
}
