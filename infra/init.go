package infra

import (
	"github.com/midnight-trigger/todo/infra/mysql"
)

func Init() {
	mysql.Init()
}
