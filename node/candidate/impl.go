package candidate

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/linguohua/titan/node/config"
	"go.uber.org/fx"

	"github.com/linguohua/titan/api"
	"github.com/linguohua/titan/node/carfile"
	"github.com/linguohua/titan/node/common"
	"github.com/linguohua/titan/node/device"
	datasync "github.com/linguohua/titan/node/sync"

	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	vd "github.com/linguohua/titan/node/validate"
	mh "github.com/multiformats/go-multihash"
)

var log = logging.Logger("candidate")

const (
	schedulerAPITimeout = 3
	validateTimeout     = 5
	tcpPackMaxLength    = 52428800
)

func cidFromData(data []byte) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("len(data) == 0")
	}

	pref := cid.Prefix{
		Version:  1,
		Codec:    uint64(cid.Raw),
		MhType:   mh.SHA2_256,
		MhLength: -1, // default length
	}

	c, err := pref.Sum(data)
	if err != nil {
		return "", err
	}

	return c.String(), nil
}

type blockWaiter struct {
	conn *net.TCPConn
	ch   chan tcpMsg
}

type Candidate struct {
	fx.In

	*common.CommonAPI
	*carfile.CarfileImpl
	*device.Device
	*vd.Validate
	*datasync.DataSync

	Scheduler      api.Scheduler
	Config         *config.CandidateCfg
	BlockWaiterMap *BlockWaiter
	TCPSrv         *TCPServer
}

type BlockWaiter struct {
	sync.Map
}

func NewBlockWaiter() *BlockWaiter {
	return &BlockWaiter{}
}

func (candidate *Candidate) WaitQuiet(ctx context.Context) error {
	log.Debug("WaitQuiet")
	return nil
}

func (candidate *Candidate) GetBlocksOfCarfile(ctx context.Context, carfileCID string, randomSeed int64, randomCount int) (map[int]string, error) {
	blockCount, err := candidate.CarfileImpl.BlockCountOfCarfile(carfileCID)
	if err != nil {
		log.Errorf("GetBlocksOfCarfile, BlockCountOfCarfile error:%s, carfileCID:%s", err.Error(), carfileCID)
		return nil, err
	}

	indices := make([]int, 0)
	indexMap := make(map[int]struct{})
	r := rand.New(rand.NewSource(randomSeed))

	for i := 0; i < randomCount; i++ {
		index := r.Intn(blockCount)

		if _, ok := indexMap[index]; !ok {
			indices = append(indices, index)
			indexMap[index] = struct{}{}
		}
	}

	return candidate.CarfileImpl.GetBlocksOfCarfile(carfileCID, indices)
}

func (candidate *Candidate) ValidateNodes(ctx context.Context, req []api.ValidateReq) (string, error) {
	for _, reqValidate := range req {
		prepareValidate(&reqValidate, candidate)
	}

	address, err := candidate.Scheduler.NodeExternalServiceAddress(context.Background())
	if err != nil {
		log.Errorf("can not get external service address: %s", err.Error())
		return "", err
	}

	host, _, err := net.SplitHostPort(address)
	if err != nil {
		log.Errorf("can not get external service address: %s", err.Error())
		return "", err
	}

	_, port, err := net.SplitHostPort(candidate.Config.TCPSrvAddr)
	if err != nil {
		log.Errorf("can not get external service address: %s", err.Error())
		return "", err
	}

	candidateTCPSrvAddr := fmt.Sprintf("%s:%s", host, port)
	return candidateTCPSrvAddr, nil
}

func (candidate *Candidate) loadBlockWaiterFromMap(key string) (*blockWaiter, bool) {
	vb, exist := candidate.BlockWaiterMap.Load(key)
	if exist {
		return vb.(*blockWaiter), exist
	}
	return nil, exist
}

func sendValidateResult(candidate *Candidate, result *api.ValidateResult) error {
	ctx, cancel := context.WithTimeout(context.Background(), schedulerAPITimeout*time.Second)
	defer cancel()

	return candidate.Scheduler.NodeValidatedResult(ctx, *result)
}

func waitBlock(vb *blockWaiter, req *api.ValidateReq, candidate *Candidate, result *api.ValidateResult) {
	defer func() {
		candidate.BlockWaiterMap.Delete(result.NodeID)
	}()

	size := int64(0)
	now := time.Now()
	isBreak := false
	t := time.NewTimer(time.Duration(req.Duration+validateTimeout) * time.Second)
	for {
		select {
		case tcpMsg, ok := <-vb.ch:
			if !ok {
				// log.Infof("waitblock close channel %s", result.NodeID)
				isBreak = true
				vb.ch = nil
				break
			}

			if tcpMsg.msgType == api.TCPMsgTypeCancel {
				result.IsCancel = true
				if err := sendValidateResult(candidate, result); err != nil {
					log.Errorf("node %s cancel validator, send validate result error: %s", result.NodeID, err.Error())
				}

				log.Infof("node %s cancel validator", result.NodeID)
				return
			}

			if tcpMsg.msgType == api.TCPMsgTypeBlock && len(tcpMsg.msg) > 0 {
				cid, err := cidFromData(tcpMsg.msg)
				if err != nil {
					log.Errorf("waitBlock, cidFromData error:%v", err)
				}
				result.Cids = append(result.Cids, cid)
			}
			size += int64(tcpMsg.length)
			result.RandomCount++
		case <-t.C:
			if vb.conn != nil {
				if err := vb.conn.Close(); err != nil {
					log.Errorf("close tcp error: %s", err.Error())
				}
			}
			isBreak = true
			log.Errorf("wait node %s timeout %ds, exit wait block", result.NodeID, req.Duration+validateTimeout)
		}

		if isBreak {
			break
		}

	}

	duration := time.Since(now)
	result.CostTime = int64(duration / time.Millisecond)

	if duration < time.Duration(req.Duration)*time.Second {
		duration = time.Duration(req.Duration) * time.Second
	}
	result.Bandwidth = float64(size) / float64(duration) * float64(time.Second)

	if err := sendValidateResult(candidate, result); err != nil {
		log.Errorf("send validate result error: %s", err.Error())
	}

	log.Infof("validator %s %d block, bandwidth:%f, cost time:%d, IsTimeout:%v, duration:%d, size:%d, randCount:%d",
		result.NodeID, len(result.Cids), result.Bandwidth, result.CostTime, result.IsTimeout, req.Duration, size, result.RandomCount)
}

func prepareValidate(req *api.ValidateReq, candidate *Candidate) {
	result := &api.ValidateResult{NodeID: req.NodeID, RoundID: req.RoundID, Cids: make([]string, 0)}

	if _, exist := candidate.loadBlockWaiterFromMap(req.NodeID); exist {
		log.Warnf("Already validating nodeID %s, not need to repeat do", req.NodeID)
		return
	}

	bw := &blockWaiter{conn: nil, ch: make(chan tcpMsg, 1)}
	candidate.BlockWaiterMap.Store(req.NodeID, bw)

	go waitBlock(bw, req, candidate, result)

	return
}
