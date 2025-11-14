package initialize

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/fans-club/model"
	"go.uber.org/zap"
)

// Gorm 初始化数据库表
func Gorm(ctx context.Context) {
	err := global.GVA_DB.AutoMigrate(
		&model.FansClub{},
		&model.FansClubMember{},
		&model.FansClubPost{},
	)
	if err != nil {
		global.GVA_LOG.Error("fans-club plugin: register table failed", zap.Error(err))
	} else {
		global.GVA_LOG.Info("fans-club plugin: register table success")
	}
}
