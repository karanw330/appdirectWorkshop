import { useState, useEffect } from 'react'
import { attendeesAPI } from '../../services/api'
import { PieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip } from 'recharts'

const COLORS = ['#3b82f6', '#8b5cf6', '#ec4899', '#f59e0b', '#10b981', '#ef4444', '#06b6d4', '#84cc16', '#6366f1']

function Analytics() {
  const [attendees, setAttendees] = useState([])
  const [loading, setLoading] = useState(true)
  const [chartData, setChartData] = useState([])

  useEffect(() => {
    fetchAttendees()
  }, [])

  const fetchAttendees = async () => {
    try {
      const response = await attendeesAPI.getAll()
      setAttendees(response.data)
      
      // Calculate designation breakdown
      const designationCount = {}
      response.data.forEach((attendee) => {
        designationCount[attendee.designation] = 
          (designationCount[attendee.designation] || 0) + 1
      })
      
      const data = Object.entries(designationCount).map(([name, value]) => ({
        name,
        value,
      }))
      
      setChartData(data)
    } catch (error) {
      console.error('Error fetching attendees:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="text-center py-12">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
      </div>
    )
  }

  return (
    <div className="bg-white rounded-xl shadow-lg p-6">
      <h2 className="text-2xl font-bold text-gray-900 mb-6">Analytics</h2>
      
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div>
          <h3 className="text-lg font-semibold text-gray-700 mb-4">
            Attendee Breakdown by Designation
          </h3>
          {chartData.length > 0 ? (
            <ResponsiveContainer width="100%" height={400}>
              <PieChart>
                <Pie
                  data={chartData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(0)}%`}
                  outerRadius={120}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {chartData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          ) : (
            <div className="text-center py-12 text-gray-500">
              No data available
            </div>
          )}
        </div>
        
        <div>
          <h3 className="text-lg font-semibold text-gray-700 mb-4">
            Statistics
          </h3>
          <div className="space-y-4">
            <div className="bg-blue-50 rounded-lg p-4">
              <div className="text-3xl font-bold text-blue-600 mb-1">
                {attendees.length}
              </div>
              <div className="text-sm text-gray-600">Total Attendees</div>
            </div>
            
            <div className="bg-purple-50 rounded-lg p-4">
              <div className="text-3xl font-bold text-purple-600 mb-1">
                {chartData.length}
              </div>
              <div className="text-sm text-gray-600">Unique Designations</div>
            </div>
            
            <div className="bg-indigo-50 rounded-lg p-4">
              <h4 className="font-semibold text-gray-700 mb-3">Top Designations</h4>
              <div className="space-y-2">
                {chartData
                  .sort((a, b) => b.value - a.value)
                  .slice(0, 5)
                  .map((item, index) => (
                    <div key={item.name} className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <div
                          className="w-4 h-4 rounded"
                          style={{ backgroundColor: COLORS[index % COLORS.length] }}
                        ></div>
                        <span className="text-sm text-gray-700">{item.name}</span>
                      </div>
                      <span className="text-sm font-semibold text-gray-900">
                        {item.value}
                      </span>
                    </div>
                  ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Analytics

