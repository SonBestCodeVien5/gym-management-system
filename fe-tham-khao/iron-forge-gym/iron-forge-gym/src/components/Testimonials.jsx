import styles from './Testimonials.module.css'

const TESTIMONIALS = [
  {
    name: 'Nguyễn Văn Hùng',
    role: 'Kỹ sư phần mềm',
    result: '−18kg trong 5 tháng',
    text: 'Trước đây tôi đã thử nhiều gym nhưng không kiên trì được. Iron Forge thay đổi hoàn toàn — PT ở đây không chỉ dạy kỹ thuật mà còn giúp tôi xây dựng thói quen thực sự.',
    rating: 5,
    duration: '8 tháng',
  },
  {
    name: 'Trần Thị Mai Anh',
    role: 'Nhân viên văn phòng',
    result: 'Tăng 6kg cơ trong 4 tháng',
    text: 'Môi trường tập luyện cực kỳ chuyên nghiệp và thân thiện. Lớp HIIT buổi sáng giúp tôi tỉnh táo cả ngày làm việc. Sẽ không đổi gym khác!',
    rating: 5,
    duration: '1 năm',
  },
  {
    name: 'Phạm Quốc Bảo',
    role: 'Doanh nhân',
    result: 'Chạy half marathon đầu tiên',
    text: 'Gói Elite xứng đáng từng đồng. PT riêng theo dõi sát từng buổi, chế độ dinh dưỡng được lên kế hoạch chi tiết. Sau 6 tháng tôi hoàn thành half marathon đầu đời.',
    rating: 5,
    duration: '6 tháng',
  },
]

export default function Testimonials() {
  return (
    <section id="testimonials" className={`section ${styles.section}`}>
      <div className="container">
        <p className="section-eyebrow">Đánh Giá</p>
        <h2 className={styles.title}>
          HỌ ĐÃ<br />
          <span className={styles.stroke}>THAY ĐỔI</span>
        </h2>

        <div className={styles.grid}>
          {TESTIMONIALS.map((t, i) => (
            <div key={i} className={styles.card}>
              <div className={styles.stars}>
                {'★'.repeat(t.rating)}
              </div>
              <p className={styles.text}>"{t.text}"</p>

              <div className={styles.resultTag}>{t.result}</div>

              <div className={styles.footer}>
                <div className={styles.avatar}>{t.name[0]}</div>
                <div>
                  <div className={styles.name}>{t.name}</div>
                  <div className={styles.meta}>{t.role} · {t.duration}</div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
