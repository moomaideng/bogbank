// Package objectstorage is the S3 API client.
package objectstorage

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Config is the S3 connection. Endpoint is a full URL (http:// or https://).
type Config struct {
	Endpoint        string
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
}

// Client talks to one bucket.
type Client struct {
	raw     *s3.Client
	presign *s3.PresignClient
	bucket  string
}

// New builds a path-style client so RustFS and other non-AWS endpoints resolve.
func New(ctx context.Context, cfg Config) (*Client, error) {
	region := cfg.Region
	if region == "" {
		region = "auto"
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("load s3 config: %w", err)
	}

	raw := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(strings.TrimRight(cfg.Endpoint, "/"))
		o.UsePathStyle = true
	})

	return &Client{
		raw:     raw,
		presign: s3.NewPresignClient(raw),
		bucket:  cfg.Bucket,
	}, nil
}

// PingContext checks that the bucket exists and is reachable.
func (c *Client) PingContext(ctx context.Context) error {
	_, err := c.raw.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(c.bucket),
	})
	if err != nil {
		return fmt.Errorf("head bucket %q: %w", c.bucket, err)
	}
	return nil
}

// Put writes body at key, replacing any existing object.
func (c *Client) Put(ctx context.Context, key string, body io.Reader, contentType *string) error {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: contentType,
	}
	if _, err := c.raw.PutObject(ctx, input); err != nil {
		return fmt.Errorf("put %q: %w", key, err)
	}
	return nil
}

// Get returns the object body. The caller closes it.
func (c *Client) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	out, err := c.raw.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("get %q: %w", key, err)
	}
	return out.Body, nil
}

// Delete removes key. A missing key is not an error.
func (c *Client) Delete(ctx context.Context, key string) error {
	_, err := c.raw.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("delete %q: %w", key, err)
	}
	return nil
}

// PresignGet returns a URL that GETs key until ttl elapses.
func (c *Client) PresignGet(ctx context.Context, key string, ttl time.Duration) (string, error) {
	req, err := c.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", fmt.Errorf("presign %q: %w", key, err)
	}
	return req.URL, nil
}
