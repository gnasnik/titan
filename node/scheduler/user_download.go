package scheduler

import (
	"context"

	"github.com/linguohua/titan/api/types"
	"github.com/linguohua/titan/node/cidutil"
	"github.com/linguohua/titan/node/handler"
	"golang.org/x/xerrors"
)

// UserDownloadResult result for user download
func (s *Scheduler) UserDownloadResult(ctx context.Context, result types.UserDownloadResult) error {
	nodeID := handler.GetNodeID(ctx)
	if result.Succeed {
		blockHash, err := cidutil.CIDString2HashString(result.BlockCID)
		if err != nil {
			return err
		}

		blockDownloadInfo := &types.DownloadRecordInfo{NodeID: nodeID, BlockCID: result.BlockCID, BlockSize: result.BlockSize}

		carfileInfo, _ := s.NodeManager.NodeMgrDB.LoadCarfileRecordInfo(blockHash)
		if carfileInfo != nil && carfileInfo.CarfileCID != "" {
			blockDownloadInfo.CarfileCID = result.BlockCID
		}
	}

	return nil
}

func (s *Scheduler) handleUserDownloadBlockResult(ctx context.Context, result types.UserBlockDownloadResult) error {
	// TODO: implement user download count
	return nil
}

// UserDownloadBlockResults node result for user download block
func (s *Scheduler) UserDownloadBlockResults(ctx context.Context, results []types.UserBlockDownloadResult) error {
	for _, result := range results {
		s.handleUserDownloadBlockResult(ctx, result)
	}
	return nil
}

// EdgeDownloadInfos find node
func (s *Scheduler) EdgeDownloadInfos(ctx context.Context, cid string) ([]*types.DownloadInfo, error) {
	if cid == "" {
		return nil, xerrors.New("cids is nil")
	}

	userURL := handler.GetRemoteAddr(ctx)

	log.Infof("EdgeDownloadInfos url:%s", userURL)

	infos, err := s.NodeManager.FindNodeDownloadInfos(cid, userURL)
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func (s *Scheduler) getNodesUnValidate(minute int) ([]string, error) {
	return s.NodeManager.NodeMgrDB.GetNodesByUserDownloadBlockIn(minute)
}

// func (s *Scheduler) getNodePublicKey(nodeID string) (*rsa.PublicKey, error) {
// 	pem, err := s.NodeManager.NodeMgrDB.NodePublicKey(nodeID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return titanrsa.Pem2PublicKey([]byte(pem))
// }
