package pgutil

// postgresql 과 graphql을 연결하는 함수
import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/lib/pq"
)

func MakeArrayString[T any](values []T) string {
	var builder strings.Builder
	builder.WriteString("{")
	for i, v := range values {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprint(v))
	}
	builder.WriteString("}")
	return builder.String()
}

func MarshalPqStringArray(arr pq.StringArray) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		a := []string{}
		for _, d := range arr {
			a = append(a, d)
		}
		enc, _ := json.Marshal(a)
		w.Write(enc)
	})
}

func UnmarshalPqStringArray(v interface{}) (pq.StringArray, error) {
	switch v := v.(type) {
	case []string:
		arr := pq.StringArray{}
		for _, d := range v {
			arr = append(arr, d)
		}
		return arr, nil
	case []interface{}:
		arr := pq.StringArray{}
		for _, d := range v {
			arr = append(arr, d.(string))
		}
		return arr, nil
	default:
		return nil, fmt.Errorf("%T is not a pq.StringArray", v)
	}
}
