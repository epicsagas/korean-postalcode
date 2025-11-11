package importer

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/epicsagas/korean-postalcode"
	"github.com/epicsagas/korean-postalcode/internal/service"
)

// Importer는 파일에서 우편번호 데이터를 가져오는 기능을 제공합니다.
type Importer interface {
	// 도로명주소 관련 메서드
	// ImportFromFile은 파일에서 도로명주소 데이터를 가져와 DB에 저장합니다.
	ImportFromFile(filePath string, batchSize int, progressFn postalcode.ProgressFunc) (*postalcode.ImportResult, error)

	// ParseFile은 파일을 파싱하여 postalcode.PostalCodeRoad 슬라이스로 변환합니다.
	ParseFile(filePath string) ([]postalcode.PostalCodeRoad, error)

	// 지번주소 관련 메서드
	// ImportLandFromFile은 파일에서 지번주소 데이터를 가져와 DB에 저장합니다.
	ImportLandFromFile(filePath string, batchSize int, progressFn postalcode.ProgressFunc) (*postalcode.ImportResult, error)

	// ParseLandFile은 파일을 파싱하여 postalcode.PostalCodeLand 슬라이스로 변환합니다.
	ParseLandFile(filePath string) ([]postalcode.PostalCodeLand, error)
}

// importer는 Importer 인터페이스 구현입니다.
type importer struct {
	service service.Service
}

// New는 새로운 Importer를 생성합니다.
func New(svc service.Service) Importer {
	return &importer{service: svc}
}

// countDataLines counts the number of data lines in a file (excluding header)
func countDataLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	firstLine := true

	for scanner.Scan() {
		if firstLine {
			firstLine = false
			continue // Skip header
		}
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}

// ImportFromFile은 파일에서 우편번호 데이터를 가져와 DB에 저장합니다.
func (imp *importer) ImportFromFile(filePath string, batchSize int, progressFn postalcode.ProgressFunc) (*postalcode.ImportResult, error) {
	startTime := time.Now()

	if batchSize <= 0 {
		batchSize = 1000
	}

	// Count total lines in file (excluding header)
	totalLines, err := countDataLines(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to count lines: %w", err)
	}

	// 파일 파싱
	roads, err := imp.ParseFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("file parsing failed: %w", err)
	}

	totalCount := 0
	errorCount := 0

	// 배치 처리
	for i := 0; i < len(roads); i += batchSize {
		end := i + batchSize
		if end > len(roads) {
			end = len(roads)
		}

		batch := roads[i:end]

		// DB에 저장
		if err := imp.service.BatchUpsert(batch); err != nil {
			fmt.Printf("❌ 배치 %d-%d 저장 실패: %v\n", i, end, err)
			errorCount += len(batch)
		} else {
			totalCount += len(batch)
		}

		// 진행 상황 보고
		if progressFn != nil {
			progressFn(i+len(batch), len(roads))
		}
	}

	// Parse errors = total lines - successfully parsed records
	parseErrors := totalLines - len(roads)
	errorCount += parseErrors

	duration := time.Since(startTime)
	return &postalcode.ImportResult{
		TotalCount: totalCount,
		ErrorCount: errorCount,
		Duration:   duration.String(),
	}, nil
}

// ParseFile은 파일을 파싱하여 PostalCodeRoad 슬라이스로 변환합니다.
func (imp *importer) ParseFile(filePath string) ([]postalcode.PostalCodeRoad, error) {
	// 파일 열기
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// CSV 리더 생성 (파이프 구분자)
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = '|'
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	// 헤더 읽기 (첫 줄 스킵)
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	var roads []postalcode.PostalCodeRoad
	lineNumber := 1 // 헤더 이후부터
	var parseErrors []string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			parseErrors = append(parseErrors, fmt.Sprintf("라인 %d: CSV 파싱 에러 - %v", lineNumber, err))
			lineNumber++
			continue
		}

		// 필드 수 검증
		if len(record) < 15 {
			parseErrors = append(parseErrors, fmt.Sprintf("라인 %d: 필드 수 부족 (필요: 15, 실제: %d)", lineNumber, len(record)))
			lineNumber++
			continue
		}

		// 데이터 파싱
		zipCode := strings.TrimSpace(record[0])
		zipPrefix := ""
		if len(zipCode) >= 3 {
			zipPrefix = zipCode[:3]
		}

		road := postalcode.PostalCodeRoad{
			ZipCode:        zipCode,
			ZipPrefix:      zipPrefix,
			SidoName:       strings.TrimSpace(record[1]),
			SidoNameEn:     strings.TrimSpace(record[2]),
			SigunguName:    strings.TrimSpace(record[3]),
			SigunguNameEn:  strings.TrimSpace(record[4]),
			EupmyeonName:   strings.TrimSpace(record[5]),
			EupmyeonNameEn: strings.TrimSpace(record[6]),
			RoadName:       strings.TrimSpace(record[7]),
			RoadNameEn:     strings.TrimSpace(record[8]),
		}

		// 지하여부 파싱
		if underground := strings.TrimSpace(record[9]); underground == "1" {
			road.IsUnderground = true
		}

		// 시작건물번호(주) 파싱
		if startMain := strings.TrimSpace(record[10]); startMain != "" {
			if val, err := strconv.Atoi(startMain); err == nil {
				road.StartBuildingMain = val
			}
		}

		// 시작건물번호(부) 파싱
		if startSub := strings.TrimSpace(record[11]); startSub != "" && startSub != "0" {
			if val, err := strconv.Atoi(startSub); err == nil {
				road.StartBuildingSub = &val
			}
		}

		// 끝건물번호(주) 파싱
		if endMain := strings.TrimSpace(record[12]); endMain != "" {
			if val, err := strconv.Atoi(endMain); err == nil {
				road.EndBuildingMain = &val
			}
		}

		// 끝건물번호(부) 파싱
		if endSub := strings.TrimSpace(record[13]); endSub != "" && endSub != "0" {
			if val, err := strconv.Atoi(endSub); err == nil {
				road.EndBuildingSub = &val
			}
		}

		// 범위종류 파싱
		if rangeType := strings.TrimSpace(record[14]); rangeType != "" {
			if val, err := strconv.Atoi(rangeType); err == nil {
				road.RangeType = int8(val)
			}
		}

		roads = append(roads, road)
		lineNumber++
	}

	// 파싱 에러가 있으면 출력
	if len(parseErrors) > 0 {
		fmt.Printf("⚠️  파싱 중 %d개 에러 발생:\n", len(parseErrors))
		for i, errMsg := range parseErrors {
			if i < 10 { // 최대 10개만 출력
				fmt.Printf("  - %s\n", errMsg)
			}
		}
		if len(parseErrors) > 10 {
			fmt.Printf("  ... 외 %d개\n", len(parseErrors)-10)
		}
	}

	return roads, nil
}

// ============================================================
// 지번주소 관련 메서드
// ============================================================

// ImportLandFromFile은 파일에서 지번주소 데이터를 가져와 DB에 저장합니다.
func (imp *importer) ImportLandFromFile(filePath string, batchSize int, progressFn postalcode.ProgressFunc) (*postalcode.ImportResult, error) {
	startTime := time.Now()

	if batchSize <= 0 {
		batchSize = 1000
	}

	// Count total lines in file (excluding header)
	totalLines, err := countDataLines(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to count lines: %w", err)
	}

	// 파일 파싱
	lands, err := imp.ParseLandFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("file parsing failed: %w", err)
	}

	totalCount := 0
	errorCount := 0

	// 배치 처리
	for i := 0; i < len(lands); i += batchSize {
		end := i + batchSize
		if end > len(lands) {
			end = len(lands)
		}

		batch := lands[i:end]

		// DB에 저장
		if err := imp.service.BatchUpsertLand(batch); err != nil {
			fmt.Printf("❌ 배치 %d-%d 저장 실패: %v\n", i, end, err)
			errorCount += len(batch)
		} else {
			totalCount += len(batch)
		}

		// 진행 상황 보고
		if progressFn != nil {
			progressFn(i+len(batch), len(lands))
		}
	}

	// Parse errors = total lines - successfully parsed records
	parseErrors := totalLines - len(lands)
	errorCount += parseErrors

	duration := time.Since(startTime)
	return &postalcode.ImportResult{
		TotalCount: totalCount,
		ErrorCount: errorCount,
		Duration:   duration.String(),
	}, nil
}

// ParseLandFile은 파일을 파싱하여 PostalCodeLand 슬라이스로 변환합니다.
func (imp *importer) ParseLandFile(filePath string) ([]postalcode.PostalCodeLand, error) {
	// 파일 열기
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// CSV 리더 생성 (파이프 구분자)
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = '|'
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	// 헤더 읽기 (첫 줄 스킵)
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	var lands []postalcode.PostalCodeLand
	lineNumber := 1 // 헤더 이후부터
	var parseErrors []string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			parseErrors = append(parseErrors, fmt.Sprintf("라인 %d: CSV 파싱 에러 - %v", lineNumber, err))
			lineNumber++
			continue
		}

		// 필드 수 검증
		if len(record) < 14 {
			parseErrors = append(parseErrors, fmt.Sprintf("라인 %d: 필드 수 부족 (필요: 14, 실제: %d)", lineNumber, len(record)))
			lineNumber++
			continue
		}

		// 데이터 파싱
		zipCode := strings.TrimSpace(record[0])
		zipPrefix := ""
		if len(zipCode) >= 3 {
			zipPrefix = zipCode[:3]
		}

		land := postalcode.PostalCodeLand{
			ZipCode:            zipCode,
			ZipPrefix:          zipPrefix,
			SidoName:           strings.TrimSpace(record[1]),
			SidoNameEn:         strings.TrimSpace(record[2]),
			SigunguName:        strings.TrimSpace(record[3]),
			SigunguNameEn:      strings.TrimSpace(record[4]),
			EupmyeondongName:   strings.TrimSpace(record[5]),
			EupmyeondongNameEn: strings.TrimSpace(record[6]),
			RiName:             strings.TrimSpace(record[7]),
			HaengjeongdongName: strings.TrimSpace(record[9]),
		}

		// 산여부 파싱
		if mountain := strings.TrimSpace(record[8]); mountain == "1" {
			land.IsMountain = true
		}

		// 시작주번지 파싱
		if startMain := strings.TrimSpace(record[10]); startMain != "" {
			if val, err := strconv.Atoi(startMain); err == nil {
				land.StartJibunMain = val
			}
		}

		// 시작부번지 파싱
		if startSub := strings.TrimSpace(record[11]); startSub != "" && startSub != "0" {
			if val, err := strconv.Atoi(startSub); err == nil {
				land.StartJibunSub = &val
			}
		}

		// 끝주번지 파싱
		if endMain := strings.TrimSpace(record[12]); endMain != "" {
			if val, err := strconv.Atoi(endMain); err == nil {
				land.EndJibunMain = &val
			}
		}

		// 끝부번지 파싱
		if endSub := strings.TrimSpace(record[13]); endSub != "" && endSub != "0" {
			if val, err := strconv.Atoi(endSub); err == nil {
				land.EndJibunSub = &val
			}
		}

		lands = append(lands, land)
		lineNumber++
	}

	// 파싱 에러가 있으면 출력
	if len(parseErrors) > 0 {
		fmt.Printf("⚠️  파싱 중 %d개 에러 발생:\n", len(parseErrors))
		for i, errMsg := range parseErrors {
			if i < 10 { // 최대 10개만 출력
				fmt.Printf("  - %s\n", errMsg)
			}
		}
		if len(parseErrors) > 10 {
			fmt.Printf("  ... 외 %d개\n", len(parseErrors)-10)
		}
	}

	return lands, nil
}
