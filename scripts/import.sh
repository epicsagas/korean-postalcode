#!/bin/bash

# PostalCode 데이터 Import 스크립트
# Usage: ./import.sh -file <file_path> [-type road|land] [-dsn <dsn>] [-batch <size>]
#
# Example:
#   ./import.sh -file data/postal_codes.txt -type road
#   ./import.sh -file data/postal_codes.txt -type land -batch 2000
#   ./import.sh -file data/postal_codes.txt -dsn "user:pass@tcp(localhost:3306)/dbname"

set -e

# 색상 코드
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 기본값
DSN=""
FILE_PATH=""
DATA_TYPE="road"
BATCH_SIZE="1000"

# 플래그 파싱 (공백 구분 및 = 구분 모두 지원)
while [[ $# -gt 0 ]]; do
    case $1 in
        -dsn)
            DSN="$2"
            shift 2
            ;;
        -dsn=*)
            DSN="${1#*=}"
            shift
            ;;
        -file)
            FILE_PATH="$2"
            shift 2
            ;;
        -file=*)
            FILE_PATH="${1#*=}"
            shift
            ;;
        -type)
            DATA_TYPE="$2"
            shift 2
            ;;
        -type=*)
            DATA_TYPE="${1#*=}"
            shift
            ;;
        -batch)
            BATCH_SIZE="$2"
            shift 2
            ;;
        -batch=*)
            BATCH_SIZE="${1#*=}"
            shift
            ;;
        *)
            echo -e "${RED}❌ 알 수 없는 옵션: $1${NC}"
            echo ""
            echo -e "${YELLOW}Usage: $0 -file <file_path> [-type road|land] [-dsn <dsn>] [-batch <size>]${NC}"
            echo ""
            echo -e "${YELLOW}Example:${NC}"
            echo "  $0 -file data/postal_codes.txt -type road"
            echo "  $0 -file=data/postal_codes.txt -type=land -batch=2000"
            echo "  $0 -file data/postal_codes.txt -dsn \"user:pass@tcp(localhost:3306)/dbname\""
            echo ""
            exit 1
            ;;
    esac
done

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}PostalCode Data Import Tool${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# 필수 파라미터 확인
if [ -z "$FILE_PATH" ]; then
    echo -e "${RED}❌ 오류: -file 파라미터는 필수입니다${NC}"
    echo ""
    echo -e "${YELLOW}Usage: $0 -file <file_path> [-type road|land] [-dsn <dsn>] [-batch <size>]${NC}"
    echo ""
    echo -e "${YELLOW}Example:${NC}"
    echo "  $0 -file data/postal_codes.txt -type road"
    echo "  $0 -file data/postal_codes.txt -type land -batch 2000"
    echo "  $0 -file data/postal_codes.txt -dsn \"user:pass@tcp(localhost:3306)/dbname\""
    echo ""
    exit 1
fi

# 데이터 타입 검증
if [ "$DATA_TYPE" != "road" ] && [ "$DATA_TYPE" != "land" ]; then
    echo -e "${RED}❌ 오류: -type 은 'road' 또는 'land' 여야 합니다${NC}"
    exit 1
fi

# 파일 존재 확인
if [ ! -f "$FILE_PATH" ]; then
    echo -e "${RED}❌ 오류: 파일을 찾을 수 없습니다: $FILE_PATH${NC}"
    exit 1
fi

# 상대 경로를 절대 경로로 변환
FILE_PATH=$(cd "$(dirname "$FILE_PATH")" && pwd)/$(basename "$FILE_PATH")

# 파일 정보 출력
FILE_SIZE=$(du -h "$FILE_PATH" | cut -f1)
LINE_COUNT=$(wc -l < "$FILE_PATH")

echo -e "${BLUE}📂 파일 정보:${NC}"
echo -e "  - 경로: $FILE_PATH"
echo -e "  - 크기: $FILE_SIZE"
echo -e "  - 라인 수: $LINE_COUNT"
echo ""

# 데이터 타입 한글 표시
if [ "$DATA_TYPE" == "road" ]; then
    TYPE_KOREAN="도로명주소"
else
    TYPE_KOREAN="지번주소"
fi

echo -e "${BLUE}⚙️  설정:${NC}"
if [ -z "$DSN" ]; then
    echo -e "  - DSN: .env 파일에서 로드"
else
    echo -e "  - DSN: ${DSN%%:*}:***@..."
fi
echo -e "  - 데이터 타입: $DATA_TYPE ($TYPE_KOREAN)"
echo -e "  - 배치 사이즈: $BATCH_SIZE"
echo ""

# 실행 전 확인
read -p "계속하시겠습니까? (y/N): " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}⚠️  작업이 취소되었습니다.${NC}"
    exit 0
fi

echo ""
echo -e "${YELLOW}🔨 바이너리 빌드 중...${NC}"

# 현재 스크립트 위치 찾기
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PKG_DIR="$( dirname "$SCRIPT_DIR" )"
CMD_DIR="$PKG_DIR/cmd/postalcode-import"

# 바이너리 빌드
cd "$CMD_DIR"
go build -o "$PKG_DIR/bin/postalcode-import" .

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ 빌드 실패${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 빌드 완료${NC}"
echo ""

# 실행
echo -e "${YELLOW}🚀 Import 시작...${NC}"
echo ""

# DSN이 있으면 플래그로 전달, 없으면 .env에서 자동 로드
if [ -z "$DSN" ]; then
    "$PKG_DIR/bin/postalcode-import" \
        -file "$FILE_PATH" \
        -type "$DATA_TYPE" \
        -batch "$BATCH_SIZE"
else
    "$PKG_DIR/bin/postalcode-import" \
        -dsn "$DSN" \
        -file "$FILE_PATH" \
        -type "$DATA_TYPE" \
        -batch "$BATCH_SIZE"
fi

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✅ Import 완료!${NC}"
else
    echo ""
    echo -e "${RED}❌ Import 실패${NC}"
    exit 1
fi
