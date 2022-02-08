package xoidc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/jinares/xpkg/xtools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"strings"
)

func getMdMap(ctx context.Context) map[string]string {
	mds := make(map[string]string)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return mds
	}
	for k, v := range md {
		if len(v) > 0 {
			mds[strings.ToLower(k)] = v[0]
		}
	}
	return mds
}

func parseJWT(token string, iter interface{}) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return "", xtools.XErr(
			codes.InvalidArgument,
			fmt.Sprintf("parseToke: malformed jwt, expected 3 parts got %d - %s - %s", len(parts), token, parts),
			true,
		)
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("parseToken: malformed jwt payload: %v", err)
	}
	if err := json.Unmarshal(payload, iter); err != nil {
		return string(payload), xtools.XErr(
			codes.Internal,
			fmt.Sprintf("parseToken: failed to %s unmarshal claims: %v", string(payload), err),
			true,
		)
	}
	return string(payload), nil
}
