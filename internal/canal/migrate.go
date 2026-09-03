package canal

import (
	"context"
	"gocanal/config"
	"gocanal/pkg/canal"
	"gocanal/pkg/canal/outer"
	"gocanal/pkg/log"

	"github.com/spf13/viper"
)

func Migrate(ctx context.Context) {
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

	l := log.NewRestLogger()
	canal.SetLogger(l)
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

	for _, migrateCfg := range cfg.Migrates {
		sc, err := canal.NewMigrateConnector(
			canal.WithMigrateDbIp(migrateCfg.DbIp),
			canal.WithMigrateDbPort(migrateCfg.DbPort),
			canal.WithMigrateDbUser(migrateCfg.DbUser),
			canal.WithMigrateDbPwd(migrateCfg.DbPwd),
			canal.WithMigrateDbName(migrateCfg.DbName),
			canal.WithMigrateDbInstance(migrateCfg.Instance),
			canal.WithMigrateIndex(migrateCfg.Index),
			canal.WithMigrateTableFilter(migrateCfg.Filter),
			canal.WithMigrateTableArgs(migrateCfg.TableArgs),
			canal.WithMigrateBatch(migrateCfg.Batch),
			canal.WithMigrateTimeout(migrateCfg.Timeout),
			canal.WithMigrateOuter(esOuter),
			canal.WithMigrateLogger(l),
		)
		if err != nil {
			log.Error(ctx, "NewMigrateConnector: %v", err)
			return
		}
		err = sc.Run(ctx)
		if err != nil {
			log.Error(ctx, "NewMigrateConnector Run: %v", err)
			return
		}
	}

}
