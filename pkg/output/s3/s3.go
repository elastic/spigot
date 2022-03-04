package s3

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/elastic/go-ucfg"
	"github.com/leehinman/spigot/pkg/output"
)

var (
	doOnce   sync.Once
	uploader *manager.Uploader
)

type S3Output struct {
	delimiter string
	bucket    string
	key       string
	buf       *bytes.Buffer
	gw        *gzip.Writer
}

func init() {
	output.Register("s3", New)
}

func New(cfg *ucfg.Config) (s output.Output, err error) {
	c := defaultConfig()
	if err := cfg.Unpack(&c); err != nil {
		return nil, err
	}

	doOnce.Do(func() {
		cfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion(c.Region))
		if err != nil {
			panic(err)
		}
		uploader = manager.NewUploader(s3.NewFromConfig(cfg))
	})
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	key := fmt.Sprintf("%s_%19d_%3d.gz", c.Prefix, time.Now().UnixNano(), rand.Intn(1000))

	s = &S3Output{
		delimiter: c.Delimiter,
		bucket:    c.Bucket,
		key:       key,
		buf:       &buf,
		gw:        gw,
	}
	return s, nil
}

func (s *S3Output) Write(b []byte) (n int, err error) {
	j, err := s.gw.Write(b)
	if err != nil {
		return j, err
	}
	k, err := s.gw.Write([]byte(s.delimiter))
	return j + k, err
}

func (s *S3Output) Close() error {
	s.gw.Close()

	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.key),
		Body:   s.buf,
	})
	return err
}
