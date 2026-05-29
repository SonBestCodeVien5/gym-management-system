import styles from './Services.module.css'

const SERVICES = [
  {
    id: '01',
    icon: '🏋️',
    title: 'Strength Training',
    desc: 'Khu vực tạ tự do và máy tập chuyên dụng với đầy đủ thiết bị nhập khẩu tiêu chuẩn quốc tế.',
    tags: ['Free Weight', 'Barbell', 'Dumbbell'],
  },
  {
    id: '02',
    icon: '🔥',
    title: 'HIIT & Cardio',
    desc: 'Lớp học nhóm cường độ cao kết hợp cardio và sức mạnh. Đốt cháy calo tối đa trong thời gian ngắn.',
    tags: ['HIIT', 'Cardio', 'Boxing'],
  },
  {
    id: '03',
    icon: '🥊',
    title: 'Combat Sports',
    desc: 'Võ thuật tổng hợp, boxing, muay thai với các huấn luyện viên có chứng chỉ quốc tế.',
    tags: ['Boxing', 'MMA', 'Muay Thai'],
  },
  {
    id: '04',
    icon: '🧘',
    title: 'Recovery & Yoga',
    desc: 'Phòng tập riêng biệt cho yoga, stretching và phục hồi. Cân bằng cơ thể sau những buổi tập nặng.',
    tags: ['Yoga', 'Stretch', 'Meditation'],
  },
  {
    id: '05',
    icon: '🎯',
    title: 'Personal Training',
    desc: 'Chương trình tập luyện thiết kế riêng 1-1 với PT, bám sát mục tiêu và theo dõi tiến trình hàng tuần.',
    tags: ['1-1 PT', 'Custom Plan', 'Nutrition'],
  },
  {
    id: '06',
    icon: '🏃',
    title: 'Endurance & Run',
    desc: 'Treadmill cao cấp, bài tập chạy bộ nâng cao, chuẩn bị cho các giải marathon và triathlon.',
    tags: ['Marathon', 'Treadmill', 'Cycling'],
  },
]

export default function Services() {
  return (
    <section id="services" className={`section ${styles.services}`}>
      <div className="container">
        <p className="section-eyebrow">Chương Trình</p>
        <div className={styles.heading}>
          <h2 className={styles.title}>
            TẤT CẢ NHỮNG GÌ<br />
            <span className={styles.stroke}>BẠN CẦN</span>
          </h2>
          <p className={styles.desc}>
            Từ sức mạnh đến sự dẻo dai — chúng tôi có đầy đủ chương trình
            để bạn đạt được mục tiêu của mình.
          </p>
        </div>

        <div className={styles.grid}>
          {SERVICES.map(s => (
            <div key={s.id} className={styles.card}>
              <div className={styles.cardNum}>{s.id}</div>
              <div className={styles.cardIcon}>{s.icon}</div>
              <h3 className={styles.cardTitle}>{s.title}</h3>
              <p className={styles.cardDesc}>{s.desc}</p>
              <div className={styles.tagRow}>
                {s.tags.map(t => (
                  <span key={t} className={styles.tag}>{t}</span>
                ))}
              </div>
              <div className={styles.cardArrow}>→</div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
