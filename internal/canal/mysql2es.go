package canal

import (
	"context"
	"gocanal/config"
	"gocanal/pkg/canal"
	"gocanal/pkg/canal/outer"

	"gocanal/pkg/log"

	"github.com/spf13/viper"
)

func Go(raw context.Context) {
	go func(ctx context.Context) {
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
			log.Error(ctx, "Failed to unmarshal config: %v", err)
			return
		}
		for _, syncCfg := range cfg.Syncs {
			sc, err := canal.NewSyncConnector(
				canal.WithSyncCanalIp(cfg.CanalIp),
				canal.WithSyncCanalPort(cfg.CanalPort),
				canal.WithSyncCanalUser(cfg.CanalUser),
				canal.WithSyncCanalPwd(cfg.CanalPwd),
				canal.WithSyncCanalInstance(syncCfg.Instance),
				canal.WithSyncIndex(syncCfg.Index),
				canal.WithSyncTableFilter(syncCfg.Filter),
				canal.WithSyncBatch(syncCfg.Batch),
				canal.WithSyncTimeout(syncCfg.Timeout),
				canal.WithSyncOuter(esOuter),
				canal.WithSyncLogger(log.NewRestLogger()),
			)
			if err != nil {
				log.Error(ctx, "NewSyncConnector: %v", err)
				return
			}
			err = sc.Run(ctx)
			if err != nil {
				log.Error(ctx, "NewSyncConnector Run: %v", err)
				return
			}
		}
	}(raw)

}
