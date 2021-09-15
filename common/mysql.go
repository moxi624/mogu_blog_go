package common

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/10 2:43 下午
 * @version 1.0
 */
import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB
var err error

func init() {
	mysqladmin, _ := beego.AppConfig.String("mysqladmin")
	mysqlpwd, _ := beego.AppConfig.String("mysqlpwd")
	mysqldb, _ := beego.AppConfig.String("mysqldb")
	mysqlurls, _ := beego.AppConfig.String("mysqlurls")
	dsn := mysqladmin + ":" + mysqlpwd + "@tcp(" + mysqlurls + ")/" + mysqldb + "?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,  // 慢 SQL 阈值
			LogLevel:      logger.Error, // Log level
			Colorful:      true,         // 彩色打印
		},
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		logs.Error(err)
	}

}
