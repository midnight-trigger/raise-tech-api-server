package domain

import (
	"errors"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/labstack/echo"
	"github.com/midnight-trigger/raise-tech-api-server/api/error_handling"
	"github.com/midnight-trigger/raise-tech-api-server/infra/s3"
	"github.com/midnight-trigger/raise-tech-api-server/logger"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

type Image struct {
	Base
}

func GetNewImageService() *Image {
	image := new(Image)
	return image
}

func (s *Image) PostImage(fileName string) (r Result) {
	r.New()

	return
}

func GetImageFileName(ctx echo.Context) (fileName string, err error) {
	r := ctx.Request()
	r.Body = http.MaxBytesReader(ctx.Response(), r.Body, viper.GetInt64("image.maxSize"))

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		if err.Error() == "http: request body too large" {
			limit := strconv.FormatInt(viper.GetInt64("image.maxSize")/(1024*1000), 10)
			err = errors.New(error_handling.GetValidationErrorMessage("", "imageSize", limit))
		} else {
			err = errors.New(err.Error())
		}
		logger.L.Error(err)
		return
	}
	defer file.Close()

	imgSrc, err := imaging.Decode(file, imaging.AutoOrientation(true))
	format := strings.Split(fileHeader.Filename, ".")[1]
	if err != nil {
		err = errors.New(error_handling.GetValidationErrorMessage("", "imageType", ""))
		logger.L.Error(err)
		return
	}

	bucket := viper.GetString("image.Bucket")
	var f multipart.File
	dst := "image." + format

	rctSrc := imgSrc.Bounds()
	width := rctSrc.Dx()
	height := rctSrc.Dy()
	maxWidth := viper.GetInt("image.maxWidth")
	maxHeight := viper.GetInt("image.maxHeight")

	if width > maxWidth || height > maxHeight {
		widthRatio := decimal.NewFromFloat(float64(maxWidth)).Div(decimal.NewFromFloat(float64(width)))
		heightRatio := decimal.NewFromFloat(float64(maxHeight)).Div(decimal.NewFromFloat(float64(height)))

		var ratio decimal.Decimal
		wr, _ := widthRatio.Float64()
		hr, _ := heightRatio.Float64()
		if wr < hr {
			ratio = widthRatio
		} else {
			ratio = heightRatio
		}

		newWidth := int(ratio.Mul(decimal.NewFromFloat(float64(width))).IntPart())
		newHeight := int(ratio.Mul(decimal.NewFromFloat(float64(height))).IntPart())

		imgSrc = imaging.Resize(imgSrc, newWidth, newHeight, imaging.Lanczos)
		err = imaging.Save(imgSrc, dst)
		if err != nil {
			err = errors.New(err.Error())
			logger.L.Error(err)
			return
		}
		f, err = os.Open(dst)
		if err != nil {
			err = errors.New(err.Error())
			logger.L.Error(err)
			return
		}
	} else {
		f, _, err = r.FormFile("image")
		if err != nil {
			err = errors.New(err.Error())
			logger.L.Error(err)
			return
		}
	}

	fileName, err = s3.UploadImageToS3(f, format, bucket)

	defer func() {
		f.Close()
		os.Remove(dst)
	}()

	return
}
