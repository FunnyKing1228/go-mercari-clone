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
- [x] **6.1 API Docs**: Auto-generate API documentation using Swagger (swaggo) for frontend collaboration
- [x] **6.2 Web Frontend**: Build a simple TypeScript/React UI to fetch and display items from the Go API
- [x] **6.3 Caching (Bonus)**: Introduce Redis to cache hot items and reduce PostgreSQL load

### Phase 7: Core Business Logic & Relational Data
- [x] **7.1 Real User System**: Implement true user registration and secure passwords using `bcrypt` hashing
- [x] **7.2 Relational DB**: Use GORM Preload to establish One-to-Many relationships (e.g., link an `Item` to a specific `User` as the seller)
- [x] **7.3 Advanced Search**: Implement keyword search and category filtering in the `GET /items` API

### Phase 8: Production Readiness & Load Testing
- [ ] **8.1 Frontend Containerization**: Write a Dockerfile for the React app and integrate it into the `docker-compose` network
- [ ] **8.2 Load Testing**: Introduce load testing tools (e.g., `k6` or `JMeter`) to simulate high concurrency and stress test the API
- [ ] **8.3 Performance Tuning**: Analyze bottlenecks from load tests and optimize database queries or Redis caching strategies

---
### 🎓 Final Challenge: Architecture Deep Dive & Interview Prep
*(This is a conceptual validation phase, focusing on explaining the "Why" behind the code for Tier-1 tech interviews.)*
- [ ] **Concurrency Control**: Explain the mechanism of `FOR UPDATE` (Pessimistic Lock) and how it prevents race conditions in e-commerce.
- [ ] **Caching Strategy**: Defend the choice of the Cache-Aside pattern with Redis and explain how to handle cache invalidation and stale data.
- [ ] **Security & Auth**: Articulate the flow of JWT authentication and why it is chosen over stateful session cookies.
---