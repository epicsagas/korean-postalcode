-- 우편번호별 지번주소 범위 테이블 생성
CREATE TABLE IF NOT EXISTS postal_code_lands (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT 'PK',

    -- 우편번호 (주 조회 키)
    zip_code VARCHAR(5) NOT NULL COMMENT '우편번호 (5자리)',
    zip_prefix CHAR(3) NOT NULL COMMENT '우편번호 앞 3자리',

    -- 행정구역
    sido_name VARCHAR(40) NOT NULL COMMENT '시도명',
    sido_name_en VARCHAR(40) DEFAULT NULL COMMENT '시도명 영문',
    sigungu_name VARCHAR(40) NOT NULL COMMENT '시군구명',
    sigungu_name_en VARCHAR(40) DEFAULT NULL COMMENT '시군구명 영문',
    eupmyeondong_name VARCHAR(40) NOT NULL COMMENT '읍면동명',
    eupmyeondong_name_en VARCHAR(40) DEFAULT NULL COMMENT '읍면동명 영문',

    -- 리명
    ri_name VARCHAR(40) DEFAULT NULL COMMENT '리명',

    -- 산여부
    is_mountain TINYINT(1) NOT NULL DEFAULT 0 COMMENT '산여부 (0=일반, 1=산)',

    -- 행정동
    haengjeongdong_name VARCHAR(40) DEFAULT NULL COMMENT '행정동명',

    -- 주번지 범위
    start_jibun_main INT NOT NULL COMMENT '시작주번지',
    start_jibun_sub INT DEFAULT NULL COMMENT '시작부번지',
    end_jibun_main INT DEFAULT NULL COMMENT '끝주번지',
    end_jibun_sub INT DEFAULT NULL COMMENT '끝부번지',

    -- 타임스탬프
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '생성일시',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정일시',

    -- 인덱스 (우편번호 조회 최적화)
    INDEX idx_land_zipcode (zip_code),
    INDEX idx_land_zip_prefix (zip_prefix),
    INDEX idx_land_sido (sido_name),
    INDEX idx_land_sigungu (sigungu_name),
    INDEX idx_land_eupmyeondong (eupmyeondong_name),
    INDEX idx_land_ri (ri_name),

    -- 유니크 인덱스 (중복 방지 및 무결성 보장)
    -- 모든 필드를 포함하여 완전히 동일한 레코드만 중복으로 간주
    UNIQUE INDEX idx_land_unique (
        zip_code, sido_name, sigungu_name, eupmyeondong_name, ri_name,
        is_mountain, start_jibun_main, start_jibun_sub,
        end_jibun_main, end_jibun_sub
    )
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='우편번호별 지번주소 범위 정보';

-- UPSERT를 위한 INSERT ... ON DUPLICATE KEY UPDATE 문 예시
-- idx_land_unique (우편번호 + 시도 + 시군구 + 읍면동 + 리명 + 산여부 + 시작주번지)로 중복 판단
-- INSERT INTO postal_code_lands (
--     zip_code, zip_prefix, sido_name, sido_name_en, sigungu_name, sigungu_name_en,
--     eupmyeondong_name, eupmyeondong_name_en, ri_name, is_mountain,
--     haengjeongdong_name, start_jibun_main, start_jibun_sub,
--     end_jibun_main, end_jibun_sub
-- ) VALUES (
--     ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
-- ) ON DUPLICATE KEY UPDATE
--     sido_name_en = VALUES(sido_name_en),
--     sigungu_name_en = VALUES(sigungu_name_en),
--     eupmyeondong_name_en = VALUES(eupmyeondong_name_en),
--     haengjeongdong_name = VALUES(haengjeongdong_name),
--     start_jibun_sub = VALUES(start_jibun_sub),
--     end_jibun_main = VALUES(end_jibun_main),
--     end_jibun_sub = VALUES(end_jibun_sub),
--     updated_at = CURRENT_TIMESTAMP;
