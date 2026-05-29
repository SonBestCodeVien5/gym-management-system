import { useState, useEffect } from 'react'
import styles from './Navbar.module.css'

const NAV_LINKS = [
  { label: 'Lịch Tập',   href: '#services'  },
  { label: 'Huấn Luyện', href: '#trainers'  },
  { label: 'Gói Tập',    href: '#pricing'   },
  { label: 'Đánh Giá',   href: '#testimonials' },
  { label: 'Liên Hệ',   href: '#footer'    },
]

export default function Navbar() {
  const [scrolled,  setScrolled]  = useState(false)
  const [menuOpen,  setMenuOpen]  = useState(false)

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 40)
    window.addEventListener('scroll', onScroll)
    return () => window.removeEventListener('scroll', onScroll)
  }, [])

  return (
    <header className={`${styles.navbar} ${scrolled ? styles.scrolled : ''}`}>
      <div className={`container ${styles.inner}`}>
        {/* Logo */}
        <a href="#" className={styles.logo}>
          IRON<span className={styles.logoAccent}> FORGE</span>
        </a>

        {/* Desktop nav */}
        <nav className={styles.nav}>
          {NAV_LINKS.map(link => (
            <a key={link.href} href={link.href} className={styles.navLink}>
              {link.label}
            </a>
          ))}
        </nav>

        {/* CTA */}
        <a href="#pricing" className={`btn-primary ${styles.ctaBtn}`}>
          Đăng Ký Ngay
        </a>

        {/* Hamburger */}
        <button
          className={`${styles.hamburger} ${menuOpen ? styles.open : ''}`}
          onClick={() => setMenuOpen(v => !v)}
          aria-label="Toggle menu"
        >
          <span /><span /><span />
        </button>
      </div>

      {/* Mobile menu */}
      <div className={`${styles.mobileMenu} ${menuOpen ? styles.mobileMenuOpen : ''}`}>
        {NAV_LINKS.map(link => (
          <a key={link.href} href={link.href} className={styles.mobileLink}
             onClick={() => setMenuOpen(false)}>
            {link.label}
          </a>
        ))}
        <a href="#pricing" className={`btn-primary ${styles.mobileCta}`}
           onClick={() => setMenuOpen(false)}>
          Đăng Ký Ngay
        </a>
      </div>
    </header>
  )
}
