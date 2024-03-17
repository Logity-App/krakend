package grpc

import (
	"context"
	"fmt"
	ssov1 "github.com/Logity-App/contracts/gen/go/sso"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Client struct {
	api ssov1.AuthClient
}

func New(
	ctx context.Context,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &Client{
		api: ssov1.NewAuthClient(cc),
	}, nil
}

func (c *Client) VerifyNewPhoneNumber(ctx context.Context, phone string) (*ssov1.VerifyNewPhoneNumberResponse, error) {
	resp, err := c.api.VerifyNewPhoneNumber(ctx, &ssov1.VerifyNewPhoneNumberRequest{
		Phone: phone,
	})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return resp, nil
}

func (c *Client) SendSmsCode(ctx context.Context) (*ssov1.Empty, error) {
	resp, err := c.api.SendSmsCode(ctx, &ssov1.Empty{})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return resp, nil
}

func (c *Client) SignUpByPhone(ctx context.Context, phone string, birthdayDate string, defaultTag string) (*ssov1.SignUpByPhoneResponse, error) {
	resp, err := c.api.SignUpByPhone(ctx, &ssov1.SignUpByPhoneRequest{
		Phone:        phone,
		BirthdayDate: birthdayDate,
		DefaultTag:   defaultTag,
	})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return resp, nil
}

func (c *Client) VerifyPhoneNumber(ctx context.Context, phone string) (*ssov1.VerifyPhoneNumberResponse, error) {
	resp, err := c.api.VerifyPhoneNumber(ctx, &ssov1.VerifyPhoneNumberRequest{
		Phone: phone,
	})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return resp, nil
}

func (c *Client) SignInByPhone(ctx context.Context, phone string, smsCode string) (*ssov1.SignInByPhoneResponse, error) {
	resp, err := c.api.SignInByPhone(ctx, &ssov1.SignInByPhoneRequest{
		Phone:   phone,
		SmsCode: smsCode,
	})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return resp, nil
}

func InterceptorLogger() grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		fmt.Println(msg) // TODO add logger
	})
}
