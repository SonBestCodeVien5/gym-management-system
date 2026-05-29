import styles from './Hero.module.css'

export default function Hero() {
  return (
    <section className={`${styles.hero} grid-bg`}>
      {/* Decorative accent bar */}
      <div className={styles.accentBar} />

      <div className={`container ${styles.inner}`}>
        {/* Left content */}
        <div className={styles.content}>
          <p className={`section-eyebrow fade-up ${styles.eyebrow}`}>
            Hà Nội · Tập Luyện Chuyên Nghiệp
          </p>

          <h1 className={`${styles.headline} fade-up delay-1`}>
            RÈN<br />
            <span className={styles.stroke}>THÉP</span><br />
            <span className={styles.thirdLine}>THÂN</span>
          </h1>

          <p className={`${styles.sub} fade-up delay-2`}>
            Không gian luyện tập dành cho những người muốn thực sự thay đổi.
            Không phô trương. Chỉ là kết quả.
          </p>

          <div className={`${styles.ctaRow} fade-up delay-3`}>
            <a href="#pricing" className="btn-primary">
              Đăng Ký Ngay →
            </a>
            <a href="#services" className="btn-outline">
              Xem Chương Trình
            </a>
          </div>

          {/* Trust badges */}
          <div className={`${styles.badges} fade-up delay-4`}>
            <span className={styles.badge}>✓ Không cần hợp đồng</span>
            <span className={styles.badge}>✓ Buổi tập thử miễn phí</span>
            <span className={styles.badge}>✓ PT chuyên nghiệp</span>
          </div>
        </div>

        {/* Right – visual block */}
        <div className={styles.visual}>
          <div className={styles.imgBlock}>
            {/* Replace src with your actual gym photo */}
            <div className={styles.imgPlaceholder}>
              <span className={styles.imgLabel}>GYM PHOTO</span>
              <p className={styles.imgHint}>Thay bằng <code>{'<img src="..." />'}</code></p>
            </div>

            {/* Floating tag */}
            <div className={styles.floatTag}>
              <span className={styles.floatNum}>24/7</span>
              <span className={styles.floatText}>Luôn Mở Cửa</span>
            </div>
          </div>

          {/* Decorative lines */}
          <div className={styles.decoLines}>
            {[...Array(8)].map((_, i) => (
              <div key={i} className={styles.decoLine} style={{ opacity: 0.04 + i * 0.015 }} />
            ))}
          </div>
        </div>
      </div>

      {/* Bottom gradient fade */}
      <div className={styles.bottomFade} />
    </section>
  )
}
