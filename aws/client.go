// pkg/aws/client.go (예시)
package aws

type Config struct {
	Region          string
	AccessKey       string
	SecretAccessKey string
}

type Client struct {
	cfg Config
}

// 라이브러리 사용자가 설정을 직접 주입함
func NewClient(cfg Config) *Client {
	return &Client{cfg: cfg}
}
