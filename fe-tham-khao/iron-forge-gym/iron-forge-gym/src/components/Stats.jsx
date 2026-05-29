import styles from './Stats.module.css'

const STATS = [
  { num: '12+',  label: 'Năm Kinh Nghiệm'  },
  { num: '3.8K', label: 'Hội Viên Hiện Tại' },
  { num: '48',   label: 'PT Chuyên Nghiệp'  },
  { num: '99%',  label: 'Hài Lòng'          },
  { num: '24/7', label: 'Mở Cửa'            },
]

export default function Stats() {
  return (
    <div className={styles.statsBar}>
      <div className={`container ${styles.grid}`}>
        {STATS.map((s, i) => (
          <div key={i} className={styles.stat}>
            <span className={styles.num}>{s.num}</span>
            <span className={styles.label}>{s.label}</span>
          </div>
        ))}
      </div>
    </div>
  )
}
