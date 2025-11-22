import { motion } from 'framer-motion'
import { MapPin, Calendar, Clock } from 'lucide-react'

function Location() {
  return (
    <section id="location" className="py-20 px-4 bg-white">
      <div className="max-w-7xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.6 }}
          className="text-center mb-12"
        >
          <h2 className="text-4xl md:text-5xl font-bold text-gray-900 mb-4">
            Event Location
          </h2>
          <p className="text-xl text-gray-600">
            Join us at our office for an immersive AI experience
          </p>
        </motion.div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.6 }}
            className="space-y-6"
          >
            <div className="bg-gradient-to-br from-blue-50 to-indigo-50 rounded-xl p-6 shadow-lg">
              <div className="flex items-start gap-4 mb-4">
                <MapPin className="w-6 h-6 text-blue-600 mt-1" />
                <div>
                  <h3 className="text-xl font-bold text-gray-900 mb-2">
                    Address
                  </h3>
                  <p className="text-gray-700">
                    AppDirect India
                    <br />
                    Pune, Maharashtra, India
                  </p>
                </div>
              </div>
            </div>

            <div className="bg-gradient-to-br from-purple-50 to-pink-50 rounded-xl p-6 shadow-lg">
              <div className="flex items-start gap-4 mb-4">
                <Calendar className="w-6 h-6 text-purple-600 mt-1" />
                <div>
                  <h3 className="text-xl font-bold text-gray-900 mb-2">Date</h3>
                  <p className="text-gray-700">To be announced</p>
                </div>
              </div>
            </div>

            <div className="bg-gradient-to-br from-indigo-50 to-blue-50 rounded-xl p-6 shadow-lg">
              <div className="flex items-start gap-4 mb-4">
                <Clock className="w-6 h-6 text-indigo-600 mt-1" />
                <div>
                  <h3 className="text-xl font-bold text-gray-900 mb-2">Time</h3>
                  <p className="text-gray-700">Full day event</p>
                </div>
              </div>
            </div>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, x: 20 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.6 }}
            className="rounded-xl overflow-hidden shadow-lg"
          >
            <iframe
              src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3930.310034205593!2d73.92600257972352!3d18.515585566381823!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x3bc2c18cf4eaad8d%3A0xc5835f1d9e3a91d3!2sAppDirect%20India!5e0!3m2!1sen!2sin!4v1762854087901!5m2!1sen!2sin"
              width="100%"
              height="100%"
              style={{ border: 0, minHeight: '450px' }}
              allowFullScreen=""
              loading="lazy"
              referrerPolicy="no-referrer-when-downgrade"
              className="w-full"
            ></iframe>
          </motion.div>
        </div>
      </div>
    </section>
  )
}

export default Location

