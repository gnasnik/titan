package carfile

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
	legacy "github.com/ipfs/go-ipld-legacy"
	"github.com/ipfs/go-libipfs/blocks"
	logging "github.com/ipfs/go-log/v2"
	"github.com/ipfs/go-merkledag"
	dagpb "github.com/ipld/go-codec-dagpb"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/linguohua/titan/api"
	"github.com/linguohua/titan/api/types"
	"github.com/linguohua/titan/node/carfile/cache"
	"github.com/linguohua/titan/node/carfile/fetcher"
	"github.com/linguohua/titan/node/carfile/store"
	"github.com/linguohua/titan/node/device"
	mh "github.com/multiformats/go-multihash"
)

var log = logging.Logger("carfile")

const (
	batch = 5
)

type CarfileImpl struct {
	scheduler       api.Scheduler
	device          *device.Device
	cacheMgr        *cache.Manager
	carfileStore    *store.CarfileStore
	TotalBlockCount int
}

func NewCarfileImpl(carfileStore *store.CarfileStore, scheduler api.Scheduler, bFetcher fetcher.BlockFetcher, device *device.Device) *CarfileImpl {
	cfImpl := &CarfileImpl{
		scheduler:    scheduler,
		device:       device,
		carfileStore: carfileStore,
	}

	opts := &cache.ManagerOptions{CarfileStore: carfileStore, BFetcher: bFetcher, DownloadBatch: batch}
	cfImpl.cacheMgr = cache.NewManager(opts)

	totalBlockCount, err := carfileStore.BlockCount()
	if err != nil {
		log.Panicf("NewCarfileImpl block count error:%s", err.Error())
	}
	cfImpl.TotalBlockCount = totalBlockCount

	legacy.RegisterCodec(cid.DagProtobuf, dagpb.Type.PBNode, merkledag.ProtoNodeConverter)
	legacy.RegisterCodec(cid.Raw, basicnode.Prototype.Bytes, merkledag.RawNodeConverter)

	return cfImpl
}

func (cfImpl *CarfileImpl) CacheCarfile(ctx context.Context, rootCID string, dss []*types.DownloadSource) error {
	if types.RunningNodeType == types.NodeEdge && len(dss) == 0 {
		return fmt.Errorf("download source can not empty")
	}

	root, err := cid.Decode(rootCID)
	if err != nil {
		return err
	}

	if has := cfImpl.carfileStore.HasCarfile(root); has {
		log.Debugf("CacheCarfile %s already exist", root.String())
		return nil
	}

	log.Debugf("CacheCarfile cid:%s", rootCID)

	cfImpl.cacheMgr.AddToWaitList(root, dss)
	return nil
}

func (cfImpl *CarfileImpl) removeResult() error {
	if count, err := cfImpl.carfileStore.BlockCount(); err == nil {
		cfImpl.TotalBlockCount = count
	}

	_, diskUsage := cfImpl.device.GetDiskUsageStat()
	ret := types.RemoveAssetResult{BlocksCount: cfImpl.TotalBlockCount, DiskUsage: diskUsage}

	return cfImpl.scheduler.RemoveAssetResult(context.Background(), ret)
}

func (cfImpl *CarfileImpl) deleteCarfile(ctx context.Context, c cid.Cid) error {
	ok, err := cfImpl.cacheMgr.DeleteCarFromWaitList(c)
	if err != nil {
		log.Errorf("delete car %s from wait list error:%s", c.String(), err.Error())
		return err
	}

	if ok {
		log.Debugf("delete carfile %s from wait list", c.String())
		return cfImpl.removeResult()
	}

	log.Debugf("delete carfile %s", c.String())

	if err := cfImpl.carfileStore.DeleteCarfile(c); err != nil {
		return err
	}
	return cfImpl.removeResult()
}

func (cfImpl *CarfileImpl) DeleteCarfile(ctx context.Context, carfileCID string) error {
	c, err := cid.Decode(carfileCID)
	if err != nil {
		log.Errorf("DeleteCarfile, decode carfile cid %s error :%s", carfileCID, err.Error())
		return err
	}

	go cfImpl.deleteCarfile(ctx, c) //nolint:errcheck

	return nil
}

func (cfImpl *CarfileImpl) GetBlock(ctx context.Context, cidStr string) ([]byte, error) {
	c, err := cid.Decode(cidStr)
	if err != nil {
		return nil, err
	}

	blk, err := cfImpl.carfileStore.Block(c)
	if err != nil {
		return nil, err
	}
	return blk.RawData(), nil
}

func (cfImpl *CarfileImpl) QueryCacheStat(ctx context.Context) (*types.CacheStat, error) {
	carfileCount, err := cfImpl.carfileStore.CarfileCount()
	if err != nil {
		return nil, err
	}

	cacheStat := &types.CacheStat{}
	cacheStat.TotalBlockCount = cfImpl.TotalBlockCount
	cacheStat.TotalAssetCount = carfileCount
	cacheStat.WaitCacheAssetCount = cfImpl.cacheMgr.WaitListLen()
	_, cacheStat.DiskUsage = cfImpl.device.GetDiskUsageStat()

	carfileCache := cfImpl.cacheMgr.CachingCar()
	if carfileCache != nil {
		cacheStat.CachingAssetCID = carfileCache.Root().String()
	}

	log.Debugf("cacheStat:%#v", *cacheStat)

	return cacheStat, nil
}

func (cfImpl *CarfileImpl) QueryCachingCarfile(ctx context.Context) (*types.CachingAsset, error) {
	carfileCache := cfImpl.cacheMgr.CachingCar()
	if carfileCache == nil {
		return nil, fmt.Errorf("caching carfile not exist")
	}

	ret := &types.CachingAsset{}
	ret.CID = carfileCache.Root().Hash().String()
	ret.TotalSize = carfileCache.TotalSize()
	ret.DoneSize = carfileCache.DoneSize()

	return ret, nil
}

func (cfImpl *CarfileImpl) GetBlocksOfCarfile(carfileCID string, indices []int) (map[int]string, error) {
	c, err := cid.Decode(carfileCID)
	if err != nil {
		return nil, err
	}

	cids, err := cfImpl.carfileStore.BlocksOfCarfile(c)
	if err != nil {
		return nil, err
	}

	ret := make(map[int]string, len(indices))
	for _, index := range indices {
		if index >= len(cids) {
			return nil, fmt.Errorf("index is out of blocks count")
		}

		c := cids[index]
		ret[index] = c.String()
	}
	return ret, nil
}

func (cfImpl *CarfileImpl) BlockCountOfCarfile(carfileCID string) (int, error) {
	c, err := cid.Decode(carfileCID)
	if err != nil {
		return 0, err
	}

	return cfImpl.carfileStore.BlockCountOfCarfile(c)
}

func (cfImpl *CarfileImpl) CacheCarForSyncData(carfiles []string) error {
	switch types.RunningNodeType {
	case types.NodeCandidate:
		for _, hash := range carfiles {
			multihash, err := mh.FromHexString(hash)
			if err != nil {
				return err
			}

			cid := cid.NewCidV1(cid.Raw, multihash)
			if err := cfImpl.CacheCarfile(context.Background(), cid.String(), nil); err != nil {
				return err
			}
		}
	case types.NodeEdge:
		return fmt.Errorf("not implement")
	default:
		return fmt.Errorf("unsupport node type:%s", types.RunningNodeType)
	}

	return nil
}

func (cfImpl *CarfileImpl) CachedProgresses(ctx context.Context, carfileCIDs []string) (*types.PullResult, error) {
	progresses := make([]*types.AssetPullProgress, 0, len(carfileCIDs))
	for _, carfileCID := range carfileCIDs {
		root, err := cid.Decode(carfileCID)
		if err != nil {
			return nil, err
		}

		progress, err := cfImpl.carProgress(root)
		if err != nil {
			return nil, err
		}

		progresses = append(progresses, progress)
	}

	result := &types.PullResult{
		Progresses:       progresses,
		TotalBlocksCount: cfImpl.TotalBlockCount,
	}

	if count, err := cfImpl.carfileStore.CarfileCount(); err == nil {
		result.AssetCount = count
	}
	_, result.DiskUsage = cfImpl.device.GetDiskUsageStat()

	return result, nil
}

func (cfImpl *CarfileImpl) progressForSucceededCar(root cid.Cid) (*types.AssetPullProgress, error) {
	progress := &types.AssetPullProgress{
		CID:    root.String(),
		Status: types.ReplicaStatusSucceeded,
	}

	if count, err := cfImpl.carfileStore.BlockCountOfCarfile(root); err == nil {
		progress.BlocksCount = count
		progress.DoneBlocksCount = count
	}

	blk, err := cfImpl.carfileStore.Block(root)
	if err != nil {
		return nil, err
	}

	blk = blocks.NewBlock(blk.RawData())
	node, err := legacy.DecodeNode(context.Background(), blk)
	if err != nil {
		return nil, err
	}

	linksSize := uint64(len(blk.RawData()))
	for _, link := range node.Links() {
		linksSize += link.Size
	}

	progress.Size = int64(linksSize)
	progress.DoneSize = int64(linksSize)

	return progress, nil
}

func (cfImpl *CarfileImpl) carProgress(root cid.Cid) (*types.AssetPullProgress, error) {
	status, err := cfImpl.cacheMgr.CachedStatus(root)
	if err != nil {
		return nil, err
	}

	switch status {
	case types.ReplicaStatusWaiting:
		return &types.AssetPullProgress{CID: root.String(), Status: types.ReplicaStatusWaiting}, nil
	case types.ReplicaStatusPulling:
		return cfImpl.cacheMgr.CachingCar().Progress(), nil
	case types.ReplicaStatusFailed:
		return cfImpl.cacheMgr.ProgressForFailedCar(root)
	case types.ReplicaStatusSucceeded:
		return cfImpl.progressForSucceededCar(root)
	}
	return nil, nil
}
