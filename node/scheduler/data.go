package scheduler

import (
	"sync"
	"time"

	"github.com/linguohua/titan/api"
	"github.com/linguohua/titan/node/scheduler/db/persistent"
	"golang.org/x/xerrors"
)

// Data Data
type Data struct {
	nodeManager     *NodeManager
	dataManager     *DataManager
	cid             string
	cacheMap        sync.Map
	reliability     int
	needReliability int
	totalSize       int
	cacheCount      int
	// rootCacheID     string
	totalBlocks int
	nodes       int
	expiredTime time.Time
}

func newData(nodeManager *NodeManager, dataManager *DataManager, cid string, reliability int) *Data {
	return &Data{
		nodeManager:     nodeManager,
		dataManager:     dataManager,
		cid:             cid,
		reliability:     0,
		needReliability: reliability,
		cacheCount:      0,
		totalBlocks:     1,
		// rootCacheID:     "",
	}
}

func loadData(cid string, nodeManager *NodeManager, dataManager *DataManager) *Data {
	dInfo, err := persistent.GetDB().GetDataInfo(cid)
	if err != nil {
		log.Errorf("loadData %s err :%s", cid, err.Error())
		return nil
	}
	if dInfo != nil {
		data := &Data{}
		data.cid = cid
		data.nodeManager = nodeManager
		data.dataManager = dataManager
		data.totalSize = dInfo.TotalSize
		data.needReliability = dInfo.NeedReliability
		data.reliability = dInfo.Reliability
		data.cacheCount = dInfo.CacheCount
		// data.rootCacheID = dInfo.RootCacheID
		data.totalBlocks = dInfo.TotalBlocks
		data.nodes = dInfo.Nodes
		data.expiredTime = dInfo.ExpiredTime

		idList, err := persistent.GetDB().GetCachesWithData(cid)
		if err != nil {
			log.Warnf("loadData GetCacheWithData err:%s", err.Error())
			return data
		}

		for _, cacheID := range idList {
			if cacheID == "" {
				continue
			}
			c := loadCache(cacheID, cid, nodeManager, data)
			if c == nil {
				continue
			}

			data.cacheMap.Store(cacheID, c)
		}

		return data
	}

	return nil
}

func (d *Data) haveRootCache() bool {
	have := false

	d.cacheMap.Range(func(key, value interface{}) bool {
		if have {
			return true
		}

		if value != nil {
			c := value.(*Cache)
			if c != nil {
				have = c.isRootCache && c.status == persistent.CacheStatusSuccess
			}
		}

		return true
	})

	return have
}

func (d *Data) createCache(isRootCache bool) (*Cache, error) {
	cache, err := newCache(d.nodeManager, d, d.cid, isRootCache)
	if err != nil {
		return nil, xerrors.Errorf("new cache err:%s", err.Error())
	}

	return cache, nil
}

func (d *Data) updateAndSaveCacheingInfo(blockInfo *persistent.BlockInfo, info *api.CacheResultInfo, cache *Cache, createBlocks []*persistent.BlockInfo) error {
	if !d.haveRootCache() {
		d.totalSize = cache.totalSize
		d.totalBlocks = cache.totalBlocks
	}

	dInfo := &persistent.DataInfo{
		CID:         d.cid,
		TotalSize:   d.totalSize,
		TotalBlocks: d.totalBlocks,
		Reliability: d.reliability,
		CacheCount:  d.cacheCount,
	}

	cInfo := &persistent.CacheInfo{
		// ID:          cache.dbID,
		CarfileID:   cache.carfileCid,
		CacheID:     cache.cacheID,
		DoneSize:    cache.doneSize,
		Status:      int(cache.status),
		DoneBlocks:  cache.doneBlocks,
		Reliability: cache.reliability,
		TotalSize:   cache.totalSize,
		TotalBlocks: cache.totalBlocks,
		// RemoveBlocks: cache.removeBlocks,
	}

	return persistent.GetDB().SaveCacheingResults(dInfo, cInfo, blockInfo, createBlocks)
}

func (d *Data) updateAndSaveCacheEndInfo(cache *Cache) error {
	if cache.status == persistent.CacheStatusSuccess {
		d.reliability += cache.reliability
		// if !d.haveRootCache() {
		// 	d.rootCacheID = cache.cacheID
		// }
	}

	cNodes, err := persistent.GetDB().GetNodesFromCache(cache.cacheID)
	if err != nil {
		log.Warnf("updateAndSaveCacheEndInfo GetNodesFromCache err:%s", err.Error())
	}

	dNodes, err := persistent.GetDB().GetNodesFromData(d.cid)
	if err != nil {
		log.Warnf("updateAndSaveCacheEndInfo GetNodesFromData err:%s", err.Error())
	}

	d.nodes = dNodes
	dInfo := &persistent.DataInfo{
		CID:         d.cid,
		TotalSize:   d.totalSize,
		TotalBlocks: d.totalBlocks,
		Reliability: d.reliability,
		CacheCount:  d.cacheCount,
		Nodes:       d.nodes,
	}

	cache.nodes = cNodes
	cInfo := &persistent.CacheInfo{
		// ID:          cache.dbID,
		CarfileID:   cache.carfileCid,
		CacheID:     cache.cacheID,
		DoneSize:    cache.doneSize,
		Status:      int(cache.status),
		DoneBlocks:  cache.doneBlocks,
		Reliability: cache.reliability,
		TotalSize:   cache.totalSize,
		TotalBlocks: cache.totalBlocks,
		Nodes:       cache.nodes,
	}

	return persistent.GetDB().SaveCacheEndResults(dInfo, cInfo)
}

func (d *Data) findOrCreateCache(cacheID string) (string, error) {
	var err error
	var cache *Cache
	var list map[string]string

	if cacheID != "" {
		cacheI, ok := d.cacheMap.Load(cacheID)
		if !ok || cacheI == nil {
			err = xerrors.Errorf("Not Found CacheID :%s", cacheID)
			return cacheID, err
		}
		cache = cacheI.(*Cache)

		list, err = persistent.GetDB().GetUndoneBlocks(cacheID)
		if err != nil {
			return cacheID, err
		}
	} else {
		cache, err = d.createCache(!d.haveRootCache())
		if err != nil {
			return "", err
		}

		d.cacheMap.Store(cache.cacheID, cache)
		cacheID = cache.cacheID

		list = map[string]string{cache.carfileCid: ""}
	}

	d.cacheCount++

	return cacheID, cache.startCache(list)
}

func (d *Data) cacheEnd(c *Cache) {
	var err error
	cacheID := ""
	defer func() {
		if err != nil {
			d.dataManager.dataTaskEnd(d.cid, err.Error(), c.cacheID)
		} else {
			d.dataManager.dataTaskStart(d.cid, cacheID)
		}
	}()

	err = d.updateAndSaveCacheEndInfo(c)
	if err != nil {
		err = xerrors.Errorf("saveCacheEndResults err:%s", err.Error())
		return
	}

	if d.cacheCount > d.needReliability {
		err = xerrors.Errorf("cacheCount:%d reach needReliability:%d", d.cacheCount, d.needReliability)
		return
	}

	if d.needReliability <= d.reliability {
		err = xerrors.Errorf("reliability is enough:%d/%d", d.reliability, d.needReliability)
		return
	}

	// old cache
	d.cacheMap.Range(func(key, value interface{}) bool {
		c := value.(*Cache)

		if c.status != persistent.CacheStatusSuccess {
			cacheID = c.cacheID
		}

		return true
	})

	cacheID, err = d.findOrCreateCache(cacheID)
}
