import styles from './Footer.module.css'

const LINKS = {
  'Chương Trình': ['Strength Training', 'HIIT & Cardio', 'Combat Sports', 'Yoga & Recovery', 'Personal Training'],
  'Thông Tin':    ['Về Chúng Tôi', 'Đội Ngũ PT', 'Cơ Sở Vật Chất', 'Blog Sức Khoẻ', 'Tuyển Dụng'],
  'Hỗ Trợ':      ['Câu Hỏi Thường Gặp', 'Chính Sách Hoàn Tiền', 'Điều Khoản Sử Dụng', 'Bảo Mật', 'Liên Hệ'],
}

export default function Footer() {
  return (
    <footer id="footer" className={styles.footer}>
      <div className={`container ${styles.top}`}>
        {/* Brand */}
        <div className={styles.brand}>
          <div className={styles.logo}>
            IRON <span className={styles.logoAccent}>FORGE</span>
          </div>
          <p className={styles.tagline}>
            Rèn thân. Rèn tâm.<br />Không giới hạn.
          </p>
          <div className={styles.contact}>
            <p>📍 123 Láng Hạ, Đống Đa, Hà Nội</p>
            <p>📞 0912 345 678</p>
            <p>✉️ hello@ironforge.vn</p>
          </div>
          <div className={styles.socials}>
            {['FB', 'IG', 'YT', 'TK'].map(s => (
              <a key={s} href="#" className={styles.social}>{s}</a>
            ))}
          </div>
        </div>

        {/* Links */}
        {Object.entries(LINKS).map(([heading, items]) => (
          <div key={heading} className={styles.col}>
            <h4 className={styles.colHeading}>{heading}</h4>
            <ul className={styles.linkList}>
              {items.map(item => (
                <li key={item}>
                  <a href="#" className={styles.link}>{item}</a>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </div>

      {/* Bottom bar */}
      <div className={styles.bottom}>
        <div className="container">
          <div className={styles.bottomInner}>
            <p className={styles.copy}>© 2025 Iron Forge Gym. All rights reserved.</p>
            <p className={styles.made}>Thiết kế với 🔥 tại Hà Nội</p>
          </div>
        </div>
      </div>
    </footer>
  )
}
