# Nhat Ky Phat Trien Du An (Phase 2)

## 1. Muc tieu tai lieu
Tai lieu nay ghi lai:
- Tien do phat trien theo commit.
- Cac van de da gap khi code/chay local.
- Cach da xu ly va bai hoc rut ra.
- Cac kien thuc nen nam de tiep tuc mo rong du an.

Xem them phan hoi-dap kien truc: [faq_why.md](faq_why.md)
Huong dan viet bao cao phan tich-thiet ke: [system_analysis_design_guide.md](system_analysis_design_guide.md)

## 2. Tien do theo timeline (branch main, tuyen tinh)
| Thu tu | Commit | Noi dung chinh |
|---|---|---|
| 1 | `dc276ac` | Khoi tao go module |
| 2 | `99365d1` | Tao project structure, docker MongoDB, model member ban dau |
| 3 | `29009c0` | Chuan hoa line ending LF |
| 4 | `9f76da0` | Bo sung cac model moi |
| 5 | `74c4fd0` | Them ket noi MongoDB + member repository |
| 6 | `aede7ba` | Bo tham so khong dung trong constructor repo |
| 7 | `4b8eb51` | Bo sung tai lieu Go |
| 8 | `6d17db4` | Them truong `IsSuspended` cho member |
| 9 | `bdb00ee` | Them service/handler + registration flow |
| 10 | `95c0f66` | Cap nhat va test lai MongoDB URI |

## 2.1. Tien do moi trong workspace hien tai (sau commit cuoi)
1. `internal/repository/member_repo.go`: them unique index `ccid` trong constructor repository.
2. `internal/repository/course_repo.go`: them repo doc `GetByID` cho `courses`.
3. `internal/repository/branch_repo.go`: them repo doc `GetByID` cho `branches`.
4. `internal/repository/subscription_repo.go`: them repository cho `subscriptions`.
5. `internal/service/subscription_service.go`: them rule create subscription va kiem tra tham chieu member/course/branch.
6. `internal/handlers/subscription_handler.go`: them handler API cho subscription, parse `RFC3339`.
7. `cmd/server/main.go`: wire day du member + subscription flow vao Gin routes.

## 3. Trang thai hien tai cua he thong
Da hoan thanh:
1. Chay duoc backend Go + Gin.
2. Ket noi MongoDB local qua Docker.
3. Route `GET /ping` hoat dong.
4. Route `POST /api/v1/registration` hoat dong.
5. Da co 3 tang cho module member: repository, service, handler.
6. Da co block subscription ban dau: repository, service, handler, wire vao main.
7. Da giam coupling nhe o service bang cach dung `repository.ErrNotFound` thay vi check loi Mongo truc tiep.
8. Da them unique index `ccid` de chong race condition trung member.

Dang co:
1. Member moi co registration va get-by-id.
2. Subscription moi co create va get-by-id.
3. Chua trien khai day du payment/suspension/resume/attendance theo Phase 2.

## 4. Nhat ky van de da gap va cach xu ly

### Van de 1: Authentication failed khi ket noi MongoDB
- Trieu chung: `unable to authenticate` khi `go run`.
- Nguyen nhan: Mat khau trong URI khong khop `docker-compose.yml`.
- Cach xu ly:
  1. Dong bo lai `MONGO_INITDB_ROOT_PASSWORD` va `.env`.
  2. Xac nhan bang `docker exec ... mongosh`.
- Bai hoc: Luon doi chieu 3 diem: `docker-compose.yml`, `.env`, log runtime.
- Kien thuc lien quan (FAQ):
  1. [mongosh interactive vs --eval](faq_why.md#faq-dev-mongosh-modes)
  2. [localhost vs 127.0.0.1](faq_why.md#faq-dev-localhost-vs-loopback)

### Van de 2: Build loi do typo trong handler
- Trieu chung: `undefined: registerMemberRequest`.
- Nguyen nhan: Sai ten struct (`resister...` vs `register...`).
- Cach xu ly: Sua ten struct va gofmt lai file.
- Bai hoc: Loi chinh ta trong ten type rat de xay ra, nen build thu ngay sau moi block code.
- Kien thuc lien quan (FAQ):
  1. [go build vs go run](faq_why.md#faq-dev-build-vs-run)

### Van de 3: Khong hien `Send Request` cho request thu 2 trong file `.http`
- Trieu chung: Chi thay nut gui o block dau.
- Nguyen nhan: File test thieu dau phan tach `###`.
- Cach xu ly:
  1. Tach tung request bang `###`.
  2. Dung extension REST Client.
- Bai hoc: File `.http` can format dung parser cua extension.
- Kien thuc lien quan (FAQ):
  1. [REST Client va dau ###](faq_why.md#faq-dev-rest-client-separator)

### Van de 4: Nham giua `go build` va `go run`
- Trieu chung: Build thanh cong nhung terminal "khong in gi".
- Nguyen nhan: `go build` chi compile, khong chay app.
- Cach xu ly:
  1. Dung `go build ./... && echo "BUILD OK"` de check nhanh.
  2. Dung `go run cmd/server/main.go` de test runtime.
- Bai hoc: Tach ro buoc compile va runtime giup debug nhanh hon.
- Kien thuc lien quan (FAQ):
  1. [go build vs go run](faq_why.md#faq-dev-build-vs-run)

### Van de 5: MongoDB Compass ket noi fail voi `localhost`
- Trieu chung: Cau hinh giong nhau nhung luc ket noi duoc, luc khong.
- Nguyen nhan kha nang cao: Phan giai `localhost` sang IPv6/forwarding khong on dinh trong WSL + Docker.
- Cach xu ly:
  1. Chuyen host sang `127.0.0.1`.
  2. Them `directConnection=true` va `authSource=admin`.
  3. Tat TLS/SSH trong Compass neu khong dung.
- Bai hoc: O local WSL, uu tien `127.0.0.1` cho ket noi on dinh.
- Kien thuc lien quan (FAQ):
  1. [localhost vs 127.0.0.1](faq_why.md#faq-dev-localhost-vs-loopback)

### Van de 6: Kiem tra record MongoDB nhung "khong thay output"
- Trieu chung: Chay `docker exec ... mongosh <uri>` nhung khong in du lieu.
- Nguyen nhan: Day la che do interactive shell, khong co query thi khong in.
- Cach xu ly:
  1. Dung `--eval` de one-shot query.
  2. Hoac vao shell roi go `db.members.find().pretty()`.
- Bai hoc: Phan biet ro interactive mode va one-shot mode.
- Kien thuc lien quan (FAQ):
  1. [mongosh interactive vs --eval](faq_why.md#faq-dev-mongosh-modes)

### Van de 7: Handler subscription parse ngay gio bi sai kieu
- Trieu chung: `parseTimeValue` tra ve sai kieu so voi `time.Time` trong model.
- Nguyen nhan: Handler ban dau dung nham `primitive.DateTime` thay vi `time.Time`.
- Cach xu ly:
  1. Doi helper parse thanh `time.Time`.
  2. Dung `time.Parse(time.RFC3339, ...)` cho input JSON.
- Bai hoc: Khi model da chon `time.Time`, handler phai parse ve cung kieu de khop service va Mongo.

## 5. Kien thuc bo sung can nam

### 5.1 Service vs Repository (vi sao co 2 ham ten gan giong nhau)
- `GetMemberByID` o service la nghiep vu cap ung dung.
- `GetByID` o repository la truy cap du lieu cap DB.
- Service goi repository de tach business logic va storage logic.
- Wiki lien quan: [Service va Repository khac nhau gi](faq_why.md#faq-arch-service-vs-repository)

### 5.2 Coupling va huong giam phu thuoc
Da lam:
1. Mapping loi not-found o repository -> `ErrNotFound` trung lap.
2. Service khong can import package Mongo de check not-found.
3. Handler member khong import repository nua, chi map loi service.

Con de sau (refactor nang):
1. Tach domain model va persistence model.
2. Tach strategy tao ID khoi service.
- Wiki lien quan: [Coupling va cach giam](faq_why.md#faq-arch-coupling)

### 5.4 Block subscription hien tai da co gi
1. `CourseRepository` va `BranchRepository` chi can `GetByID` vi Mongo da tu co index `_id`.
2. `SubscriptionRepository` ho tro `Create` va `GetByID`.
3. `SubscriptionService` kiem tra tham chieu `member/course/branch` va lay gia tu `Course`.
4. `SubscriptionHandler` parse `start_date`, `end_date` theo RFC3339.

### 5.3 Extension MongoDB vs MongoDB Compass
- Extension VS Code: tien cho query nhanh trong luong code.
- Compass: manh hon cho quan tri truc quan (index/schema/explain).
- Cach dung toi uu: dev thuong ngay bang extension, debug/admin sau bang Compass.
- Wiki lien quan: [localhost vs 127.0.0.1](faq_why.md#faq-dev-localhost-vs-loopback)

## 6. Checklist truoc khi tiep tuc buoc lon tiep theo
1. `go build ./...` pass.
2. `GET /ping` pass.
3. `POST /api/v1/registration` pass.
4. Query trong `gym_management.members` thay record vua tao.
5. URI trong `.env` da on dinh voi local (`127.0.0.1`).
6. `POST /api/v1/subscriptions` co the chuan bi test sau khi chot payload mau.

## 7. De xuat buoc tiep theo
1. Hoan thien test API cho `GET /api/v1/members/:id` va `GET /api/v1/subscriptions/:id`.
2. Chot payload mau cho `POST /api/v1/subscriptions`.
3. Bat dau module `payment` / `suspension` theo business rule Phase 2.

---

Cap nhat lan cuoi: theo trang thai code va commit den ngay 2026-04-28.