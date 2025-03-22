package helper

import (
	"peramalan-stok-be/src/helper/awss3"
	"peramalan-stok-be/src/helper/response"
	"peramalan-stok-be/src/helper/validator"

	viperHelper "peramalan-stok-be/src/helper/viper"
	"time"
)

type Helper struct {
	AwsS3        awss3.AwsS3Helper
	Response     response.Interface
	Config       viperHelper.Interface
	Validator    validator.ValidatorHelper
	TimeLocation *time.Location
}
