import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { useAuth } from '../context/AuthContext'
import { attendeesAPI, speakersAPI, sessionsAPI } from '../services/api'
import AttendeeList from '../components/admin/AttendeeList'
import SpeakerManagement from '../components/admin/SpeakerManagement'
import SessionManagement from '../components/admin/SessionManagement'
import Analytics from '../components/admin/Analytics'
import { LogOut, Users, Mic, Calendar, BarChart3 } from 'lucide-react'

function AdminDashboard() {
  const [activeTab, setActiveTab] = useState('attendees')
  const { logout } = useAuth()
  const navigate = useNavigate()

  useEffect(() => {
    // Check authentication on mount
    const authStatus = localStorage.getItem('adminAuthenticated')
    if (authStatus !== 'true') {
      navigate('/admin/login')
    }
  }, [navigate])

  const handleLogout = () => {
    logout()
    navigate('/admin/login')
  }

  const tabs = [
    { id: 'attendees', label: 'Attendees', icon: Users },
    { id: 'speakers', label: 'Speakers', icon: Mic },
    { id: 'sessions', label: 'Sessions', icon: Calendar },
    { id: 'analytics', label: 'Analytics', icon: BarChart3 },
  ]

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold text-gray-900">
            Admin Dashboard
          </h1>
          <button
            onClick={handleLogout}
            className="flex items-center gap-2 px-4 py-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
          >
            <LogOut className="w-5 h-5" />
            Logout
          </button>
        </div>
      </header>

      {/* Tabs */}
      <div className="bg-white border-b">
        <div className="max-w-7xl mx-auto px-4">
          <div className="flex space-x-1 overflow-x-auto">
            {tabs.map((tab) => {
              const Icon = tab.icon
              return (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`flex items-center gap-2 px-6 py-4 font-semibold transition-colors whitespace-nowrap ${
                    activeTab === tab.id
                      ? 'text-blue-600 border-b-2 border-blue-600'
                      : 'text-gray-600 hover:text-gray-900'
                  }`}
                >
                  <Icon className="w-5 h-5" />
                  {tab.label}
                </button>
              )
            })}
          </div>
        </div>
      </div>

      {/* Content */}
      <main className="max-w-7xl mx-auto px-4 py-8">
        <motion.div
          key={activeTab}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3 }}
        >
          {activeTab === 'attendees' && <AttendeeList />}
          {activeTab === 'speakers' && <SpeakerManagement />}
          {activeTab === 'sessions' && <SessionManagement />}
          {activeTab === 'analytics' && <Analytics />}
        </motion.div>
      </main>
    </div>
  )
}

export default AdminDashboard

