import { useState, useEffect, useRef } from 'react'
import { api, getUser } from '../api'

export default function Messages({ addToast }) {
  const [convs, setConvs] = useState([])
  const me = getUser()
  const [chat, setChat] = useState(null)
  const [msgs, setMsgs] = useState([])
  const [input, setInput] = useState('')
  const bubblesRef = useRef(null)

  useEffect(() => {
    api('/messages/conversations').then(d => setConvs(d?.items || [])).catch(() => {})
  }, [])

  async function openChat(uid, name) {
    setChat({ id: uid, name })
    try {
      const data = await api(`/messages?sender_id=${uid}&page=1&page_size=50`)
      setMsgs((data?.items || []).reverse())
    } catch (err) { addToast(err.message, 'error') }
  }

  async function send() {
    if (!input.trim() || !chat) return
    try {
      await api('/messages', { method: 'POST', body: { receiver_id: chat.id, content: input.trim() } })
      setMsgs(m => [...m, { sender_id: 'me', content: input.trim(), id: Date.now() }])
      setInput('')
      setTimeout(() => bubblesRef.current?.scrollTo(0, bubblesRef.current.scrollHeight), 50)
    } catch (err) { addToast(err.message, 'error') }
  }

  return (
    <div className="msg-layout">
      <div className="msg-list">
        {convs.length === 0 ? <div className="empty">暂无会话</div> :
         convs.map(c => (
           <div key={c.user_id} className={`msg-item ${chat?.id===c.user_id?'active':''}`} onClick={()=>openChat(c.user_id,c.nickname||c.username)}>
             <div className="avatar" style={{width:36,height:36,fontSize:14}}>{(c.nickname||c.username||'?')[0].toUpperCase()}</div>
             <div className="flex1" style={{minWidth:0}}>
               <div className="flex" style={{justifyContent:'space-between'}}>
                 <span style={{fontWeight:600,fontSize:13}}>{c.nickname||c.username}</span>
                 <span className="text-muted text-sm">{c.last_time?.slice(-8)||''}</span>
               </div>
               <div className="preview">{c.last_message||''}</div>
             </div>
           </div>
         ))}
      </div>
      <div className="msg-chat">
        {!chat ? <div style={{flex:1,display:'flex',alignItems:'center',justifyContent:'center',color:'var(--text3)'}}>选择一个会话</div> : <>
          <div style={{padding:'10px 16px',borderBottom:'1px solid var(--border)',fontWeight:600}}>{chat.name}</div>
          <div className="bubbles" ref={bubblesRef}>
            {msgs.map(m => (
              <div key={m.id} className={`bubble ${m.sender_id===me.user_id?'mine':'theirs'}`}>{m.content}</div>
            ))}
          </div>
          <div className="input-bar">
            <input value={input} onChange={e=>setInput(e.target.value)} onKeyDown={e=>e.key==='Enter'&&send()} placeholder="输入消息..." />
            <button className="btn btn-sm" onClick={send}>发送</button>
          </div>
        </>}
      </div>
    </div>
  )
}
