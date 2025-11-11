-- 우편번호별 도로명 범위 테이블 생성
CREATE TABLE IF NOT EXISTS postal_code_roads (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT 'PK',

    -- 우편번호 (주 조회 키)
    zip_code VARCHAR(5) NOT NULL COMMENT '우편번호 (5자리)',
    zip_prefix CHAR(3) NOT NULL COMMENT '우편번호 앞 3자리',

    -- 행정구역
    sido_name VARCHAR(40) NOT NULL COMMENT '시도명',
    sido_name_en VARCHAR(40) DEFAULT NULL COMMENT '시도명 영문',
    sigungu_name VARCHAR(40) NOT NULL COMMENT '시군구명',
    sigungu_name_en VARCHAR(40) DEFAULT NULL COMMENT '시군구명 영문',
    eupmyeon_name VARCHAR(40) DEFAULT NULL COMMENT '읍면명',
    eupmyeon_name_en VARCHAR(40) DEFAULT NULL COMMENT '읍면명 영문',

    -- 도로명
    road_name VARCHAR(80) NOT NULL COMMENT '도로명',
    road_name_en VARCHAR(80) DEFAULT NULL COMMENT '도로명 영문',

    -- 지하여부
    is_underground TINYINT(1) NOT NULL DEFAULT 0 COMMENT '지하여부 (0=지상, 1=지하)',

    -- 건물번호 범위
    start_building_main INT NOT NULL COMMENT '시작건물번호(주)',
    start_building_sub INT DEFAULT NULL COMMENT '시작건물번호(부)',
    end_building_main INT DEFAULT NULL COMMENT '끝건물번호(주)',
    end_building_sub INT DEFAULT NULL COMMENT '끝건물번호(부)',

    -- 범위종류
    range_type TINYINT NOT NULL COMMENT '범위종류 (1=단일번호, 2=구간, 3=다수)',

    -- 타임스탬프
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '생성일시',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정일시',

    -- 인덱스 (우편번호 조회 최적화)
    INDEX idx_zipcode (zip_code),
    INDEX idx_zip_prefix (zip_prefix),
    INDEX idx_sido (sido_name),
    INDEX idx_sigungu (sigungu_name),
    INDEX idx_road (road_name),

    -- 유니크 인덱스 (중복 방지 및 무결성 보장)
    -- 모든 필드를 포함하여 완전히 동일한 레코드만 중복으로 간주
    UNIQUE INDEX idx_postal_unique (
        zip_code, sido_name, sigungu_name, eupmyeon_name, road_name,
        is_underground, start_building_main, start_building_sub,
        end_building_main, end_building_sub, range_type
    )
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='우편번호별 도로명 범위 정보';

-- UPSERT를 위한 INSERT ... ON DUPLICATE KEY UPDATE 문 예시
-- idx_postal_unique (우편번호 + 시도 + 시군구 + 도로명 + 시작건물번호)로 중복 판단
-- INSERT INTO postal_code_roads (
--     zip_code, zip_prefix, sido_name, sido_name_en, sigungu_name, sigungu_name_en,
--     eupmyeon_name, eupmyeon_name_en, road_name, road_name_en,
--     is_underground, start_building_main, start_building_sub,
--     end_building_main, end_building_sub, range_type
-- ) VALUES (
--     ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
-- ) ON DUPLICATE KEY UPDATE
--     sido_name_en = VALUES(sido_name_en),
--     sigungu_name_en = VALUES(sigungu_name_en),
--     eupmyeon_name = VALUES(eupmyeon_name),
--     eupmyeon_name_en = VALUES(eupmyeon_name_en),
--     road_name_en = VALUES(road_name_en),
--     is_underground = VALUES(is_underground),
--     start_building_sub = VALUES(start_building_sub),
--     end_building_main = VALUES(end_building_main),
--     end_building_sub = VALUES(end_building_sub),
--     range_type = VALUES(range_type),
--     updated_at = CURRENT_TIMESTAMP;
