# PostalCode API Server

Gin 기반 한국 우편번호 및 주소 조회 REST API 서버입니다.

## 🚀 빠른 시작

### 1. 빌드

```bash
cd cmd/postalcode-api
go build -o postalcode-api
```

### 2. 환경 설정

프로젝트 루트에 `.env` 파일 생성:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=postalcode

# Server Configuration (optional)
SERVER_PORT=8080
```

### 3. 실행

```bash
# .env 파일 사용
./postalcode-api

# 또는 플래그로 직접 설정
./postalcode-api -dsn "user:pass@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True"

# 포트 변경
./postalcode-api -port 9000

# 호스트 변경
./postalcode-api -host 127.0.0.1 -port 8080
```

### 4. 옵션

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | `8080` | 서버 포트 |
| `-host` | `0.0.0.0` | 서버 호스트 |
| `-dsn` | `""` | 데이터베이스 DSN (.env 우선순위 오버라이드) |
| `-env` | `"."` | .env 파일이 있는 디렉토리 경로 |

## 📡 API 엔드포인트

### 도로명주소 API

#### 1. 우편번호로 정확 조회
```bash
GET /api/v1/postal-codes/zipcode/{code}

# Example
curl http://localhost:8080/api/v1/postal-codes/zipcode/01000
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "zip_code": "01000",
      "sido_name": "서울특별시",
      "sigungu_name": "강북구",
      "road_name": "삼양로177길",
      ...
    }
  ],
  "total": 3
}
```

#### 2. 우편번호 앞 3자리로 빠른 검색 (권장)
```bash
GET /api/v1/postal-codes/prefix/{prefix}?page=1&limit=10

# Example
curl http://localhost:8080/api/v1/postal-codes/prefix/010?limit=20
```

#### 3. 복합 조건 검색
```bash
GET /api/v1/postal-codes/search?sido_name={시도}&sigungu_name={시군구}&road_name={도로명}&page=1&limit=10

# Example
curl 'http://localhost:8080/api/v1/postal-codes/search?sido_name=서울&sigungu_name=강북&limit=10'
```

**Query Parameters:**
- `zip_code`: 우편번호 (5자리, 정확 매칭)
- `zip_prefix`: 우편번호 앞 3자리 (빠른 검색)
- `sido_name`: 시도명 (부분 매칭)
- `sigungu_name`: 시군구명 (부분 매칭)
- `road_name`: 도로명 (부분 매칭)
- `page`: 페이지 번호 (기본 1)
- `limit`: 페이지당 결과 개수 (기본 10, 최대 100)

### 지번주소 API

#### 1. 우편번호로 지번주소 조회
```bash
GET /api/v1/postal-codes/land/zipcode/{code}

# Example
curl http://localhost:8080/api/v1/postal-codes/land/zipcode/25627
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "zip_code": "25627",
      "sido_name": "강원특별자치도",
      "sigungu_name": "강릉시",
      "eupmyeondong_name": "강동면",
      "ri_name": "모전리",
      "is_mountain": false,
      ...
    }
  ],
  "total": 2
}
```

#### 2. 우편번호 앞 3자리로 지번주소 빠른 검색
```bash
GET /api/v1/postal-codes/land/prefix/{prefix}?page=1&limit=10

# Example
curl http://localhost:8080/api/v1/postal-codes/land/prefix/256?limit=20
```

#### 3. 지번주소 복합 조건 검색
```bash
GET /api/v1/postal-codes/land/search?sido_name={시도}&eupmyeondong_name={읍면동}&ri_name={리명}

# Example
curl 'http://localhost:8080/api/v1/postal-codes/land/search?sido_name=강원&eupmyeondong_name=강동면&ri_name=모전리'
```

**Query Parameters:**
- `zip_code`: 우편번호 (5자리, 정확 매칭)
- `zip_prefix`: 우편번호 앞 3자리 (빠른 검색)
- `sido_name`: 시도명 (부분 매칭)
- `sigungu_name`: 시군구명 (부분 매칭)
- `eupmyeondong_name`: 읍면동명 (부분 매칭)
- `ri_name`: 리명 (부분 매칭)
- `page`: 페이지 번호 (기본 1)
- `limit`: 페이지당 결과 개수 (기본 10, 최대 100)

### 헬스 체크

```bash
GET /health

# Example
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "ok",
  "service": "korean-postalcode",
  "version": "1.0.0"
}
```

## 🔧 설정

### 환경변수 (.env)

```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=postalcode

# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
```

### 데이터베이스 마이그레이션

서버 시작 시 자동으로 테이블이 생성됩니다 (AutoMigrate).

수동으로 마이그레이션하려면:

```bash
# 도로명주소 테이블
mysql -u root -p postalcode < ../../migrations/create_postal_code_roads.sql

# 지번주소 테이블
mysql -u root -p postalcode < ../../migrations/create_postal_code_lands.sql
```

## 📦 데이터 Import

### 1. 데이터 다운로드

우체국에서 최신 데이터를 다운로드합니다:
- [우체국 우편번호 서비스](https://www.epost.go.kr/search/zipcode/areacdAddressDown.jsp)
- **"범위주소 DB"** 다운로드 후 압축해제

### 2. 데이터 Import

서버 실행 전에 데이터를 import하세요:

```bash
# 도로명주소 데이터
cd ../postalcode-import
./postalcode-import \
    -file "../../data/road_address.txt" \
    -type road \
    -batch 1000

# 지번주소 데이터
./postalcode-import \
    -dsn "user:pass@tcp(localhost:3306)/dbname" \
    -file "../../data/land_address.txt" \
    -type land \
    -batch 1000
```

⚠️ **주의**: Import 시 기존 데이터가 자동으로 TRUNCATE되고 새 데이터로 대체됩니다.

## 🐳 Docker (선택사항)

### Dockerfile 예시

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o postalcode-api ./cmd/postalcode-api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/postalcode-api .
COPY .env .

EXPOSE 8080
CMD ["./postalcode-api"]
```

### 빌드 및 실행

```bash
docker build -t postalcode-api .
docker run -p 8080:8080 --env-file .env postalcode-api
```

## 🛡️ 보안

- CORS가 기본적으로 활성화되어 있습니다 (`*` 허용)
- 프로덕션 환경에서는 CORS 설정을 수정하세요
- Database credentials는 환경변수로 관리하세요
- HTTPS를 사용하는 것을 권장합니다

## ⚡ 성능

- 우편번호 prefix 검색은 인덱스 최적화로 3-5배 빠릅니다
- Limit/Offset 페이징 지원
- Connection pooling 자동 설정 (GORM)

## 📝 라이센스

MIT License
