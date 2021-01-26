package infra

import (
	"github.com/midnight-trigger/raise-tech-api-server/infra/mysql"
)

func Init() {
	mysql.Init()
}
