package election

import (
	"context"
	"math"
	"sort"
	"sync"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"github.com/linguohua/titan/node/modules/dtypes"
	"github.com/linguohua/titan/node/scheduler/db"
	"github.com/linguohua/titan/node/scheduler/node"
)

var log = logging.Logger("election")

var (
	firstElectInterval = 5 * time.Minute    // Time of the first election
	electionCycle      = 5 * 24 * time.Hour // election cycle
)

// Election election
type Election struct {
	opts     EleOption
	updateCh chan struct{}
	nodeMgr  *node.Manager
	lock     sync.RWMutex
	config   dtypes.GetSchedulerConfigFunc
	*db.SQLDB
	dtypes.ServerID
}

// EleOption election option
type EleOption struct {
	ctx context.Context
}

func newEleOption(opts ...Option) EleOption {
	opt := EleOption{
		ctx: context.Background(),
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Option election option
type Option func(opt *EleOption)

// NewElection return new election instance
func NewElection(manager *node.Manager, configFunc dtypes.GetSchedulerConfigFunc, sdb *db.SQLDB, serverID dtypes.ServerID) *Election {
	ele := &Election{
		nodeMgr:  manager,
		opts:     newEleOption(),
		updateCh: make(chan struct{}, 1),
		config:   configFunc,
		SQLDB:    sdb,
		ServerID: serverID,
	}

	go ele.electTicker()

	return ele
}

func (v *Election) electTicker() {
	validators, err := v.LoadValidators(v.ServerID)
	if err != nil {
		log.Errorf("fetch current validators: %v", err)
		return
	}

	expiration := electionCycle
	if len(validators) <= 0 {
		expiration = firstElectInterval
	}

	electTicker := time.NewTicker(expiration)
	defer electTicker.Stop()

	doElect := func() {
		electTicker.Reset(electionCycle)
		err := v.elect()
		if err != nil {
			log.Errorf("elect err:%s", err.Error())
		}
	}

	for {
		select {
		case <-electTicker.C:
			doElect()
		case <-v.updateCh:
			doElect()
		case <-v.opts.ctx.Done():
			return
		}
	}
}

func (v *Election) elect() error {
	log.Debugln("start elect ")
	validators := v.electValidators()

	return v.UpdateValidators(validators, v.ServerID)
}

// StartElect elect
func (v *Election) StartElect() {
	v.updateCh <- struct{}{}
}

func (v *Election) getValidatorRatio() float64 {
	cfg, err := v.config()
	if err != nil {
		log.Errorf("schedulerConfig err:%s", err.Error())
		return 0
	}

	return cfg.ValidatorRatio
}

func (v *Election) electValidators() (out []string) {
	out = make([]string, 0)

	// TODO problematic
	candidates := v.getAllCandidates()
	candidateCount := len(candidates)

	needValidatorCount := int(math.Ceil(float64(candidateCount) * v.getValidatorRatio()))
	if needValidatorCount <= 0 {
		return
	}

	if needValidatorCount >= candidateCount {
		for _, candidate := range candidates {
			out = append(out, candidate.NodeID)
		}
		return
	}

	sort.Slice(candidates, func(i, j int) bool {
		// TODO Consider node reliability
		return candidates[i].BandwidthDown > candidates[j].BandwidthDown
	})

	for i := 0; i < needValidatorCount; i++ {
		candidate := candidates[i]
		out = append(out, candidate.NodeID)
	}

	return out
}

func (v *Election) getAllCandidates() []*node.Candidate {
	candidates := make([]*node.Candidate, 0)

	v.nodeMgr.CandidateNodes.Range(func(key, value interface{}) bool {
		node := value.(*node.Candidate)
		candidates = append(candidates, node)

		return true
	})

	return candidates
}
