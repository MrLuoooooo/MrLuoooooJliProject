import { useState, useEffect, useCallback } from 'react'
import { Routes, Route, Navigate, useNavigate } from 'react-router-dom'
import { api, setToken, setUser, getUser, getToken, clearAuth, getWsUrl } from './api'
import Layout from './components/Layout'
import Toast from './components/Toast'
import AuthForm from './components/AuthForm'
import Feed from './pages/Feed'
import AIChat from './pages/AIChat'
import PostDetail from './pages/PostDetail'
import CreatePost from './pages/CreatePost'
import Messages from './pages/Messages'
import Notifications from './pages/Notifications'
import Search from './pages/Search'
import Profile from './pages/Profile'
import Tags from './pages/Tags'
import Admin from './pages/Admin'

export default function App() {
  const [auth, setAuth] = useState(false)
  const [loading, setLoading] = useState(true)
  const [toasts, setToasts] = useState([])
  const navigate = useNavigate()

  const addToast = useCallback((msg, type = 'info') => {
    const id = Date.now()
    setToasts(p => [...p, { id, msg, type }])
    setTimeout(() => setToasts(p => p.filter(t => t.id !== id)), 2500)
  }, [])

  useEffect(() => {
    console.log('[APP] init: token=' + (getToken() ? 'YES' : 'NO'))
    if (getToken()) {
      api('/users/profile')
        .then(u => { console.log('[APP] profile loaded', u); setUser(u); setAuth(true) })
        .catch(e => { console.error('[APP] profile error:', e.message); clearAuth() })
        .finally(() => setLoading(false))
    } else {
      setLoading(false)
    }
  }, [])
  // 登录成功后自动跳转
  useEffect(() => {
    console.log('[APP] auth effect: auth=' + auth + ' loading=' + loading)
    if (auth && !loading) {
      console.log('[APP] navigating to /feed')
      navigate('/feed', { replace: true })
    }
  }, [auth, loading])

  function onLogin(userData, tok) {
    console.log('[APP] onLogin called')
    setToken(tok)
    setUser(userData)
    setAuth(true)
    addToast('登录成功', 'success')
    console.log('[APP] onLogin: redirecting to /feed')
    window.location.href = '/feed'
  }

  function onLogout() {
    clearAuth()
    setAuth(false)
  }

  if (loading) return <div style={{display:'flex',alignItems:'center',justifyContent:'center',height:'100vh',background:'var(--bg)'}}><div className="spinner"/></div>
  if (!auth) return <AuthForm onLogin={onLogin} addToast={addToast} />

  return (
    <>
      <Toast toasts={toasts} />
      <Layout user={getUser()} onLogout={onLogout} addToast={addToast}>
        <Routes>
          <Route path="/feed" element={<Feed addToast={addToast} />} />
          <Route path="/post" element={<PostDetail addToast={addToast} />} />
          <Route path="/create" element={<CreatePost addToast={addToast} />} />
          <Route path="/messages" element={<Messages addToast={addToast} />} />
          <Route path="/notifications" element={<Notifications addToast={addToast} />} />
          <Route path="/search" element={<Search addToast={addToast} />} />
          <Route path="/profile" element={<Profile addToast={addToast} />} />
          <Route path="/tags" element={<Tags addToast={addToast} />} />
          <Route path="/ai" element={<AIChat addToast={addToast} />} />
          <Route path="/admin" element={<Admin addToast={addToast} />} />
          <Route path="*" element={<Navigate to="/feed" />} />
        </Routes>
      </Layout>
    </>
  )
}
