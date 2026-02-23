## 🗺️ Learning Roadmap & Progress

這是一個為了熟悉 Go Backend 生態系與 Mercari 技術堆疊而建立的實戰專案。

### Phase 1: HTTP Basics & Framework
- [x] **1.1 Web Server**: Setup Gin framework and implement `GET /ping`
- [x] **1.2 Data Struct**: Define `Item` struct with JSON tags
- [x] **1.3 RESTful API**: Implement basic CRUD (Create, Read) for items (In-memory)

### Phase 2: Database & Docker
- [x] **2.1 Docker**: Containerize the app using `Dockerfile`
- [x] **2.2 Database**: Setup PostgreSQL using `docker-compose`
- [x] **2.3 Persistence**: Connect Go with DB using GORM/sqlx and migrate schema

### Phase 3: Architecture & Quality
- [x] **3.1 Layered Arch**: Refactor code into Router -> Controller -> Service -> Repository
- [x] **3.2 Testing**: Add Unit Tests for Service layer (using `testify` or `gomock`)
- [x] **3.3 Configuration**: Manage env variables (Viper or Godotenv)

### Phase 4: Advanced Features (Mercari-like)
- [x] **4.1 Auth**: Implement Signup/Login with JWT
- [x] **4.2 Transaction**: Handle concurrent purchases (Database Transaction)
- [x] **4.3 Logging**: Implement structured logging

### Phase 5: CI/CD & Real-World Infrastructure
- [x] **5.1 CI Pipeline**: Setup GitHub Actions to automate `go test` and code linting on every push
- [x] **5.2 Image Upload**: Handle file uploads (multipart/form-data) for item pictures and store them
- [x] **5.3 Pagination**: Implement `limit` and `offset` query parameters for `GET /items` to handle large databases

### Phase 6: Full-Stack Integration & Optimization
- [ ] **6.1 API Docs**: Auto-generate API documentation using Swagger (swaggo) for frontend collaboration
- [ ] **6.2 Web Frontend**: Build a simple TypeScript/React UI to fetch and display items from the Go API
- [ ] **6.3 Caching (Bonus)**: Introduce Redis to cache hot items and reduce PostgreSQL load
---