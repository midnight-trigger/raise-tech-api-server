package s3

import (
	"fmt"
	"mime/multipart"
	"os"

	"github.com/midnight-trigger/raise-tech-api-server/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/midnight-trigger/raise-tech-api-server/configs"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var uploader *s3manager.Uploader
var svc *s3.S3

func Init() {

	var sess *session.Session
	if configs.IsLocal() {
		sess = session.Must(session.NewSession(&aws.Config{
			Endpoint:         aws.String(viper.GetString("aws.endpoint")),
			S3ForcePathStyle: aws.Bool(true),
		}))
	} else {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Credentials: credentials.NewStaticCredentialsFromCreds(credentials.Value{
					AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
					SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
				}),
				Region: aws.String(os.Getenv("AWS_REGION")),
			},
		}))
	}
	uploader = s3manager.NewUploader(sess)
	svc = s3.New(sess)
}

func UploadImageToS3(file multipart.File, format string, bucket string) (fileName string, err error) {
	fileName = fmt.Sprintf("%s.%s", uuid.New(), format)

	upParams := &s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fileName),
		Body:        file,
		ACL:         aws.String("public-read"),
		ContentType: aws.String("image/" + format),
	}
	fmt.Println(fileName)

	_, err = uploader.Upload(upParams)
	if err != nil {
		err = errors.New(err.Error())
		logger.L.Error(err)
	}

	return
}

func DeleteImage(fileName string) (err error) {
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(viper.GetString("image.bucket")),
		Key:    aws.String(fileName),
	})
	if err != nil {
		err = errors.New(err.Error())
		logger.L.Error(err)
		return
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(viper.GetString("image.bucket")),
		Key:    aws.String(fileName),
	})
	if err != nil {
		err = errors.New(err.Error())
		logger.L.Error(err)
	}

	return
}
