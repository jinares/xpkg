package xoidc

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/jinares/xpkg/xtools"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

var (
	// OIDCVerifier 备注：一定要初始化到全局
	OIDCVerifier *oidc.IDTokenVerifier
	cfg          = &Config{
		OIDCAuth: AuthConfig{
			Issuer:         "",
			AllowProjects:  []string{},
			Scopes:         []string{},
			AllowedIssuers: []string{},
		},
	}
)

//InitOIDC InitOIDC oidcissue 内网填写 https://apis.xbase.xyz ， 外网填写：https://apis.xbase.cloud
func InitOIDC(path string) error {
	err := xtools.BindYAML(path, cfg)
	if err != nil {
		fmt.Println("oidc-config-err")
		return err
	}
	if cfg.OIDCAuth.Issuer == "" {
		return xtools.XErr(codes.Internal, "oidc-config-err")
	}
	provider, err := oidc.NewProvider(context.Background(), cfg.OIDCAuth.Issuer)
	if err != nil {
		return err
	}
	oidcConfig := &oidc.Config{
		SkipClientIDCheck: true,
		SkipIssuerCheck:   true,
	}
	OIDCVerifier = provider.Verifier(oidcConfig)
	return nil
}
func InitByConfig(iss string, cfg *oidc.Config) error {
	provider, err := oidc.NewProvider(context.Background(), iss)
	if err != nil {
		return err
	}
	oidcConfig := &oidc.Config{
		SkipClientIDCheck: true,
		SkipIssuerCheck:   true,
	}
	if cfg == nil {
		cfg = oidcConfig
	}
	OIDCVerifier = provider.Verifier(oidcConfig)
	return nil
}

//IDToken IDToken
type (
	IDToken struct {
		Issue    string `protobuf:"bytes,1,opt,name=iss,json=Issue,proto3" json:"iss,omitempty"`
		ClientID string `protobuf:"bytes,2,opt,name=aud,json=Aud,proto3" json:"aud,omitempty"`
		Exp      int64  `protobuf:"bytes,3,opt,name=exp,json=Exp,proto3" json:"exp,omitempty"`
		Iat      int64  `protobuf:"bytes,4,opt,name=iat,json=Iat,proto3" json:"iat,omitempty"`
		AtHash   string `protobuf:"bytes,5,opt,name=at_hash,json=AtHash,proto3" json:"at_hash,omitempty"`
		Scope    string `protobuf:"bytes,6,opt,name=scope,json=Scope,proto3" json:"scope,omitempty"`
		// 	ClientID    string `protobuf:"bytes,7,opt,name=client_id,json=ClientID,proto3" json:"client_id,omitempty`
		ProjectID   string `protobuf:"bytes,8,opt,name=project_id,json=ProjectID,proto3" json:"project_id,omitempty"`
		ServiceID   string `protobuf:"bytes,9,opt,name=service_id,json=ServiceID,proto3" json:"service_id,omitempty"`
		ServiceName string `protobuf:"bytes,10,opt,name=service_name,json=ServiceName,proto3" json:"service_name,omitempty"`
		Sub         string `protobuf:"bytes,10,opt,name=sub,json=Sub,proto3" json:"sub,omitempty"`
	}
	OidcAuthFuncHandler func(ctx context.Context, fullmethod string, token *IDToken) error
)

func DefaultAuthHandler() OidcAuthFuncHandler {
	return func(ctx context.Context, fullmethod string, token *IDToken) (err error) {
		if token == nil {
			return nil
		}
		err = VerifyHas(token.Issue, cfg.OIDCAuth.AllowedIssuers)
		if err != nil {
			return err
		}
		err = VerifyHas(token.ProjectID, cfg.OIDCAuth.AllowProjects)
		if err != nil {
			return err
		}
		err = VerifyHas(token.Scope, cfg.OIDCAuth.Scopes)
		if err != nil {
			return err
		}
		return nil
	}
}

// Authorize 检查授权
func Authorize(ctx context.Context, fullMethodName string, mds map[string]string, handler OidcAuthFuncHandler) (context.Context, error) {
	//
	//mds := getMdMap(ctx)
	//// 没有就默认放过？
	if mds[AUTHORIZATION] == "" {
		return ctx, nil
	}

	// oidc || xlsoa
	spls := strings.SplitN(mds[AUTHORIZATION], " ", 2)
	if len(spls) != 2 || spls[0] != "Bearer" {
		return ctx, status.Error(codes.InvalidArgument, "token not found")
	}

	// token
	claims, err := VerifyIDToken(ctx, spls[1])
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "token %s verify error %s", spls[1], err)
	}
	err = handler(ctx, fullMethodName, claims)
	if err != nil {
		return ctx, err
	}
	ctx = context.WithValue(ctx, FULL_METHOD_NAME, fullMethodName)
	ctx = context.WithValue(ctx, SOA_SERVICE_ID, claims.ServiceID)
	ctx = context.WithValue(ctx, SOA_SERVICE_NAME, claims.ServiceName)

	md, _ := metadata.FromIncomingContext(ctx)

	if claims.Scope != "" {
		md.Set(X_SCOPE, claims.Scope)
	}
	if claims.ClientID != "" {
		md.Set(X_CLIENT_ID, claims.ClientID)
	}
	if claims.ProjectID != "" {
		md.Set(X_PROJECT_ID, claims.ProjectID)
	}
	if claims.ServiceName != "" {
		md.Set(X_SERVICE_NAME, claims.ServiceName)
	}
	ctx = metadata.NewIncomingContext(ctx, md)
	return ctx, nil

}

// VerifyIDToken VerifyIDToken the id token
func VerifyIDToken(ctx context.Context, token string) (*IDToken, error) {
	if nil == OIDCVerifier {
		fmt.Println("dddd")
		var claims IDToken
		_, err := parseJWT(token, &claims)
		if nil != err {
			return nil, status.Error(codes.InvalidArgument, "token not found")
		}
		return &claims, nil
	}
	tk, err := OIDCVerifier.Verify(ctx, token)
	fmt.Println("aaaaaaaaaaassss")
	if err != nil {
		return nil, err
	}
	var claims IDToken
	err = tk.Claims(&claims)
	if err != nil {
		return nil, err
	}
	return &claims, nil
}
