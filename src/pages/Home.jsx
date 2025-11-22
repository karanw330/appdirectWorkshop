import Hero from '../components/Hero'
import SessionsSpeakers from '../components/SessionsSpeakers'
import RegistrationForm from '../components/RegistrationForm'
import Location from '../components/Location'
import Footer from '../components/Footer'

function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
      <Hero />
      <SessionsSpeakers />
      <RegistrationForm />
      <Location />
      <Footer />
    </div>
  )
}

export default Home

