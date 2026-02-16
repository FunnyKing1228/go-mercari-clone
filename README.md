## 🗺️ Learning Roadmap & Progress

這是一個為了熟悉 Go Backend 生態系與 Mercari 技術堆疊而建立的實戰專案。

### Phase 1: HTTP Basics & Framework
- [x] **1.1 Web Server**: Setup Gin framework and implement `GET /ping`
- [x] **1.2 Data Struct**: Define `Item` struct with JSON tags
- [x] **1.3 RESTful API**: Implement basic CRUD (Create, Read) for items (In-memory)

### Phase 2: Database & Docker
- [x] **2.1 Docker**: Containerize the app using `Dockerfile`
- [ ] **2.2 Database**: Setup PostgreSQL using `docker-compose`
- [ ] **2.3 Persistence**: Connect Go with DB using GORM/sqlx and migrate schema

### Phase 3: Architecture & Quality
- [ ] **3.1 Layered Arch**: Refactor code into Router -> Controller -> Service -> Repository
- [ ] **3.2 Testing**: Add Unit Tests for Service layer (using `testify` or `gomock`)
- [ ] **3.3 Configuration**: Manage env variables (Viper or Godotenv)

### Phase 4: Advanced Features (Mercari-like)
- [ ] **4.1 Auth**: Implement Signup/Login with JWT
- [ ] **4.2 Transaction**: Handle concurrent purchases (Database Transaction)
- [ ] **4.3 Logging**: Implement structured logging

---