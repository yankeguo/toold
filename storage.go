package toold

import (
	"net/http"
	"net/url"
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
