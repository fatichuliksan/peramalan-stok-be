package awss3

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"peramalan-stok-be/src/helper/viper"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// AwsS3 ...
type AwsS3Helper struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
}

func (ctrl *AwsS3Helper) Connect() *session.Session {
	ctrl.Setting()
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(ctrl.Region),
			Credentials: credentials.NewStaticCredentials(
				ctrl.AccessKeyID,
				ctrl.SecretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})
	if err != nil {
		panic(err)
	}
	return sess
}

func (ctrl *AwsS3Helper) Upload(path string, files []*multipart.FileHeader) (pathFile string, location string, err error) {
	sess := ctrl.Connect()
	uploader := s3manager.NewUploader(sess)

	fileHeader := files[0]
	fileName := fileHeader.Filename
	key := ctrl.RandStringBytes(8) + "_" + fileName

	file, err := fileHeader.Open()
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	//upload to the s3 bucket
	uploaded, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(ctrl.BucketName + path),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		return "", "", err
	}

	fmt.Println("uploaded to AWS S3 : " + uploaded.Location)
	return path + key, uploaded.Location, nil
}

func (ctrl *AwsS3Helper) Upload2(path string, fileName string, file io.Reader) (string, string, error) {
	sess := ctrl.Connect()
	uploader := s3manager.NewUploader(sess)
	now := time.Now()
	key := now.Format("150405") + ctrl.RandStringBytes(5) + "_" + fileName

	//upload to the s3 bucket
	uploaded, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(ctrl.BucketName + path),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return "", "", err
	}

	fmt.Println("uploaded to AWS S3 : " + uploaded.Location)
	return path + key, uploaded.Location, nil
}

func (ctrl *AwsS3Helper) Upload3(path string, files []*multipart.FileHeader) (pathFile string, location string, err error) {
	sess := ctrl.Connect()
	uploader := s3manager.NewUploader(sess)

	fileHeader := files[0]
	now := time.Now()
	key := now.Format("150405") + ctrl.RandStringBytes(6)

	file, err := fileHeader.Open()
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	//upload to the s3 bucket
	uploaded, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(ctrl.BucketName + path),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		return "", "", err
	}

	fmt.Println("uploaded to AWS S3 : " + uploaded.Location)
	return path + key, uploaded.Location, nil
}

func (ctrl *AwsS3Helper) Download(key string) (string, error) {
	sess := ctrl.Connect()

	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(ctrl.BucketName),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(1 * time.Minute)

	if err != nil {
		fmt.Println("failed to get url file from AWS S3 for key : " + key)
		return "", err
	}

	// fmt.Println("url file request : " + urlStr)
	return urlStr, nil
}

func (ctrl *AwsS3Helper) Setting() {
	config := viper.NewViper("config.json", "json")

	region := config.GetString("aws.region")
	accessKeyID := config.GetString("aws.access_key_id")
	secretAccessKeyID := config.GetString("aws.secret_access_key")
	bucketName := config.GetString("aws.bucket_name")

	ctrl.Region = region
	ctrl.AccessKeyID = accessKeyID
	ctrl.SecretAccessKey = secretAccessKeyID
	ctrl.BucketName = bucketName
}

func (ctrl *AwsS3Helper) RandStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (ctrl *AwsS3Helper) UploadByte(path string, binary string, name string) (string, error) {
	sess := ctrl.Connect()
	uploader := s3manager.NewUploader(sess)

	key := ctrl.RandStringBytes(8) + "_" + name

	//upload to the s3 bucket
	uploaded, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(ctrl.BucketName + path),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(binary)),
	})

	if err != nil {
		return "", err
	}

	fmt.Println("uploaded to AWS S3 : " + uploaded.Location)
	return path + key, nil
}
