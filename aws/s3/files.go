package s3

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/engineer-ssu/go-module/config"
)

// TransferObjectIfNotExist 는 배포 경로에 파일이 없다면 임시 경로에서 복사해옵니다.
func (s *S3Service) TransferObjectIfNotExist(ctx context.Context, filename string) error {
	destKey := fmt.Sprintf("%s/%s", s.config.MediaPrefix, filename)

	// 1. 파일 존재 여부 확인 (HeadObject가 GetObject보다 비용이 저렴하고 효율적입니다)
	headInput := &s3.HeadObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(destKey),
	}

	_, err := s.svc.HeadObjectWithContext(ctx, headInput)
	if err == nil {
		// 파일이 이미 존재함
		return nil
	}

	// 2. 파일이 없을 경우 CopyObject 수행 (기존 TransferObject 로직을 내재화)
	sourceKey := fmt.Sprintf("%s/%s/%s", s.config.Bucket, s.config.TempPrefix, filename)

	copyInput := &s3.CopyObjectInput{
		Bucket:     aws.String(s.config.Bucket),
		CopySource: aws.String(sourceKey),
		Key:        aws.String(destKey),
	}

	_, err = s.svc.CopyObjectWithContext(ctx, copyInput)
	if err != nil {
		return fmt.Errorf("failed to transfer object from %s to %s: %w", sourceKey, destKey, err)
	}

	return nil
}

// ParseImgSrc 텍스트에디터에서 로컬파일url을 S3url로 바꾸고 파일도 복사하기
func (s *S3Service) ParseImgSrc(cfg config.Config, content *string, prefix string) *string {
	if content == nil {
		return nil
	}
	re := regexp.MustCompile(`src="([^"]+)"[^>]*>`)
	newContent := *content
	for _, match := range re.FindAllStringSubmatch(*content, -1) {
		key := match[1]
		key = strings.ReplaceAll(key, cfg.String("cdn_uri")+"/", "")
		key = filepath.Base(key)
		bucket := cfg.String("s3_bucket")
		source := filepath.Join(bucket, cfg.String("s3_temp_prefix"), key)
		dest := filepath.Join("media", prefix, key)
		err := s.TransferObject(bucket, source, dest)
		if err != nil {
			fmt.Println(err)
		}
		newSrc := fmt.Sprintf("%s/%s", cfg.String("cdn_uri"), dest)
		newContent = strings.ReplaceAll(newContent, match[1], newSrc)
	}
	return &newContent
}

// TransferObject 임시저장소에 있는 오브젝트를 배포저장소로 복사하기
func (s *S3Service) TransferObject(bucket string, source string, dest string) error {
	input := &s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		Key:        aws.String(dest),
		CopySource: aws.String(source),
	}
	_, err := s.svc.CopyObject(input)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
