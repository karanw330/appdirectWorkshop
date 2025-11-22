import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Lock } from 'lucide-react'

function Footer() {
  return (
    <footer className="bg-gray-900 text-white py-12 px-4">
      <div className="max-w-7xl mx-auto">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-8">
          <div>
            <h3 className="text-xl font-bold mb-4">AppDirect India</h3>
            <p className="text-gray-400">
              Empowering businesses through innovative AI solutions and workshops.
            </p>
          </div>
          <div>
            <h3 className="text-xl font-bold mb-4">Quick Links</h3>
            <ul className="space-y-2 text-gray-400">
              <li>
                <a href="#sessions" className="hover:text-white transition-colors">
                  Sessions
                </a>
              </li>
              <li>
                <a
                  href="#registration"
                  className="hover:text-white transition-colors"
                >
                  Register
                </a>
              </li>
              <li>
                <a href="#location" className="hover:text-white transition-colors">
                  Location
                </a>
              </li>
            </ul>
          </div>
          <div>
            <h3 className="text-xl font-bold mb-4">Admin</h3>
            <Link
              to="/admin/login"
              className="inline-flex items-center gap-2 text-gray-400 hover:text-white transition-colors"
            >
              <Lock className="w-4 h-4" />
              Admin Login
            </Link>
          </div>
        </div>
        <div className="border-t border-gray-800 pt-8 text-center text-gray-400">
          <p>&copy; 2024 AppDirect India. All rights reserved.</p>
        </div>
      </div>
    </footer>
  )
}

export default Footer

