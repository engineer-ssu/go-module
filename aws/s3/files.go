package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
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
