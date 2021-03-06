package mysql

import (
	"database/sql"

	"flashcat.cloud/categraf/inputs"
	"flashcat.cloud/categraf/pkg/tagx"
	"github.com/toolkits/pkg/container/list"
)

func (m *MySQL) gatherEngineInnodbStatusCompute(slist *list.SafeList, ins *Instance, db *sql.DB, globalTags map[string]string, cache map[string]float64) {
	tags := tagx.Copy(globalTags)

	pageUsed := cache["innodb_buffer_pool_pages_total"] - cache["innodb_buffer_pool_pages_free"]
	byteUsed := pageUsed * cache["innodb_page_size"]
	byteData := cache["innodb_buffer_pool_pages_data"] * cache["innodb_page_size"]
	byteDirty := cache["innodb_buffer_pool_pages_dirty"] * cache["innodb_page_size"]
	byteFree := cache["innodb_buffer_pool_pages_free"] * cache["innodb_page_size"]
	byteTotal := cache["innodb_buffer_pool_pages_total"] * cache["innodb_page_size"]
	pageUtil := float64(0)
	if cache["innodb_buffer_pool_pages_total"] != 0 {
		pageUtil = pageUsed / cache["innodb_buffer_pool_pages_total"] * 100
	}

	slist.PushFront(inputs.NewSample("global_status_buffer_pool_bytes", byteUsed, tags, map[string]string{"state": "used"}))
	slist.PushFront(inputs.NewSample("global_status_buffer_pool_bytes", byteData, tags, map[string]string{"state": "data"}))
	slist.PushFront(inputs.NewSample("global_status_buffer_pool_bytes", byteFree, tags, map[string]string{"state": "free"}))
	slist.PushFront(inputs.NewSample("global_status_buffer_pool_bytes", byteTotal, tags, map[string]string{"state": "total"}))
	slist.PushFront(inputs.NewSample("global_status_buffer_pool_bytes", byteDirty, tags, map[string]string{"state": "dirty"}))
	slist.PushFront(inputs.NewSample("global_status_buffer_pool_pages_utilization", pageUtil, tags))

	if ins.ExtraInnodbMetrics {
		slist.PushFront(inputs.NewSample("global_status_buffer_pool_pages", pageUsed, tags, map[string]string{"state": "used"}))
	}
}
