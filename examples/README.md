# PostalCode Examples

실행 가능한 예제 코드 모음입니다.

## 📁 디렉토리 구조

```
examples/
├── basic/                   # 기본 사용법
│   └── main.go             # Repository/Service 기본 사용
├── api/                     # REST API 서버
│   └── main.go             # HTTP API 서버 구현
├── import/                  # 데이터 import
│   └── main.go             # 대량 데이터 import 예제
├── project-integration/     # 기존 프로젝트 통합 (중요!)
│   └── main.go             # 기존 프로젝트 DB/Config 재사용
└── example_test.go          # Go test 예제
```

## 실행 방법

### 1. 환경 설정

먼저 `.env` 파일을 설정합니다:

```bash
cd pkg/postalcode
cp .env.example .env
# .env 파일을 열어 DB 정보 입력
```

또는 환경변수로 직접 설정:

```bash
export POSTALCODE_DB_HOST=localhost
export POSTALCODE_DB_PORT=3306
export POSTALCODE_DB_USER=root
export POSTALCODE_DB_PASSWORD=your_password
export POSTALCODE_DB_NAME=postalcode
```

### 2. 예제 실행

#### Basic Example (기본 사용법)

```bash
cd examples/basic
go run main.go
```

**출력 예시:**
```
=== Example 1: Search by Zip Code ===
Found 5 results for zip code 01000
First result: 서울특별시 강북구 삼양로177길

=== Example 2: Fast Search by Prefix ===
Found 1234 results for prefix 010

=== Example 3: Complex Search ===
Found 10 results (total: 25)
  1. [01000] 서울특별시 강북구 삼양로177길 (건물번호: 93~126)
  ...
```

#### API Server Example (REST API)

```bash
cd examples/api
go run main.go
```

**서버 시작 후:**
```
🚀 PostalCode API Server starting on :8080
📍 API endpoints:
   GET  /api/v1/postal-codes/zipcode/{code}     - Search by zip code
   GET  /api/v1/postal-codes/prefix/{prefix}    - Fast search by prefix
   GET  /api/v1/postal-codes/search             - Complex search
   GET  /health                                  - Health check
```

**테스트:**
```bash
# 다른 터미널에서
curl http://localhost:8080/api/v1/postal-codes/zipcode/01000
curl http://localhost:8080/api/v1/postal-codes/prefix/010
curl "http://localhost:8080/api/v1/postal-codes/search?sido_name=서울&limit=10"
```

#### Import Example (데이터 import)

```bash
cd examples/import
go run main.go
```

**출력 예시:**
```
🔌 Connecting to database...
🔧 Creating table if not exists...
📂 Importing from: ../../docs/addresses/20251028_도로명범위.txt
📦 Batch size: 1000

✅ Progress: 1000/311468 (0.3%)
✅ Progress: 2000/311468 (0.6%)
...

📊 Import Summary:
  ✅ Success: 311468 records
  ❌ Errors:  0 records
  ⏱️  Time:    1m 23s
  📈 Speed:   3750 records/sec

🎉 Import completed successfully!
```

#### Project Integration Example (기존 프로젝트 통합) ⭐

```bash
cd examples/project-integration
go run main.go
```

**핵심 포인트:**
- ✅ 기존 프로젝트의 `config.DatabaseConfig` 재사용
- ✅ 기존 프로젝트의 `datastore.DB` 재사용
- ✅ Gin 프레임워크 통합
- ✅ 설정 파일 하나로 통합 관리

**출력 예시:**
```
📋 Using project configuration
   DB Host: localhost
   DB Name: myproject

✅ Connected to database using project config
✅ PostalCode package initialized with project DB

🚀 Project + PostalCode Integration Server
   Server: http://localhost:8080
```

## 📝 예제 설명

### basic/main.go
- Repository/Service 패턴 사용법
- 3가지 검색 방법 (정확 우편번호, prefix, 복합 검색)
- 결과 출력 및 에러 처리
- **사용 시나리오**: Standalone 애플리케이션

### api/main.go
- HTTP REST API 서버 구현
- 표준 http.ServeMux 사용
- Health check 엔드포인트
- 실제 배포 가능한 코드
- **사용 시나리오**: 독립 API 서버

### import/main.go
- 대량 데이터 import 구현
- Progress callback 사용
- AutoMigrate 자동 테이블 생성
- 성능 측정 (records/sec)
- **사용 시나리오**: 데이터 초기화/마이그레이션

### project-integration/main.go
- **기존 프로젝트 내부에서 사용하는 방법**
- 기존 프로젝트의 config/datastore 재사용
- Gin 프레임워크 통합
- 설정 파일 하나로 통합 관리
- **사용 시나리오**: 기존 프로젝트에 postalcode API 추가

### example_test.go
- Go test 형식 예제
- 단위 테스트 작성 방법
- Mock/Stub 사용 예시
- **사용 시나리오**: 테스트 코드 작성

## 💡 Tips

1. **환경변수 우선순위**: 환경변수 > .env 파일
2. **DB 연결**: 예제는 모두 `LoadConfig()`를 사용하여 설정 로드
3. **에러 처리**: 프로덕션에서는 예제보다 더 철저한 에러 처리 필요
4. **성능**: import는 배치 크기를 조정하여 성능 최적화 가능 (1000~5000)

## 🔗 참고 문서

- [API.md](../docs/API.md) - REST API 엔드포인트 완전 가이드
- [USAGE.md](../docs/USAGE.md) - Repository/Service 사용 가이드
- [INTEGRATION.md](../docs/INTEGRATION.md) - 프로젝트 통합 가이드
- [README.md](../README.md) - 패키지 개요
