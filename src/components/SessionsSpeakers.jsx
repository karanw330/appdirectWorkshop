import { useEffect, useState } from 'react'
import { motion } from 'framer-motion'
import { sessionsAPI, speakersAPI } from '../services/api'
import { Calendar, Clock, User } from 'lucide-react'

function SessionsSpeakers() {
  const [sessions, setSessions] = useState([])
  const [speakers, setSpeakers] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [sessionsRes, speakersRes] = await Promise.all([
          sessionsAPI.getAll(),
          speakersAPI.getAll(),
        ])
        // Ensure we always have arrays, even if API returns null/undefined
        setSessions(Array.isArray(sessionsRes?.data) ? sessionsRes.data : [])
        setSpeakers(Array.isArray(speakersRes?.data) ? speakersRes.data : [])
      } catch (error) {
        console.error('Error fetching data:', error)
        // Set empty arrays on error to prevent crashes
        setSessions([])
        setSpeakers([])
      } finally {
        setLoading(false)
      }
    }
    fetchData()
  }, [])

  const getSpeakerById = (speakerId) => {
    return speakers.find((s) => s.id === speakerId)
  }

  if (loading) {
    return (
      <section id="sessions" className="py-20 px-4 bg-white">
        <div className="max-w-7xl mx-auto text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
        </div>
      </section>
    )
  }

  return (
    <section id="sessions" className="py-20 px-4 bg-white">
      <div className="max-w-7xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.6 }}
          className="text-center mb-12"
        >
          <h2 className="text-4xl md:text-5xl font-bold text-gray-900 mb-4">
            Sessions & Speakers
          </h2>
          <p className="text-xl text-gray-600">
            Explore our exciting lineup of AI workshops and expert speakers
          </p>
        </motion.div>

        {sessions.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-600 text-lg">
              No sessions available yet. Check back soon!
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {sessions.map((session, index) => {
              if (!session || !session.id) return null
              const speaker = session.speakerId ? getSpeakerById(session.speakerId) : null
              return (
                <motion.div
                  key={session.id}
                  initial={{ opacity: 0, y: 20 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  viewport={{ once: true }}
                  transition={{ duration: 0.6, delay: index * 0.1 }}
                  className="bg-gradient-to-br from-blue-50 to-indigo-50 rounded-xl p-6 shadow-lg hover:shadow-xl transition-all duration-300 border border-blue-100"
                >
                  <div className="flex items-center gap-2 mb-4">
                    <Calendar className="w-5 h-5 text-blue-600" />
                    <span className="text-sm font-semibold text-blue-600">
                      {session.date || 'TBA'}
                    </span>
                  </div>
                  <h3 className="text-2xl font-bold text-gray-900 mb-3">
                    {session.title || 'Untitled Session'}
                  </h3>
                  <p className="text-gray-600 mb-4">{session.description || 'No description available.'}</p>
                  <div className="flex items-center gap-2 mb-4">
                    <Clock className="w-4 h-4 text-gray-500" />
                    <span className="text-sm text-gray-600">{session.time || 'TBA'}</span>
                  </div>
                  {speaker && (
                    <div className="mt-4 pt-4 border-t border-blue-200">
                      <div className="flex items-center gap-3">
                        <div className="w-12 h-12 bg-gradient-to-br from-blue-500 to-purple-500 rounded-full flex items-center justify-center text-white font-bold">
                          {speaker.name?.charAt(0) || '?'}
                        </div>
                        <div>
                          <div className="flex items-center gap-2">
                            <User className="w-4 h-4 text-gray-500" />
                            <span className="font-semibold text-gray-900">
                              {speaker.name || 'Unknown Speaker'}
                            </span>
                          </div>
                          <p className="text-sm text-gray-600">{speaker.bio || 'No bio available.'}</p>
                        </div>
                      </div>
                    </div>
                  )}
                </motion.div>
              )
            })}
          </div>
        )}
      </div>
    </section>
  )
}

export default SessionsSpeakers

