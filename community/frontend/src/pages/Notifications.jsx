import { useState, useEffect } from 'react'
import { api } from '../api'

export default function Notifications({ addToast }) {
  const [items, setItems] = useState(null)

  useEffect(() => {
    api('/notifications/read-all', { method: 'PUT' }).catch(() => {})
    api('/notifications?page=1&page_size=50')
      .then(d => setItems(d?.items || []))
      .catch(e => { addToast(e.message, 'error'); setItems([]) })
  }, [])

  if (items === null) return <div className="spinner" />
  if (items.length === 0) return <div className="empty">暂无通知</div>
  return items.map(n => (
    <div key={n.id} className="card" style={{padding:'12px 16px', opacity: n.is_read ? .5 : 1}}>
      <div className="flex" style={{justifyContent:'space-between',alignItems:'center'}}>
        <span><b>{n.from_name || '系统'}</b>: {n.content}</span>
        <span className="text-muted text-sm">{n.created_at}</span>
      </div>
    </div>
  ))
}
