package canal

import (
	"context"
	"fmt"
	"sync/atomic"

	pbe "github.com/withlin/canal-go/protocol/entry"
	"google.golang.org/protobuf/proto"
)

type ColsToKVsHandle func(columns []*pbe.Column) (string, map[string]string)

func DefaultColsToDoc(columns []*pbe.Column) (string, map[string]string) {
	doc := make(map[string]string, len(columns))
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
	Sync(context.Context, string, ...pbe.Entry) (bool, error)
}

type BaseOuter struct {
	TableMap      map[string]string
	ColsToKVs     ColsToKVsHandle
	KVsCb         KVsCbHandle
	ValueChangeCb ValueChangeCbHandle
}

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

func (i *StdOuter) Sync(_ context.Context, index string, entries ...pbe.Entry) (bool, error) {
	for k := range entries {
		entry := &entries[k]
		change := new(pbe.RowChange)
		if err := proto.Unmarshal(entry.GetStoreValue(), change); err != nil {
			return false, fmt.Errorf("decode row change: %w", err)
		}
		for _, row := range change.GetRowDatas() {
			columns := row.GetAfterColumns()
			action := ActionIndex
			if change.GetEventType() == pbe.EventType_DELETE {
				columns = row.GetBeforeColumns()
				action = ActionDelete
			}
			id, doc := i.ColsToKVs(columns)
			i.Counter.Add(1)
			fmt.Printf("[%s]%s %v", action, id, doc)
		}
	}
	return true, nil
}
