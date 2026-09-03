package canal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	acConvert "github.com/illidaris/aphrodite/pkg/convert"
	group "github.com/illidaris/aphrodite/pkg/group/v2"
	"github.com/spf13/cast"
)

type MigrateConnector struct {
	MigrateConnectorOption
	Client *sql.DB
}

func (c *MigrateConnector) Id() string {
	return c.DbInstance
}

func (c *MigrateConnector) Stats(ctx context.Context) OperateStats {
	return c.Outer.Stats()
}

func (c *MigrateConnector) Close(ctx context.Context) error {
	_ = c.Outer.Close(ctx)
	return c.Client.Close()
}

func (c *MigrateConnector) Run(ctx context.Context) error {
	_, _ = group.GroupFunc(func(vs ...int) (int64, error) {
		for _, v := range vs {
			time.Sleep(time.Millisecond * 100)
			err := c.run(ctx, v)
			if err != nil {
				return 0, err
			}
		}
		return 0, nil
	}, c.TableArgs)
	return nil

}

func (c *MigrateConnector) run(ctx context.Context, tid int) error {
	field := c.CursorField
	tableName := fmt.Sprintf(c.TableFilter, tid)
	tableName, _ = acConvert.FieldFilter(tableName, acConvert.FieldFilterLevelDefault)
	sqlstr := fmt.Sprintf("SELECT * FROM %s WHERE `id` > ? ORDER BY `id` LIMIT ?", tableName)
	cursor := ConvertCursor(c.CursorPos, c.CursorType)
	for {
		docs, err := Query2Map(ctx, c.Client, sqlstr, cursor, c.Batch)
		if err != nil {
			return err
		}
		if len(docs) == 0 {
			return nil
		}
		entries := []Entry{}
		for _, doc := range docs {
			v, ok := doc[field]
			if !ok {
				return errors.New("当前表没有id主键")
			}
			delete(doc, field)
			cursor = ConvertCursor(v, c.CursorType)
			c.Log.Debug(ctx, "Migrate_Values %s %v %v", tableName, cursor, doc)
			entries = append(entries, Entry{
				Id:    cast.ToString(v),
				Index: c.Index,
				Act:   ActionIndex,
				Doc:   doc,
			})
		}
		ok, err := c.Outer.Sync(ctx, entries...)
		if err != nil {
			c.Log.Error(ctx, "[gocanal]Run_Sync canal entries to Outer: %v", err)
			return err
		}
		if !ok {
			stats := c.Outer.Stats()
			okErr := fmt.Errorf("%v stats %s", cursor, stats.GetMsg(","))
			c.Log.Error(ctx, "[gocanal]Run_SyncFunc_NoOk %s", okErr)
			return okErr
		}
	}
}

func ConvertCursor(raw any, tp int8) any {
	if tp == 1 {
		return cast.ToString(raw)
	}
	return cast.ToInt(raw)
}

func Query2Map(ctx context.Context, db *sql.DB, sqlstr string, args ...any) ([]map[string]any, error) {
	// 准备查询语句
	stmt, err := db.PrepareContext(ctx, sqlstr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close() // 确保查询语句被关闭

	// 执行查询
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // 确保查询结果集被关闭

	// 获取查询返回的列名
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 计算列数，并初始化结果集和存放列值的切片
	var (
		colCount = len(cols)
		result   = []map[string]any{}
		values   = make([]any, colCount)
		valPtrs  = make([]any, colCount)
	)

	// 遍历所有行，并将每行数据转换为map格式
	for rows.Next() {
		for i := 0; i < colCount; i++ {
			valPtrs[i] = &values[i]
		}

		if subErr := rows.Scan(valPtrs...); subErr != nil {
			return nil, subErr
		}

		entry := map[string]any{}
		for i, col := range cols {
			var v any
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b) // 将字节切片转换为字符串
			} else {
				v = val
			}
			entry[col] = v
		}

		result = append(result, entry)
	}

	return result, nil
}
