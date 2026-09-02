package canal

import (
	"context"
	"fmt"
	"sync/atomic"

	pbe "github.com/withlin/canal-go/protocol/entry"
)

type ColsToKVsHandle func(columns []*pbe.Column) (string, map[string]any)

func DefaultColsToDoc(columns []*pbe.Column) (string, map[string]any) {
	doc := make(map[string]any, len(columns))
	var id string
	for _, column := range columns {
		value := column.GetValue()
		doc[column.GetName()] = value
		if column.GetIsKey() && id == "" {
			id = value
		}
	}
	return id, doc
}

var _ = ColsToKVsHandle(DefaultColsToDoc)

type KVsCbHandle func(context.Context, string, string, map[string]string)

type ValueChangeCbHandle func(context.Context, string, []*pbe.Column, []*pbe.Column)

type IOuter interface {
	Stats() OperateStats
	Close(context.Context) error
	Sync(ctx context.Context, entries ...Entry) (bool, error)
	Check(ctx context.Context, key string) error
	SyncStruct(ctx context.Context, key, index, mapping string) error
}

type BaseOuter struct {
	TableMap map[string]string
	Log      ILogger
}

var _ = IOuter(&StdOuter{})

type StdOuter struct {
	BaseOuter
	Counter atomic.Uint64
}

func (i *StdOuter) Stats() OperateStats {
	total := i.Counter.Load()
	return OperateStats{
		NumAdded:    total,
		NumFlushed:  total,
		NumRequests: total,
	}
}

func (i *StdOuter) Close(context.Context) error {
	return nil
}

func (i *StdOuter) Sync(ctx context.Context, entries ...Entry) (bool, error) {
	for _, v := range entries {
		i.Counter.Add(1)
		fmt.Printf("[%s]%s %v", v.Act, v.Id, v.Doc)
	}
	return true, nil
}

func (i *StdOuter) Check(ctx context.Context, key string) error {
	return nil
}

func (i *StdOuter) SyncStruct(ctx context.Context, key, index, mapping string) error {
	return nil
}
