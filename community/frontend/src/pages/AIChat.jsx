import { useState, useRef, useEffect } from 'react'
import { api } from '../api'

export default function AIChat({ addToast }) {
  const [messages, setMessages] = useState(() => {
    const saved = localStorage.getItem('ai_messages')
    return saved ? JSON.parse(saved) : [
      { role: 'system', content: '你好！我是 AI 助手，有什么可以帮你的？' }
    ]
  })
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)
  const bottomRef = useRef(null)

  useEffect(() => {
    localStorage.setItem('ai_messages', JSON.stringify(messages))
  }, [messages])

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  async function send() {
    if (!input.trim() || loading) return
    const text = input
    setInput('')
    setMessages(prev => [...prev, { role: 'user', content: text }])
    setLoading(true)
    try {
      const data = await api('/ai/chat', {
        method: 'POST',
        body: {
          messages: [{ role: 'user', content: text }],
          stream: false
        }
      })
      setMessages(prev => [...prev, { role: 'assistant', content: data.content }])
    } catch (err) {
      addToast(err.message, 'error')
      setMessages(prev => [...prev, { role: 'assistant', content: 'AI 回复失败：' + err.message }])
    }
    setLoading(false)
  }

  return (
    <div className="chat-page">
      <div className="chat-messages">
        {messages.map((m, i) => (
          <div key={i} className={`msg ${m.role}`}>
            <div className="msg-avatar">{m.role === 'assistant' ? '🤖' : '👤'}</div>
            <div className="msg-bubble">{m.content}</div>
          </div>
        ))}
        {loading && (
          <div className="msg assistant">
            <div className="msg-avatar">🤖</div>
            <div className="msg-bubble"><span className="typing">思考中...</span></div>
          </div>
        )}
        <div ref={bottomRef} />
      </div>
      <div className="chat-input" onKeyDown={e => e.key === 'Enter' && !e.shiftKey && (e.preventDefault(), send())}>
        <textarea value={input} onChange={e => setInput(e.target.value)} placeholder="输入消息..." rows={2} />
        <button className="btn" onClick={send} disabled={loading}>发送</button>
      </div>
    </div>
  )
}
