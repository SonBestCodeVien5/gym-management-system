import Navbar     from './components/Navbar'
import Hero       from './components/Hero'
import Stats      from './components/Stats'
import Services   from './components/Services'
import Trainers   from './components/Trainers'
import Pricing    from './components/Pricing'
import Testimonials from './components/Testimonials'
import Footer     from './components/Footer'

export default function App() {
  return (
    <>
      <Navbar />
      <main>
        <Hero />
        <Stats />
        <Services />
        <Trainers />
        <Pricing />
        <Testimonials />
      </main>
      <Footer />
    </>
  )
}
