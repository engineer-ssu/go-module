package s3

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// Config 는 S3 관련 설정을 외부에서 주입받기 위한 구조체입니다.
type Config struct {
	Bucket      string
	TempPrefix  string // 임시 저장소 경로 (예: "temp")
	MediaPrefix string // 배포 저장소 경로 (예: "media")
	CdnUri      string
}

// S3Service 는 라이브러리의 메인 구조체입니다.
type S3Service struct {
	svc    s3iface.S3API // 인터페이스를 사용하여 테스트 시 Mocking 가능하게 함
	config Config
}

// NewS3Service 는 외부에서 S3 API와 설정을 주입받아 객체를 생성합니다.
func NewS3Service(svc s3iface.S3API, cfg Config) *S3Service {
	return &S3Service{
		svc:    svc,
		config: cfg,
	}
}
