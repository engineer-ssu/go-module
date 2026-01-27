# Go-Module (Internal Shared Library)

이 프로젝트는 고(Go) 기반 백엔드 컨테이너 서비스에서 공통으로 사용되는 라이브러리입니다.

## 📁 패키지 구성 및 역할
- ***`/db/clause`**: SQL Join, OrderBy 등 동적 쿼리 구문 생성을 담당합니다.
- **`/db/schema`**: 모델 구조체(Struct)의 태그를 분석하여 Select 필드를 자동 추출합니다.
- **`/db/pgutil`**: PostgreSQL 전용 데이터 타입(StringArray 등) 변환 기능을 제공합니다.
- **`/db/utils`**: 특정 프레임워크에 의존하지 않는 순수 유틸리티 모음입니다.

## 🚀 시작하기

### 1. 설치
```bash
go get github.com/engineer-ssu/go-module@X.Y.Z
```


Requires Go >= 1.24 (toolchain supported)

### 2. import 적용
```bash
import (
  clause "github.com/engineer-ssu/go-module/db/clause"
)
```
