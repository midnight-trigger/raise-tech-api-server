package infra

import (
	"github.com/midnight-trigger/raise-tech-api-server/infra/mysql"
	"github.com/midnight-trigger/raise-tech-api-server/infra/s3"
)

func Init() {
	mysql.Init()
	s3.Init()
}
