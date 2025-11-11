package postalcode

import (
	"time"
)

// PostalCodeRoad는 우편번호별 도로명 범위 정보를 나타냅니다.
// 행정안전부 도로명주소 데이터를 저장합니다.
// @Description 한국 우편번호 및 도로명 주소 정보 (행정안전부 도로명주소 데이터)
type PostalCodeRoad struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement" example:"1"`

	// 우편번호 (주 조회 키)
	ZipCode   string `json:"zip_code" gorm:"type:varchar(5);not null;index:idx_zipcode;uniqueIndex:idx_postal_unique,priority:1" example:"01000"`
	ZipPrefix string `json:"zip_prefix" gorm:"type:char(3);not null;index:idx_zip_prefix" example:"010"`

	// 행정구역
	SidoName       string `json:"sido_name" gorm:"type:varchar(40);not null;index:idx_sido;uniqueIndex:idx_postal_unique,priority:2" example:"서울특별시"`
	SidoNameEn     string `json:"sido_name_en" gorm:"type:varchar(40)" example:"Seoul"`
	SigunguName    string `json:"sigungu_name" gorm:"type:varchar(40);not null;index:idx_sigungu;uniqueIndex:idx_postal_unique,priority:3" example:"강북구"`
	SigunguNameEn  string `json:"sigungu_name_en" gorm:"type:varchar(40)" example:"Gangbuk-gu"`
	EupmyeonName   string `json:"eupmyeon_name" gorm:"type:varchar(40)" example:""`
	EupmyeonNameEn string `json:"eupmyeon_name_en" gorm:"type:varchar(40)" example:""`

	// 도로명
	RoadName   string `json:"road_name" gorm:"type:varchar(80);not null;index:idx_road;uniqueIndex:idx_postal_unique,priority:4" example:"삼양로177길"`
	RoadNameEn string `json:"road_name_en" gorm:"type:varchar(80)" example:"Samyang-ro 177-gil"`

	// 지하여부
	IsUnderground bool `json:"is_underground" gorm:"type:tinyint(1);default:0" example:"false"`

	// 건물번호 범위
	StartBuildingMain int  `json:"start_building_main" gorm:"type:int;not null;uniqueIndex:idx_postal_unique,priority:5" example:"93"`
	StartBuildingSub  *int `json:"start_building_sub" gorm:"type:int" example:"0"`
	EndBuildingMain   *int `json:"end_building_main" gorm:"type:int" example:"126"`
	EndBuildingSub    *int `json:"end_building_sub" gorm:"type:int" example:"0"`

	// 범위종류
	RangeType int8 `json:"range_type" gorm:"type:tinyint;not null" example:"3"`

	// 타임스탬프
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
}

// TableName은 테이블 이름을 명시적으로 지정합니다.
func (PostalCodeRoad) TableName() string {
	return "postal_code_roads"
}

// SearchParams는 우편번호 검색 파라미터입니다.
// @Description 우편번호 복합 검색 파라미터
type SearchParams struct {
	ZipCode     string `json:"zip_code" form:"zip_code" example:"01000"`
	ZipPrefix   string `json:"zip_prefix" form:"zip_prefix" example:"010"`
	SidoName    string `json:"sido_name" form:"sido_name" example:"서울특별시"`
	SigunguName string `json:"sigungu_name" form:"sigungu_name" example:"강북구"`
	RoadName    string `json:"road_name" form:"road_name" example:"삼양로"`
	Page        int    `json:"page" form:"page" example:"1"`
	Limit       int    `json:"limit" form:"limit" example:"10"`
}

// PostalCodeLand는 우편번호별 지번주소 범위 정보를 나타냅니다.
// 행정안전부 지번주소 데이터를 저장합니다.
// @Description 한국 우편번호 및 지번주소 정보 (행정안전부 지번주소 데이터)
type PostalCodeLand struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement" example:"1"`

	// 우편번호 (주 조회 키)
	ZipCode   string `json:"zip_code" gorm:"type:varchar(5);not null;index:idx_land_zipcode;uniqueIndex:idx_land_unique,priority:1" example:"25627"`
	ZipPrefix string `json:"zip_prefix" gorm:"type:char(3);not null;index:idx_land_zip_prefix" example:"256"`

	// 행정구역
	SidoName           string `json:"sido_name" gorm:"type:varchar(40);not null;index:idx_land_sido;uniqueIndex:idx_land_unique,priority:2" example:"강원특별자치도"`
	SidoNameEn         string `json:"sido_name_en" gorm:"type:varchar(40)" example:"Gangwon-do"`
	SigunguName        string `json:"sigungu_name" gorm:"type:varchar(40);not null;index:idx_land_sigungu;uniqueIndex:idx_land_unique,priority:3" example:"강릉시"`
	SigunguNameEn      string `json:"sigungu_name_en" gorm:"type:varchar(40)" example:"Gangneung-si"`
	EupmyeondongName   string `json:"eupmyeondong_name" gorm:"type:varchar(40);not null;index:idx_land_eupmyeondong;uniqueIndex:idx_land_unique,priority:4" example:"강동면"`
	EupmyeondongNameEn string `json:"eupmyeondong_name_en" gorm:"type:varchar(40)" example:"Gangdong-myeon"`

	// 리명
	RiName string `json:"ri_name" gorm:"type:varchar(40);index:idx_land_ri;uniqueIndex:idx_land_unique,priority:5" example:"모전리"`

	// 산여부
	IsMountain bool `json:"is_mountain" gorm:"type:tinyint(1);default:0;uniqueIndex:idx_land_unique,priority:6" example:"false"`

	// 행정동
	HaengjeongdongName string `json:"haengjeongdong_name" gorm:"type:varchar(40)" example:""`

	// 주번지 범위
	StartJibunMain int  `json:"start_jibun_main" gorm:"type:int;not null;uniqueIndex:idx_land_unique,priority:7" example:"2"`
	StartJibunSub  *int `json:"start_jibun_sub" gorm:"type:int" example:"3"`
	EndJibunMain   *int `json:"end_jibun_main" gorm:"type:int" example:"878"`
	EndJibunSub    *int `json:"end_jibun_sub" gorm:"type:int" example:"0"`

	// 타임스탬프
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime" example:"2024-01-01T00:00:00Z"`
}

// TableName은 테이블 이름을 명시적으로 지정합니다.
func (PostalCodeLand) TableName() string {
	return "postal_code_lands"
}

// SearchParamsLand는 지번주소 우편번호 검색 파라미터입니다.
// @Description 지번주소 우편번호 복합 검색 파라미터
type SearchParamsLand struct {
	ZipCode          string `json:"zip_code" form:"zip_code" example:"25627"`
	ZipPrefix        string `json:"zip_prefix" form:"zip_prefix" example:"256"`
	SidoName         string `json:"sido_name" form:"sido_name" example:"강원특별자치도"`
	SigunguName      string `json:"sigungu_name" form:"sigungu_name" example:"강릉시"`
	EupmyeondongName string `json:"eupmyeondong_name" form:"eupmyeondong_name" example:"강동면"`
	RiName           string `json:"ri_name" form:"ri_name" example:"모전리"`
	Page             int    `json:"page" form:"page" example:"1"`
	Limit            int    `json:"limit" form:"limit" example:"10"`
}

// ============================================================
// Import 관련 타입 (Import Types)
// ============================================================

// ImportResult는 import 결과를 나타냅니다.
type ImportResult struct {
	TotalCount int
	ErrorCount int
	Duration   string
}

// ProgressFunc는 진행 상황을 보고하는 콜백 함수입니다.
type ProgressFunc func(current, total int)
