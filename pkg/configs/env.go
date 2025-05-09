package configs

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type EnvConfig struct {
	Server struct {
		Port      string `mapstructure:"PORT"`
		WorkerURL string `mapstructure:"WORKER_URL"`
	}
	AWS struct {
		AccessKeyID      string `mapstructure:"AWS_ACCESS_KEY_ID"`
		SecretAccessKey  string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
		Region           string `mapstructure:"AWS_REGION"`
		DynamoDBEndpoint string `mapstructure:"AWS_DYNAMODB_ENDPOINT"`
		Tables           struct {
			OCRCache string `mapstructure:"AWS_DYNAMODB_TABLE_OCR_CACHE"`
		}
	}
	Naver struct {
		ClientID     string `mapstructure:"NAVER_CLIENT_ID"`
		ClientSecret string `mapstructure:"NAVER_CLIENT_SECRET"`
		SearchURL    string `mapstructure:"NAVER_SEARCH_URL"`
	}
	OCR struct {
		TesseractPath string `mapstructure:"OCR_TESSERACT_PATH"`
		TempDir       string `mapstructure:"OCR_TEMP_DIR"`
	}
}

var (
	configInstance *EnvConfig
	once           sync.Once
)

// loadConfig는 환경 변수를 로드하고 검증하는 내부 함수
func loadConfig() *EnvConfig {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.AutomaticEnv()

	// 필수 환경 변수 확인
	requiredEnvVars := []string{
		"PORT",
		"WORKER_URL",
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_REGION",
		"AWS_DYNAMODB_ENDPOINT",
		"AWS_DYNAMODB_TABLE_OCR_CACHE",
		"NAVER_CLIENT_ID",
		"NAVER_CLIENT_SECRET",
	}

	missingVars := []string{}
	for _, envVar := range requiredEnvVars {
		if !viper.IsSet(envVar) {
			missingVars = append(missingVars, envVar)
		}
	}

	if len(missingVars) > 0 {
		log.Fatalf("필수 환경 변수가 설정되지 않았습니다: %s", strings.Join(missingVars, ", "))
	}

	// 기본값 설정
	viper.SetDefault("NAVER_SEARCH_URL", "https://openapi.naver.com/v1/search/blog.json")
	viper.SetDefault("OCR_TESSERACT_PATH", "/usr/local/bin/tesseract")
	viper.SetDefault("OCR_TEMP_DIR", "/tmp")

	// 환경 변수 키-구조체 필드 매핑 정의
	config := &EnvConfig{}
	envMapping := map[string]*string{
		"PORT":                         &config.Server.Port,
		"WORKER_URL":                   &config.Server.WorkerURL,
		"AWS_ACCESS_KEY_ID":            &config.AWS.AccessKeyID,
		"AWS_SECRET_ACCESS_KEY":        &config.AWS.SecretAccessKey,
		"AWS_REGION":                   &config.AWS.Region,
		"AWS_DYNAMODB_ENDPOINT":        &config.AWS.DynamoDBEndpoint,
		"AWS_DYNAMODB_TABLE_OCR_CACHE": &config.AWS.Tables.OCRCache,
		"NAVER_CLIENT_ID":              &config.Naver.ClientID,
		"NAVER_CLIENT_SECRET":          &config.Naver.ClientSecret,
		"NAVER_SEARCH_URL":             &config.Naver.SearchURL,
		"OCR_TESSERACT_PATH":           &config.OCR.TesseractPath,
		"OCR_TEMP_DIR":                 &config.OCR.TempDir,
	}

	fmt.Println("환경 변수 로드 중...")
	// 필드에 환경 변수 값 매핑 - 문자열 필드
	for key, field := range envMapping {
		*field = viper.GetString(key)
		fmt.Printf("%s: '%s'\n", key, *field) // 디버깅용
	}

	return config
}

// GetConfig는 EnvConfig의 싱글톤 인스턴스를 반환합니다.
// 처음 호출 시에만 환경 변수를 로드하고 이후 호출에서는 캐시된 인스턴스를 반환합니다.
func GetConfig() *EnvConfig {
	once.Do(func() {
		configInstance = loadConfig()
		fmt.Println("환경 변수 로드 완료")
	})
	return configInstance
}
