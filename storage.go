package toold

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/yankeguo/rg"
)

type Storage struct {
	cos *cos.Client
}

func NewStorage(opts Options) (st *Storage, err error) {
	defer rg.Guard(&err)

	st = &Storage{}

	st.cos = cos.NewClient(
		&cos.BaseURL{
			BucketURL: rg.Must(url.Parse(opts.COS.BucketURL)),
		},
		&http.Client{
			Timeout: 60 * time.Second,
			Transport: &cos.AuthorizationTransport{
				SecretID:  opts.COS.SecretID,
				SecretKey: opts.COS.SecretKey,
			},
		},
	)

	return
}

func (s *Storage) CreateSignedURL(ctx context.Context, file string, ttl time.Duration) (u string, err error) {
	defer rg.Guard(&err)

	u = rg.Must(s.cos.Object.GetPresignedURL3(ctx, "GET", file, ttl, nil)).String()

	return
}

func (s *Storage) ListFiles(ctx context.Context, dir string) (files []string, err error) {
	defer rg.Guard(&err)

	dir = strings.TrimPrefix(dir, "/")
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	var marker string

next:
	res, _ := rg.Must2(s.cos.Bucket.Get(ctx, &cos.BucketGetOptions{
		Prefix:    dir,
		Delimiter: "/",
		MaxKeys:   1000,
		Marker:    marker,
	}))

	for _, c := range res.Contents {
		if strings.HasSuffix(c.Key, "/") {
			continue
		}
		files = append(files, path.Base(c.Key))
	}

	if res.IsTruncated {
		marker = res.NextMarker
		goto next
	}

	return
}
