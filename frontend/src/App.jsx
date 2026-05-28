const stats = [
  { value: '12+', label: 'Nam kinh nghiem' },
  { value: '3.8K', label: 'Hoi vien dang tap' },
  { value: '48', label: 'PT chuyen nghiep' },
  { value: '24/7', label: 'Mo cua moi ngay' },
]

const programs = [
  {
    title: 'Strength Floor',
    meta: 'Powerlifting · Compound',
    desc: 'Khu tang suc manh voi rack, deadlift platform va chuong trinh theo tien do ro rang.',
  },
  {
    title: 'Forge Conditioning',
    meta: 'HIIT · Athletic',
    desc: 'Buoi tap cuong do cao, toi uu suc ben va kha nang dot nang luong trong 45 phut.',
  },
  {
    title: 'Personal Coaching',
    meta: '1:1 · Nutrition',
    desc: 'PT theo sat muc tieu, chi so, lich tap va thoi quen dinh duong cua tung hoi vien.',
  },
]

const schedule = [
  ['05:30', 'Morning Forge', 'Strength'],
  ['12:00', 'Lunch Burn', 'Conditioning'],
  ['18:30', 'Iron Class', 'Power'],
  ['21:00', 'Late Grind', 'Open gym'],
]

function App() {
  return (
    <main>
      <section className="hero grid-bg" id="home">
        <div className="hero-accent" aria-hidden="true" />
        <header className="site-header">
          <a className="brand" href="#home" aria-label="Iron Forge Gym">
            IRON FORGE
          </a>
          <nav className="nav-links" aria-label="Dieu huong chinh">
            <a href="#programs">Lich tap</a>
            <a href="#coaching">Huan luyen</a>
            <a href="#pricing">Goi tap</a>
            <a href="#contact">Lien he</a>
          </nav>
          <a className="nav-cta" href="#contact">
            Thu ngay
          </a>
        </header>

        <div className="hero-shell container">
          <div className="hero-copy">
            <p className="section-eyebrow fade-up">Ha Noi · Tap luyen chuyen nghiep</p>
            <h1 className="hero-title fade-up delay-1">
              REN
              <span>THEP</span>
              THAN
            </h1>
            <p className="hero-sub fade-up delay-2">
              Khong gian luyen tap danh cho nhung nguoi muon thuc su thay doi.
              Khong pho truong. Chi la ky luat, tien do va ket qua.
            </p>
            <div className="hero-actions fade-up delay-3">
              <a className="btn-primary" href="#pricing">
                Dang ky ngay
              </a>
              <a className="btn-outline" href="#programs">
                Xem chuong trinh
              </a>
            </div>
          </div>

          <div className="hero-visual fade-up delay-2" aria-label="Khu tap Iron Forge">
            <div className="visual-frame">
              <div className="plate plate-lg" />
              <div className="barbell" />
              <div className="rack rack-left" />
              <div className="rack rack-right" />
              <div className="floor-lines" />
              <div className="visual-label">
                <span>Forge floor</span>
                <strong>820 m2</strong>
              </div>
            </div>
          </div>
        </div>

        <div className="stats-strip">
          {stats.map((item) => (
            <div className="stat-card" key={item.label}>
              <strong>{item.value}</strong>
              <span>{item.label}</span>
            </div>
          ))}
        </div>
      </section>

      <section className="section programs" id="programs">
        <div className="container">
          <div className="section-head">
            <p className="section-eyebrow">Chu luc tap luyen</p>
            <h2>Moi buoi tap co muc tieu.</h2>
          </div>
          <div className="program-grid">
            {programs.map((program) => (
              <article className="program-card" key={program.title}>
                <p>{program.meta}</p>
                <h3>{program.title}</h3>
                <span>{program.desc}</span>
              </article>
            ))}
          </div>
        </div>
      </section>

      <section className="section split-section" id="coaching">
        <div className="container split-grid">
          <div>
            <p className="section-eyebrow">Coaching</p>
            <h2>PT theo du lieu, khong tap theo cam tinh.</h2>
          </div>
          <div className="coaching-copy">
            <p>
              Moi hoi vien co lich tap, muc tieu va chi so theo doi rieng. Doi ngu
              huan luyen vien lap ke hoach theo suc manh nen tang, lich sinh hoat va
              thoi gian co the duy tri.
            </p>
            <a className="btn-primary" href="#contact">
              Gap PT
            </a>
          </div>
        </div>
      </section>

      <section className="section schedule-section" id="pricing">
        <div className="container schedule-grid">
          <div className="schedule-panel">
            <p className="section-eyebrow">Lich hom nay</p>
            {schedule.map(([time, name, type]) => (
              <div className="schedule-row" key={time}>
                <strong>{time}</strong>
                <span>{name}</span>
                <em>{type}</em>
              </div>
            ))}
          </div>
          <div className="price-panel">
            <span>Membership</span>
            <h2>1.290K</h2>
            <p>Truy cap full gym, lop nhom va tracking co ban trong 30 ngay.</p>
            <a className="btn-primary" href="#contact">
              Bat dau
            </a>
          </div>
        </div>
      </section>

      <footer className="site-footer" id="contact">
        <div className="container footer-grid">
          <div>
            <strong>IRON FORGE</strong>
            <p>Mo cua 24/7 · Hotline 0900 000 000 · Ha Noi</p>
          </div>
          <a className="btn-outline" href="mailto:hello@ironforge.test">
            Dat lich tu van
          </a>
        </div>
      </footer>
    </main>
  )
}

export default App
