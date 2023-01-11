package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/fatih/structs"
	"github.com/linguohua/titan/api"

	"github.com/go-redis/redis/v8"
	redigo "github.com/gomodule/redigo/redis"
)

const (
	// redisKeyWaitingDataTaskList  server name
	redisKeyWaitingDataTaskList = "Titan:WaitingDataTaskList:%s"
	// redisKeyValidatorList server name
	redisKeyValidatorList = "Titan:ValidatorList:%s"
	// redisKeyValidateRoundID server name
	redisKeyValidateRoundID = "Titan:ValidateRoundID:%s"
	// redisKeyVerifyingList server name
	redisKeyVerifyingList = "Titan:VerifyingList:%s"
	// redisKeyNodeInfo  deviceID
	redisKeyNodeInfo = "Titan:NodeInfo:%s"
	// redisKeyBlockDownloadRecord serial number
	redisKeyBlockDownloadRecord = "Titan:BlockDownloadRecord:%d"
	// redisKeyBlockDownloadSN
	redisKeyBlockDownloadSN       = "Titan:BlockDownloadRecordSN"
	redisKeyCarfileLatestDownload = "Titan:LatestDownload:%s"

	// redisKeyCacheingNodeList  server name
	redisKeyCacheingNodeList = "Titan:CacheingNodeList:%s"
	// redisKeyNodeCacheList  server name:deviceID
	redisKeyNodeCacheList = "Titan:NodeCacheList:%s:%s"
	// redisKeyCacheingNode  server name:deviceID
	redisKeyCacheingNode = "Titan:CacheingNode:%s:%s"

	// redisKeyBaseInfo  server name
	redisKeyBaseInfo = "Titan:BaseInfo:%s"
)

const cacheErrorExpiration = 72 //hour

// TypeRedis redis
func TypeRedis() string {
	return "Redis"
}

type redisDB struct {
	cli *redis.Client
}

// InitRedis init redis pool
func InitRedis(url string) (DB, error) {
	// fmt.Printf("redis init url:%v", url)

	redisDB := &redisDB{redis.NewClient(&redis.Options{
		Addr:      url,
		Dialer:    nil,
		OnConnect: nil,
	})}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := redisDB.cli.Ping(ctx).Result()

	return redisDB, err
}

func (rd redisDB) SetCacheStart(hash, deviceID string, timeout int64) error {
	CacheingNodeList := fmt.Sprintf(redisKeyCacheingNodeList, serverName)
	nodeCacheList := fmt.Sprintf(redisKeyNodeCacheList, serverName, deviceID)
	nodeKey := fmt.Sprintf(redisKeyCacheingNode, serverName, deviceID)

	ctx := context.Background()
	_, err := rd.cli.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		pipeliner.SAdd(context.Background(), CacheingNodeList, deviceID)
		pipeliner.SAdd(context.Background(), nodeCacheList, hash)
		// Expire
		pipeliner.Set(context.Background(), nodeKey, hash, time.Second*time.Duration(timeout))

		return nil
	})

	return err
}

func (rd redisDB) SetCacheEnd(hash, deviceID string) error {
	cacheingNodeList := fmt.Sprintf(redisKeyCacheingNodeList, serverName)
	nodeCacheList := fmt.Sprintf(redisKeyNodeCacheList, serverName, deviceID)
	nodeKey := fmt.Sprintf(redisKeyCacheingNode, serverName, deviceID)

	ctx := context.Background()
	_, err := rd.cli.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {

		pipeliner.SRem(context.Background(), nodeCacheList, hash)
		// Expire
		pipeliner.Del(context.Background(), nodeKey)

		//TODO
		// if isNodeDoneAll {
		pipeliner.SRem(context.Background(), cacheingNodeList, deviceID)
		// }

		return nil
	})

	return err
}

func (rd redisDB) UpdateNodeCacheingExpireTime(hash, deviceID string, timeout int64) error {
	nodeKey := fmt.Sprintf(redisKeyCacheingNode, serverName, deviceID)
	// Expire
	_, err := rd.cli.Set(context.Background(), nodeKey, hash, time.Second*time.Duration(timeout)).Result()
	return err
}

func (rd redisDB) GetCacheingNodes() ([]string, error) {
	cacheingNodeList := fmt.Sprintf(redisKeyCacheingNodeList, serverName)

	values, err := rd.cli.SMembers(context.Background(), cacheingNodeList).Result()
	if err != nil || values == nil {
		return nil, err
	}

	return values, nil
}

func (rd redisDB) GetNodeCaches(deviceID string) ([]string, error) {
	nodeCacheList := fmt.Sprintf(redisKeyNodeCacheList, serverName, deviceID)

	values, err := rd.cli.SMembers(context.Background(), nodeCacheList).Result()
	if err != nil || values == nil {
		return nil, err
	}

	return values, nil
}

func (rd redisDB) GetCacheingCarfiles() (map[string]int, error) {
	carfileMap := make(map[string]int)

	ctx := context.Background()
	_, err := rd.cli.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		cacheingNodeList := fmt.Sprintf(redisKeyCacheingNodeList, serverName)
		nodes, err := pipeliner.SMembers(context.Background(), cacheingNodeList).Result()
		if err != nil || nodes == nil {
			return nil
		}

		for _, deviceID := range nodes {
			nodeCacheList := fmt.Sprintf(redisKeyNodeCacheList, serverName, deviceID)

			carfiles, err := pipeliner.SMembers(context.Background(), nodeCacheList).Result()
			if err != nil || carfiles == nil {
				continue
			}

			for _, carfile := range carfiles {
				carfileMap[carfile]++
			}
		}

		return nil
	})

	return carfileMap, err
}

func (rd redisDB) GetNodeCacheingCarfile(deviceID string) (string, error) {
	nodeKey := fmt.Sprintf(redisKeyCacheingNode, serverName, deviceID)
	return rd.cli.Get(context.Background(), nodeKey).Result()
}

func (rd redisDB) GetCacheTimeoutNodes() (map[string][]string, error) {
	timeoutMap := make(map[string][]string, 0)

	ctx := context.Background()
	_, err := rd.cli.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		cacheingNodeList := fmt.Sprintf(redisKeyCacheingNodeList, serverName)
		nodes, err := pipeliner.SMembers(context.Background(), cacheingNodeList).Result()
		if err != nil || nodes == nil {
			return nil
		}

		for _, deviceID := range nodes {
			nodeKey := fmt.Sprintf(redisKeyCacheingNode, serverName, deviceID)

			exist, err := pipeliner.Exists(context.Background(), nodeKey).Result()
			if err != nil {
				continue
			}

			if exist == 1 {
				continue
			}

			nodeCacheList := fmt.Sprintf(redisKeyNodeCacheList, serverName, deviceID)
			carfiles, err := pipeliner.SMembers(context.Background(), nodeCacheList).Result()
			if err != nil || carfiles == nil {
				continue
			}

			for _, carfile := range carfiles {
				list := timeoutMap[carfile]
				if list == nil {
					list = make([]string, 0)
				}

				timeoutMap[carfile] = append(list, deviceID)
			}
		}

		return nil
	})

	return timeoutMap, err
}

// waiting data list
func (rd redisDB) SetWaitingDataTask(info *api.CarfileRecordInfo) error {
	key := fmt.Sprintf(redisKeyWaitingDataTaskList, serverName)

	bytes, err := json.Marshal(info)
	if err != nil {
		return err
	}

	_, err = rd.cli.RPush(context.Background(), key, bytes).Result()
	return err
}

func (rd redisDB) GetWaitingDataTask(index int64) (*api.CarfileRecordInfo, error) {
	key := fmt.Sprintf(redisKeyWaitingDataTaskList, serverName)

	value, err := rd.cli.LIndex(context.Background(), key, index).Result()

	if value == "" {
		return nil, redis.Nil
	}

	var info api.CarfileRecordInfo
	bytes, err := redigo.Bytes(value, nil)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

func (rd redisDB) RemoveWaitingDataTasks(infos []*api.CarfileRecordInfo) error {
	key := fmt.Sprintf(redisKeyWaitingDataTaskList, serverName)

	ctx := context.Background()
	_, err := rd.cli.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		for _, info := range infos {
			bytes, err := json.Marshal(info)
			if err != nil {
				continue
			}

			pipeliner.LRem(context.Background(), key, 1, bytes)
		}

		return nil
	})

	return err
}

// validate round id ++1
func (rd redisDB) IncrValidateRoundID() (int64, error) {
	key := fmt.Sprintf(redisKeyValidateRoundID, serverName)

	return rd.cli.IncrBy(context.Background(), key, 1).Result()
}

// verifying node list
func (rd redisDB) SetNodeToVerifyingList(deviceID string) error {
	key := fmt.Sprintf(redisKeyVerifyingList, serverName)

	_, err := rd.cli.SAdd(context.Background(), key, deviceID).Result()
	return err
}

func (rd redisDB) GetNodesWithVerifyingList() ([]string, error) {
	key := fmt.Sprintf(redisKeyVerifyingList, serverName)

	return rd.cli.SMembers(context.Background(), key).Result()
}

func (rd redisDB) CountVerifyingNode(ctx context.Context) (int64, error) {
	key := fmt.Sprintf(redisKeyVerifyingList, serverName)
	return rd.cli.SCard(ctx, key).Result()
}

func (rd redisDB) RemoveNodeWithVerifyingList(deviceID string) error {
	key := fmt.Sprintf(redisKeyVerifyingList, serverName)

	_, err := rd.cli.SRem(context.Background(), key, deviceID).Result()
	return err
}

func (rd redisDB) RemoveVerifyingList() error {
	key := fmt.Sprintf(redisKeyVerifyingList, serverName)

	_, err := rd.cli.Del(context.Background(), key).Result()
	return err
}

// validator list
func (rd redisDB) SetValidatorsToList(deviceIDs []string, expiration time.Duration) error {
	key := fmt.Sprintf(redisKeyValidatorList, serverName)

	// _, err := rd.cli.SAdd(context.Background(), key, deviceIDs).Result()

	ctx := context.Background()
	_, err := rd.cli.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		pipeliner.Del(context.Background(), key)

		if len(deviceIDs) > 0 {
			pipeliner.SAdd(context.Background(), key, deviceIDs)
			// Expire
			pipeliner.PExpire(context.Background(), key, expiration)
		}

		return nil
	})

	return err
}

func (rd redisDB) GetValidatorsWithList() ([]string, error) {
	key := fmt.Sprintf(redisKeyValidatorList, serverName)

	return rd.cli.SMembers(context.Background(), key).Result()
}

func (rd redisDB) GetValidatorsAndExpirationTime() ([]string, time.Duration, error) {
	key := fmt.Sprintf(redisKeyValidatorList, serverName)

	expiration, err := rd.cli.TTL(context.Background(), key).Result()
	if err != nil {
		return nil, expiration, err
	}

	list, err := rd.cli.SMembers(context.Background(), key).Result()
	if err != nil {
		return nil, expiration, err
	}

	return list, expiration, nil
}

// device info
func (rd redisDB) IncrNodeOnlineTime(deviceID string, onlineTime int64) (float64, error) {
	key := fmt.Sprintf(redisKeyNodeInfo, deviceID)

	mTime := float64(onlineTime) / 60 // second to minute

	return rd.cli.HIncrByFloat(context.Background(), key, "OnlineTime", mTime).Result()
}

func (rd redisDB) SetDeviceInfo(info *api.DevicesInfo) error {
	key := fmt.Sprintf(redisKeyNodeInfo, info.DeviceId)

	ctx := context.Background()
	exist, err := rd.cli.Exists(ctx, key).Result()
	if err != nil {
		return err
	}

	if exist == 1 {
		// update some value
		return rd.updateDeviceInfos(info)
	}

	// _, err = rd.cli.HMSet(ctx, key, structs.Map(info)).Result()
	// if err != nil {
	// 	return err
	// }

	_, err = rd.cli.Pipelined(ctx, func(pipeline redis.Pipeliner) error {
		for field, value := range toMap(info) {
			pipeline.HSet(ctx, key, field, value)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (rd redisDB) updateDeviceInfos(info *api.DevicesInfo) error {
	key := fmt.Sprintf(redisKeyNodeInfo, info.DeviceId)

	m := make(map[string]interface{})
	m["DiskSpace"] = info.DiskSpace
	m["DiskUsage"] = info.DiskUsage

	//TODO
	m["Longitude"] = info.Longitude
	m["Latitude"] = info.Latitude

	_, err := rd.cli.HMSet(context.Background(), key, m).Result()
	return err
}

func (rd redisDB) GetDeviceInfo(deviceID string) (*api.DevicesInfo, error) {
	key := fmt.Sprintf(redisKeyNodeInfo, deviceID)

	var info api.DevicesInfo
	err := rd.cli.HGetAll(context.Background(), key).Scan(&info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (rd redisDB) UpdateDeviceInfo(deviceID, field string, value interface{}) error {
	key := fmt.Sprintf(redisKeyNodeInfo, deviceID)

	_, err := rd.cli.HSet(context.Background(), key, field, value).Result()
	if err != nil {
		return err
	}

	return err
}

func (rd redisDB) IncrByDeviceInfo(deviceID, field string, value int64) error {
	key := fmt.Sprintf(redisKeyNodeInfo, deviceID)

	_, err := rd.cli.HIncrBy(context.Background(), key, field, value).Result()
	if err != nil {
		return err
	}

	return err
}

func (rd redisDB) CacheEndRecord(dataTask *DataTask, fromDeviceID string, blockSize int, blocks int, isSuccess bool) error {
	//CacheEndRecord update node cache block info and base info
	toKey := fmt.Sprintf(redisKeyNodeInfo, dataTask.DeviceID)

	ctx := context.Background()
	_, err := rd.cli.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {

		pipeliner.HIncrBy(context.Background(), toKey, BlockCountField, int64(blocks))
		pipeliner.HIncrBy(context.Background(), toKey, TotalDownloadField, int64(blockSize))

		if fromDeviceID != "" {
			fromKey := fmt.Sprintf(redisKeyNodeInfo, fromDeviceID)
			pipeliner.HIncrBy(context.Background(), fromKey, DownloadCountField, int64(blocks))
			pipeliner.HIncrBy(context.Background(), fromKey, TotalUploadField, int64(blockSize))
		}

		if isSuccess {
			baseInfoKey := fmt.Sprintf(redisKeyBaseInfo, serverName)
			pipeliner.HIncrBy(context.Background(), baseInfoKey, CarFileCountField, 1).Result()
		}

		return nil
	})

	return err
}

// download info
func (rd redisDB) SetDownloadBlockRecord(record *DownloadBlockRecord) error {
	ctx := context.Background()
	key := fmt.Sprintf(redisKeyBlockDownloadRecord, record.SN)
	_, err := rd.cli.HMSet(ctx, key, structs.Map(record)).Result()
	if err != nil {
		return err
	}
	_, err = rd.cli.Expire(ctx, key, time.Duration(record.Timeout)*time.Second).Result()
	if err != nil {
		return err
	}

	return nil
}

func (rd redisDB) GetDownloadBlockRecord(sn int64) (*DownloadBlockRecord, error) {
	key := fmt.Sprintf(redisKeyBlockDownloadRecord, sn)

	var record DownloadBlockRecord

	err := rd.cli.HGetAll(context.Background(), key).Scan(&record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (rd redisDB) RemoveDownloadBlockRecord(sn int64) error {
	key := fmt.Sprintf(redisKeyBlockDownloadRecord, sn)
	_, err := rd.cli.Del(context.Background(), key).Result()
	return err
}

func (rd redisDB) IncrBlockDownloadSN() (int64, error) {
	// node cache tag ++1
	n, err := rd.cli.IncrBy(context.Background(), redisKeyBlockDownloadSN, 1).Result()
	if n >= math.MaxInt64 {
		rd.cli.Set(context.Background(), redisKeyBlockDownloadSN, 0, 0).Result()
	}

	return n, err
}

// latest data of download
func (rd redisDB) AddLatestDownloadCarfile(carfileCID string, userIP string) error {
	maxCount := 5
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	key := fmt.Sprintf(redisKeyCarfileLatestDownload, userIP)

	err := rd.cli.ZAdd(ctx, key, &redis.Z{Score: float64(time.Now().Unix()), Member: carfileCID}).Err()
	if err != nil {
		return err
	}

	count, err := rd.cli.ZCard(ctx, key).Result()
	if err != nil {
		return err
	}

	if count > int64(maxCount) {
		err = rd.cli.ZRemRangeByRank(ctx, key, 0, count-int64(maxCount)-1).Err()
		if err != nil {
			return err
		}
	}

	return rd.cli.Expire(ctx, key, 24*time.Hour).Err()
}

func (rd redisDB) GetLatestDownloadCarfiles(userIP string) ([]string, error) {
	key := fmt.Sprintf(redisKeyCarfileLatestDownload, userIP)

	members, err := rd.cli.ZRevRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return []string{}, err
	}
	return members, nil
}

func (rd redisDB) NodeDownloadCount(deviceID string, blockDownnloadInfo *api.BlockDownloadInfo) error {
	nodeKey := fmt.Sprintf(redisKeyNodeInfo, deviceID)

	ctx := context.Background()
	_, err := rd.cli.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {

		pipeliner.HIncrBy(context.Background(), nodeKey, DownloadCountField, 1)
		pipeliner.HIncrBy(context.Background(), nodeKey, TotalUploadField, int64(blockDownnloadInfo.BlockSize))

		// count carfile download
		if blockDownnloadInfo.BlockCID == blockDownnloadInfo.CarfileCID {
			key := fmt.Sprintf(redisKeyBaseInfo, serverName)
			pipeliner.HIncrBy(context.Background(), key, DownloadCountField, 1)
		}

		return nil
	})

	return err
}

// IsNilErr Is NilErr
func (rd redisDB) IsNilErr(err error) bool {
	return errors.Is(err, redis.Nil)
}

func toMap(info *api.DevicesInfo) map[string]interface{} {
	out := make(map[string]interface{})
	t := reflect.TypeOf(*info)
	v := reflect.ValueOf(*info)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		redisTag := field.Tag.Get("redis")
		if redisTag == "" {
			continue
		}
		out[redisTag] = v.Field(i).Interface()
	}
	return out
}

// system base info
func (rd redisDB) UpdateBaseInfo(field string, value interface{}) error {
	key := fmt.Sprintf(redisKeyBaseInfo, serverName)

	_, err := rd.cli.HSet(context.Background(), key, field, value).Result()
	if err != nil {
		return err
	}

	return err
}

func (rd redisDB) IncrByBaseInfo(field string, value int64) error {
	key := fmt.Sprintf(redisKeyBaseInfo, serverName)

	_, err := rd.cli.HIncrBy(context.Background(), key, field, value).Result()
	if err != nil {
		return err
	}

	return err
}

func (rd redisDB) GetBaseInfo() (*api.BaseInfo, error) {
	key := fmt.Sprintf(redisKeyBaseInfo, serverName)

	var info api.BaseInfo
	err := rd.cli.HGetAll(context.Background(), key).Scan(&info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (rd redisDB) RemoveCarfileRecord(dataTasks []*DataTask, carfileCount int64, nodeBlockCounts map[string]int64) error {
	ctx := context.Background()
	_, err := rd.cli.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		baseInfoKey := fmt.Sprintf(redisKeyBaseInfo, serverName)
		pipeliner.HIncrBy(context.Background(), baseInfoKey, CarFileCountField, carfileCount).Result()

		for deviceID, blocks := range nodeBlockCounts {
			key := fmt.Sprintf(redisKeyNodeInfo, deviceID)
			pipeliner.HIncrBy(context.Background(), key, BlockCountField, blocks)
		}

		for _, dataTask := range dataTasks {
			nodeCacheList := fmt.Sprintf(redisKeyNodeCacheList, serverName, dataTask.DeviceID)

			pipeliner.SRem(context.Background(), nodeCacheList, dataTask.CarfileHash)
		}

		return nil
	})

	return err
}
