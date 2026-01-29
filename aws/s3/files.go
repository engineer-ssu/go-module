package s3

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// TransferObjectIfNotExist 는 배포 경로에 파일이 없다면 임시 경로에서 복사해옵니다.
func (s *S3Service) TransferObjectIfNotExist(filename string, sourceDir string, destDir string) error {
	destKey := fmt.Sprintf("%s/%s", s.config.DestPrefix, filename) // media 가 들어가야함

	// 1. 파일 존재 여부 확인 (HeadObject가 GetObject보다 비용이 저렴하고 효율적입니다)
	headInput := &s3.HeadObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(destKey),
	}

	_, err := s.svc.HeadObject(headInput)
	if err == nil {
		return nil
	}

	// 2. 파일이 없을 경우 CopyObject 수행 (기존 TransferObject 로직을 내재화)
	sourceKey := fmt.Sprintf("%s/%s/%s", s.config.Bucket, s.config.SourcePrefix, filename)

	err = s.TransferObject(sourceKey, destKey)
	if err != nil {
		return fmt.Errorf("failed to transfer object from %s to %s: %w", sourceKey, destKey, err)
	}

	return nil
}

// ParseImgSrc 텍스트에디터에서 로컬파일url을 S3url로 바꾸고 파일도 복사하기
func (s *S3Service) ParseImgSrc(content *string, prefix string) (*string, error) {
	if content == nil {
		return nil, errors.New("content empty")
	}
	re := regexp.MustCompile(`src="([^"]+)"[^>]*>`)
	newContent := *content
	for _, match := range re.FindAllStringSubmatch(*content, -1) {
		key := match[1]
		key = strings.ReplaceAll(key, s.config.CdnUri+"/", "")
		key = filepath.Base(key)
		bucket := s.config.Bucket
		source := filepath.Join(bucket, s.config.SourcePrefix, key)
		dest := filepath.Join(s.config.DestPrefix, prefix, key)
		err := s.TransferObject(source, dest)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		newSrc := fmt.Sprintf("%s/%s", s.config.CdnUri, dest)
		newContent = strings.ReplaceAll(newContent, match[1], newSrc)
	}
	return &newContent, nil
}

// TransferObject 임시저장소에 있는 오브젝트를 배포저장소로 복사하기
func (s *S3Service) TransferObject(sourceKey string, destKey string) error {
	input := &s3.CopyObjectInput{
		Bucket:     aws.String(s.config.Bucket),
		Key:        aws.String(destKey),
		CopySource: aws.String(sourceKey),
	}
	_, err := s.svc.CopyObject(input)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
