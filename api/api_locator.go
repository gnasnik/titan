package api

import (
	"context"

	"github.com/linguohua/titan/api/types"
)

type Locator interface {
	Common
	GetAccessPoints(ctx context.Context, nodeID string) ([]string, error)                                                  //perm:read
	AddAccessPoint(ctx context.Context, areaID string, schedulerURL string, weight int, schedulerAccessToken string) error //perm:admin
	RemoveAccessPoints(ctx context.Context, areaID string) error                                                           //perm:admin                                  //perm:admin
	ListAreaIDs(ctx context.Context) (areaIDs []string, err error)                                                         //perm:admin
	ShowAccessPoint(ctx context.Context, areaID string) (AccessPoint, error)                                               //perm:admin

	SetNodeOnlineStatus(ctx context.Context, nodeID string, isOnline bool) error //perm:write

	// user api
	EdgeDownloadInfos(ctx context.Context, cid string) ([]*types.DownloadInfo, error) //perm:read
	// user send result when user download block complete
	UserDownloadBlockResults(ctx context.Context, results []types.UserBlockDownloadResult) error //perm:read

	// api for web
	RegisterNode(ctx context.Context, schedulerURL, nodeID, publicKey string, nt types.NodeType) error // perm:admin
	LoadAccessPointsForWeb(ctx context.Context) ([]AccessPoint, error)                                 // perm:admin
	LoadUserAccessPoint(ctx context.Context, userIP string) (AccessPoint, error)                       // perm:admin
}

type SchedulerInfo struct {
	URL    string
	Weight int
	// Online      bool
	// AccessToken string
}
type AccessPoint struct {
	AreaID         string
	SchedulerInfos []SchedulerInfo
}
