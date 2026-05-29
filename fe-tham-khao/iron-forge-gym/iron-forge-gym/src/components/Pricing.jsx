import { useState } from 'react'
import styles from './Pricing.module.css'

const PLANS = [
  {
    id: 'basic',
    name: 'BASIC',
    tagline: 'Bắt đầu hành trình',
    monthly: 499,
    quarterly: 420,
    features: [
      'Tập tự do không giới hạn',
      'Khu cardio & tạ tự do',
      'Tủ đồ & phòng tắm',
      'Ứng dụng theo dõi tập luyện',
      '1 buổi tư vấn PT',
    ],
    notIncluded: ['Lớp học nhóm', 'Phòng yoga', 'Dinh dưỡng tư vấn'],
    cta: 'Chọn Basic',
    highlight: false,
  },
  {
    id: 'pro',
    name: 'PRO',
    tagline: 'Được chọn nhiều nhất',
    monthly: 899,
    quarterly: 760,
    features: [
      'Tất cả tính năng Basic',
      'Lớp học nhóm không giới hạn',
      'Phòng yoga & recovery',
      '4 buổi PT / tháng',
      'Tư vấn dinh dưỡng cơ bản',
      'Đánh giá thể lực hàng tháng',
    ],
    notIncluded: ['PT không giới hạn'],
    cta: 'Chọn Pro',
    highlight: true,
  },
  {
    id: 'elite',
    name: 'ELITE',
    tagline: 'Toàn diện không giới hạn',
    monthly: 1599,
    quarterly: 1350,
    features: [
      'Tất cả tính năng Pro',
      'PT không giới hạn',
      'Lập kế hoạch dinh dưỡng chuyên sâu',
      'Phục hồi & massage 2x/tháng',
      'Guest pass 2 người/tháng',
      'Ưu tiên đặt lịch lớp nhóm',
      'Hỗ trợ 24/7 qua app',
    ],
    notIncluded: [],
    cta: 'Chọn Elite',
    highlight: false,
  },
]

export default function Pricing() {
  const [billing, setBilling] = useState('monthly') // 'monthly' | 'quarterly'

  return (
    <section id="pricing" className={`section ${styles.pricing}`}>
      <div className="container">
        <p className="section-eyebrow">Gói Tập</p>
        <div className={styles.topRow}>
          <h2 className={styles.title}>
            ĐẦU TƯ VÀO<br />
            <span className={styles.stroke}>BẢN THÂN</span>
          </h2>

          {/* Billing toggle */}
          <div className={styles.toggle}>
            <button
              className={`${styles.toggleBtn} ${billing === 'monthly' ? styles.toggleActive : ''}`}
              onClick={() => setBilling('monthly')}
            >Hàng tháng</button>
            <button
              className={`${styles.toggleBtn} ${billing === 'quarterly' ? styles.toggleActive : ''}`}
              onClick={() => setBilling('quarterly')}
            >
              3 tháng
              <span className={styles.saveBadge}>–15%</span>
            </button>
          </div>
        </div>

        <div className={styles.grid}>
          {PLANS.map(plan => {
            const price = billing === 'monthly' ? plan.monthly : plan.quarterly
            return (
              <div key={plan.id} className={`${styles.card} ${plan.highlight ? styles.cardHighlight : ''}`}>
                {plan.highlight && (
                  <div className={styles.popularBadge}>PHỔ BIẾN NHẤT</div>
                )}
                <div className={styles.planName}>{plan.name}</div>
                <div className={styles.planTagline}>{plan.tagline}</div>

                <div className={styles.priceRow}>
                  <span className={styles.currency}>₫</span>
                  <span className={styles.price}>{price.toLocaleString('vi')}</span>
                  <span className={styles.per}>/tháng</span>
                </div>

                {billing === 'quarterly' && (
                  <p className={styles.billedAs}>
                    Thanh toán ₫{(price * 3).toLocaleString('vi')} / 3 tháng
                  </p>
                )}

                <div className={styles.divider} />

                <ul className={styles.featureList}>
                  {plan.features.map(f => (
                    <li key={f} className={styles.featureItem}>
                      <span className={styles.checkMark}>✓</span> {f}
                    </li>
                  ))}
                  {plan.notIncluded.map(f => (
                    <li key={f} className={`${styles.featureItem} ${styles.featureDim}`}>
                      <span className={styles.crossMark}>✗</span> {f}
                    </li>
                  ))}
                </ul>

                <a
                  href="#"
                  className={plan.highlight ? 'btn-primary' : 'btn-outline'}
                  style={{ display: 'block', textAlign: 'center', marginTop: 'auto' }}
                >
                  {plan.cta}
                </a>
              </div>
            )
          })}
        </div>

        <p className={styles.note}>
          * Tất cả giá chưa bao gồm VAT · Thử miễn phí 3 ngày · Không cần cam kết hợp đồng
        </p>
      </div>
    </section>
  )
}
