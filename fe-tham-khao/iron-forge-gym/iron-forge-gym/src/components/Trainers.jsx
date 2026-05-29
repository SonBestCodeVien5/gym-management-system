import styles from './Trainers.module.css'

const TRAINERS = [
  {
    name: 'Minh Tuấn',
    role: 'Head Coach · Powerlifting',
    exp: '10 năm',
    certs: ['NSCA-CPT', 'CSCS'],
    specialty: 'Squat · Deadlift · Bench',
    initials: 'MT',
  },
  {
    name: 'Lan Anh',
    role: 'PT · Functional Training',
    exp: '7 năm',
    certs: ['ACE-CPT', 'FMS'],
    specialty: 'Weight Loss · Mobility',
    initials: 'LA',
  },
  {
    name: 'Quang Hải',
    role: 'Combat Coach · MMA',
    exp: '12 năm',
    certs: ['IBJJF', 'WBC Boxing'],
    specialty: 'Boxing · Muay Thai · BJJ',
    initials: 'QH',
  },
  {
    name: 'Thu Hà',
    role: 'Yoga & Recovery',
    exp: '8 năm',
    certs: ['RYT-500', 'Pilates Mat'],
    specialty: 'Hatha · Yin · Breath Work',
    initials: 'TH',
  },
]

export default function Trainers() {
  return (
    <section id="trainers" className={`section ${styles.trainers}`}>
      <div className="container">
        <p className="section-eyebrow">Đội Ngũ</p>
        <h2 className={styles.title}>
          HUẤN LUYỆN VIÊN<br />
          <span className={styles.stroke}>CHUYÊN NGHIỆP</span>
        </h2>

        <div className={styles.grid}>
          {TRAINERS.map((t, i) => (
            <div key={i} className={styles.card}>
              {/* Photo placeholder */}
              <div className={styles.photo}>
                <span className={styles.initials}>{t.initials}</span>
                <div className={styles.photoHint}>Thêm ảnh</div>
              </div>

              <div className={styles.info}>
                <div className={styles.nameRow}>
                  <h3 className={styles.name}>{t.name}</h3>
                  <span className={styles.exp}>{t.exp}</span>
                </div>
                <p className={styles.role}>{t.role}</p>
                <p className={styles.specialty}>{t.specialty}</p>
                <div className={styles.certRow}>
                  {t.certs.map(c => (
                    <span key={c} className={styles.cert}>{c}</span>
                  ))}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
