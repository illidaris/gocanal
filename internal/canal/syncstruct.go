package canal

import (
	"context"
	"fmt"
	"gocanal/config"
	"gocanal/pkg/log"

	"github.com/illidaris/aphrodite/pkg/canal/outer"

	"github.com/spf13/viper"
)

func SyncStruct(ctx context.Context, args ...string) {
	if len(args) == 0 {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			log.Error(ctx, "%v", r)
		}
	}()
	cfg := &config.Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		log.Error(ctx, "Failed to unmarshal config: %v", err)
		return
	}

	esOuter, err := outer.NewElasticOuter(
		outer.WithName(cfg.EsName),
		outer.WithEsUrls(cfg.EsUrls...),
		outer.WithESUser(cfg.ESUser),
		outer.WithESPwd(cfg.ESPwd),
	)
	if err != nil {
		log.Error(ctx, "NewElasticOuter: %v", err)
		return
	}
	defer esOuter.Close(ctx)
	for _, syncCfg := range cfg.Syncs {
		println(fmt.Sprintf("同步结构：%s", syncCfg.Instance))
		syncErr := esOuter.SyncStruct(ctx, cfg.EsName, syncCfg.Index, syncCfg.Mapping)
		println(fmt.Sprintf("%s 执行完毕 %v", syncCfg.Instance, syncErr))
	}
}
