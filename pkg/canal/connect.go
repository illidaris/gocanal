package canal

import (
	"context"
	"fmt"
	"time"

	"github.com/withlin/canal-go/client"
)

type ISyncConnector interface {
	Id() string
	Stats(ctx context.Context) OperateStats
	Close(ctx context.Context) error
	Run(ctx context.Context) error
}

type SyncConnector struct {
	SyncConnectorOption
	CanalConnector *client.SimpleCanalConnector
}

func (c *SyncConnector) Id() string {
	return c.CanalInstance
}

func (c *SyncConnector) Stats(ctx context.Context) OperateStats {
	return c.Outer.Stats()
}

func (c *SyncConnector) Close(ctx context.Context) error {
	_ = c.Outer.Close(ctx)
	return c.CanalConnector.DisConnection()
}

func (c *SyncConnector) Run(ctx context.Context) error {
	err := c.CanalConnector.Connect()
	if err != nil {
		return err
	}
	defer c.CanalConnector.DisConnection()
	// 订阅表
	err = c.CanalConnector.Subscribe(c.TableFilter)
	if err != nil {
		return err
	}

	if c.Outer == nil {
		return fmt.Errorf("outer is nil")
	}
	defer func() {
		if err := c.Outer.Close(ctx); err != nil {
			c.Logger.Error(ctx, "[gocanal]Run_OuterClose %v", err)
		}
	}()

	for {
		if err := ctx.Err(); err != nil {
			return nil
		}
		unit := int32(2)
		timout := int64(c.Timeout.Seconds())
		message, err := c.CanalConnector.GetWithOutAck(c.Batch, &timout, &unit)
		if err != nil {
			return err
		}
		batchId := message.Id
		// 批次Id -1 或者 消息集 为空， 则休眠
		if batchId == -1 || len(message.Entries) <= 0 {
			time.Sleep(300 * time.Millisecond)
			continue
		}
		ok, err := c.Outer.Sync(ctx, c.Index, DefaultColsToDoc, message.Entries...)
		if err != nil {
			c.Logger.Error(ctx, "[gocanal]Run_SyncFunc canal entries to Outer: %v", err)
			return err
		}
		if !ok {
			stats := c.Outer.Stats()
			okErr := fmt.Errorf("%d stats %s", batchId, stats.GetMsg(","))
			c.Logger.Error(ctx, "[gocanal]Run_SyncFunc_NoOk %s", okErr)
			return okErr
		}
		if err := c.CanalConnector.Ack(batchId); err != nil {
			c.Logger.Error(ctx, "[gocanal]Run_Ack batch_%d %v", batchId, err)
			return err
		}
	}
}

// statsInterval := c.Option.StatsInterval
// 	if statsInterval <= 0 {
// 		statsInterval = 10 * time.Second
// 	}
// 	ticker := time.NewTicker(statsInterval)
// 	defer ticker.Stop()
// 	go func() {
// 		for {
// 			select {
// 			case <-ticker.C:
// 				s := bIdx.Stats()
// 				resp := convertStatsToOperateStats(s)
// 				if resp.NumAdded > 0 {
// 					c.Option.RobotNotify(ctx, resp.GetTitle(), resp.GetMsg("\n"))
// 					c.Option.Logger.Info(ctx, resp.GetMsg(","))
// 				}
// 			case <-ctx.Done():
// 				return
// 			}
// 		}
// 	}()
