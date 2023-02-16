// Code generated by titan/gen/api. DO NOT EDIT.

package api

import (
	"context"
	"time"

	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/google/uuid"
	"github.com/linguohua/titan/journal/alerting"
	xerrors "golang.org/x/xerrors"
)

var ErrNotSupported = xerrors.New("method not supported")

type CandidateStruct struct {
	CommonStruct

	DeviceStruct

	DownloadStruct

	ValidateStruct

	DataSyncStruct

	CarfileOperationStruct

	Internal struct {
		GetBlocksOfCarfile func(p0 context.Context, p1 string, p2 int64, p3 int) (map[int]string, error) `perm:"read"`

		LoadBlock func(p0 context.Context, p1 string) ([]byte, error) `perm:"read"`

		ValidateNodes func(p0 context.Context, p1 []ReqValidate) error `perm:"read"`

		WaitQuiet func(p0 context.Context) error `perm:"read"`
	}
}

type CandidateStub struct {
	CommonStub

	DeviceStub

	DownloadStub

	ValidateStub

	DataSyncStub

	CarfileOperationStub
}

type CarfileOperationStruct struct {
	Internal struct {
		CacheCarfile func(p0 context.Context, p1 string, p2 []*DowloadSource) (*CacheCarfileResult, error) `perm:"write"`

		DeleteAllCarfiles func(p0 context.Context) error `perm:"admin"`

		DeleteCarfile func(p0 context.Context, p1 string) error `perm:"write"`

		QueryCacheStat func(p0 context.Context) (*CacheStat, error) `perm:"write"`

		QueryCachingCarfile func(p0 context.Context) (*CachingCarfile, error) `perm:"write"`
	}
}

type CarfileOperationStub struct {
}

type CommonStruct struct {
	Internal struct {
		AuthNew func(p0 context.Context, p1 []auth.Permission) ([]byte, error) `perm:"admin"`

		AuthVerify func(p0 context.Context, p1 string) ([]auth.Permission, error) `perm:"read"`

		Closing func(p0 context.Context) (<-chan struct{}, error) `perm:"read"`

		Discover func(p0 context.Context) (OpenRPCDocument, error) `perm:"read"`

		DownloadLogFile func(p0 context.Context) ([]byte, error) `perm:"write"`

		LogAlerts func(p0 context.Context) ([]alerting.Alert, error) `perm:"admin"`

		LogList func(p0 context.Context) ([]string, error) `perm:"write"`

		LogSetLevel func(p0 context.Context, p1 string, p2 string) error `perm:"write"`

		Session func(p0 context.Context) (uuid.UUID, error) `perm:"read"`

		ShowLogFile func(p0 context.Context) (*LogFile, error) `perm:"write"`

		Shutdown func(p0 context.Context) error `perm:"admin"`

		Version func(p0 context.Context) (APIVersion, error) `perm:"read"`
	}
}

type CommonStub struct {
}

type DataSyncStruct struct {
	Internal struct {
	}
}

type DataSyncStub struct {
}

type DeviceStruct struct {
	Internal struct {
		DeviceID func(p0 context.Context) (string, error) `perm:"read"`

		DeviceInfo func(p0 context.Context) (DevicesInfo, error) `perm:"read"`
	}
}

type DeviceStub struct {
}

type DownloadStruct struct {
	Internal struct {
		SetDownloadSpeed func(p0 context.Context, p1 int64) error `perm:"write"`
	}
}

type DownloadStub struct {
}

type EdgeStruct struct {
	CommonStruct

	DeviceStruct

	DownloadStruct

	ValidateStruct

	DataSyncStruct

	CarfileOperationStruct

	Internal struct {
		PingUser func(p0 context.Context, p1 string) error `perm:"write"`

		WaitQuiet func(p0 context.Context) error `perm:"read"`
	}
}

type EdgeStub struct {
	CommonStub

	DeviceStub

	DownloadStub

	ValidateStub

	DataSyncStub

	CarfileOperationStub
}

type LocatorStruct struct {
	CommonStruct

	Internal struct {
		AddAccessPoint func(p0 context.Context, p1 string, p2 string, p3 int, p4 string) error `perm:"admin"`

		GetAccessPoints func(p0 context.Context, p1 string) ([]string, error) `perm:"read"`

		GetDownloadInfosWithCarfile func(p0 context.Context, p1 string, p2 string) ([]*DownloadInfoResult, error) `perm:"read"`

		ListAreaIDs func(p0 context.Context) ([]string, error) `perm:"admin"`

		LoadAccessPointsForWeb func(p0 context.Context) ([]AccessPoint, error) `perm:"admin"`

		LoadUserAccessPoint func(p0 context.Context, p1 string) (AccessPoint, error) `perm:"admin"`

		RegisterNode func(p0 context.Context, p1 string, p2 string, p3 NodeType, p4 int) ([]NodeRegisterInfo, error) `perm:"admin"`

		RemoveAccessPoints func(p0 context.Context, p1 string) error `perm:"admin"`

		SetDeviceOnlineStatus func(p0 context.Context, p1 string, p2 bool) error `perm:"write"`

		ShowAccessPoint func(p0 context.Context, p1 string) (AccessPoint, error) `perm:"admin"`

		UserDownloadBlockResults func(p0 context.Context, p1 []UserBlockDownloadResult) error `perm:"read"`
	}
}

type LocatorStub struct {
	CommonStub
}

type SchedulerStruct struct {
	CommonStruct

	WebStruct

	Internal struct {
		AuthNodeNew func(p0 context.Context, p1 []auth.Permission, p2 string, p3 string) ([]byte, error) `perm:"read"`

		AuthNodeVerify func(p0 context.Context, p1 string) ([]auth.Permission, error) `perm:"read"`

		CacheCarfile func(p0 context.Context, p1 *CacheCarfileInfo) error `perm:"admin"`

		CacheResult func(p0 context.Context, p1 CacheResultInfo) error `perm:"write"`

		CandidateNodeConnect func(p0 context.Context) error `perm:"write"`

		EdgeNodeConnect func(p0 context.Context) error `perm:"write"`

		ElectionValidators func(p0 context.Context) error `perm:"admin"`

		ExecuteUndoneCarfilesTask func(p0 context.Context) error `perm:"admin"`

		GetCarfileRecordInfo func(p0 context.Context, p1 string) (CarfileRecordInfo, error) `perm:"read"`

		GetDevicesInfo func(p0 context.Context, p1 string) (DevicesInfo, error) `perm:"read"`

		GetDownloadInfo func(p0 context.Context, p1 string) ([]*BlockDownloadInfo, error) `perm:"read"`

		GetDownloadInfosWithCarfile func(p0 context.Context, p1 string, p2 string) ([]*DownloadInfoResult, error) `perm:"read"`

		GetExternalIP func(p0 context.Context) (string, error) `perm:"write"`

		GetNodeAppUpdateInfos func(p0 context.Context) (map[int]*NodeAppUpdateInfo, error) `perm:"read"`

		GetOnlineDeviceIDs func(p0 context.Context, p1 NodeTypeName) ([]string, error) `perm:"read"`

		GetPublicKey func(p0 context.Context) (string, error) `perm:"write"`

		GetRunningCarfileRecords func(p0 context.Context) ([]*CarfileRecordInfo, error) `perm:"read"`

		GetUndoneCarfileRecords func(p0 context.Context, p1 int) (*DataListInfo, error) `perm:"read"`

		ListCarfileRecords func(p0 context.Context, p1 int) (*DataListInfo, error) `perm:"read"`

		LocatorConnect func(p0 context.Context, p1 string, p2 string) error `perm:"write"`

		NodeQuit func(p0 context.Context, p1 string) error `perm:"admin"`

		NodeResultForUserDownloadBlock func(p0 context.Context, p1 NodeBlockDownloadResult) error `perm:"write"`

		RegisterNode func(p0 context.Context, p1 NodeType, p2 int) ([]NodeRegisterInfo, error) `perm:"admin"`

		RemoveCache func(p0 context.Context, p1 string, p2 string) error `perm:"admin"`

		RemoveCarfile func(p0 context.Context, p1 string) error `perm:"admin"`

		RemoveCarfileResult func(p0 context.Context, p1 RemoveCarfileResultInfo) error `perm:"write"`

		ResetBackupCacheCount func(p0 context.Context, p1 int) error `perm:"admin"`

		ResetCacheExpiredTime func(p0 context.Context, p1 string, p2 time.Time) error `perm:"admin"`

		SetNodeAppUpdateInfo func(p0 context.Context, p1 *NodeAppUpdateInfo) error `perm:"admin"`

		StopCacheTask func(p0 context.Context, p1 string) error `perm:"admin"`

		UserDownloadBlockResults func(p0 context.Context, p1 []UserBlockDownloadResult) error `perm:"read"`

		ValidateBlockResult func(p0 context.Context, p1 ValidateResults) error `perm:"write"`

		ValidateRunningState func(p0 context.Context) (bool, error) `perm:"admin"`

		ValidateStart func(p0 context.Context) error `perm:"admin"`

		ValidateSwitch func(p0 context.Context, p1 bool) error `perm:"admin"`
	}
}

type SchedulerStub struct {
	CommonStub

	WebStub
}

type ValidateStruct struct {
	Internal struct {
		BeValidate func(p0 context.Context, p1 ReqValidate, p2 string) error `perm:"read"`
	}
}

type ValidateStub struct {
}

type WebStruct struct {
	Internal struct {
		AddCacheTask func(p0 context.Context, p1 string, p2 int, p3 time.Time) error `perm:"read"`

		CancelCacheTask func(p0 context.Context, p1 string) error `perm:"read"`

		GetCacheTaskInfo func(p0 context.Context, p1 string) (CarfileRecordInfo, error) `perm:"read"`

		GetCacheTaskInfos func(p0 context.Context, p1 ListCacheInfosReq) (ListCacheInfosRsp, error) `perm:"read"`

		GetCarfileByCID func(p0 context.Context, p1 string) (WebCarfile, error) `perm:"read"`

		GetNodeInfoByID func(p0 context.Context, p1 string) (DevicesInfo, error) `perm:"read"`

		GetSummaryValidateMessage func(p0 context.Context, p1 time.Time, p2 time.Time, p3 int, p4 int) (*SummeryValidateResult, error) `perm:"read"`

		GetSystemInfo func(p0 context.Context) (SystemBaseInfo, error) `perm:"read"`

		ListBlockDownloadInfo func(p0 context.Context, p1 ListBlockDownloadInfoReq) (ListBlockDownloadInfoRsp, error) `perm:"read"`

		ListCacheTasks func(p0 context.Context, p1 int, p2 int) (ListCacheTasksRsp, error) `perm:"read"`

		ListNodeConnectionLog func(p0 context.Context, p1 ListNodeConnectionLogReq) (ListNodeConnectionLogRsp, error) `perm:"read"`

		ListNodes func(p0 context.Context, p1 int, p2 int) (ListNodesRsp, error) `perm:"read"`

		ListValidateResult func(p0 context.Context, p1 int, p2 int) (ListValidateResultRsp, error) `perm:"read"`

		RemoveCarfile func(p0 context.Context, p1 string) error `perm:"read"`

		SetupValidation func(p0 context.Context, p1 bool) error `perm:"read"`
	}
}

type WebStub struct {
}

func (s *CandidateStruct) GetBlocksOfCarfile(p0 context.Context, p1 string, p2 int64, p3 int) (map[int]string, error) {
	if s.Internal.GetBlocksOfCarfile == nil {
		return *new(map[int]string), ErrNotSupported
	}
	return s.Internal.GetBlocksOfCarfile(p0, p1, p2, p3)
}

func (s *CandidateStub) GetBlocksOfCarfile(p0 context.Context, p1 string, p2 int64, p3 int) (map[int]string, error) {
	return *new(map[int]string), ErrNotSupported
}

func (s *CandidateStruct) LoadBlock(p0 context.Context, p1 string) ([]byte, error) {
	if s.Internal.LoadBlock == nil {
		return *new([]byte), ErrNotSupported
	}
	return s.Internal.LoadBlock(p0, p1)
}

func (s *CandidateStub) LoadBlock(p0 context.Context, p1 string) ([]byte, error) {
	return *new([]byte), ErrNotSupported
}

func (s *CandidateStruct) ValidateNodes(p0 context.Context, p1 []ReqValidate) error {
	if s.Internal.ValidateNodes == nil {
		return ErrNotSupported
	}
	return s.Internal.ValidateNodes(p0, p1)
}

func (s *CandidateStub) ValidateNodes(p0 context.Context, p1 []ReqValidate) error {
	return ErrNotSupported
}

func (s *CandidateStruct) WaitQuiet(p0 context.Context) error {
	if s.Internal.WaitQuiet == nil {
		return ErrNotSupported
	}
	return s.Internal.WaitQuiet(p0)
}

func (s *CandidateStub) WaitQuiet(p0 context.Context) error {
	return ErrNotSupported
}

func (s *CarfileOperationStruct) CacheCarfile(p0 context.Context, p1 string, p2 []*DowloadSource) (*CacheCarfileResult, error) {
	if s.Internal.CacheCarfile == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.CacheCarfile(p0, p1, p2)
}

func (s *CarfileOperationStub) CacheCarfile(p0 context.Context, p1 string, p2 []*DowloadSource) (*CacheCarfileResult, error) {
	return nil, ErrNotSupported
}

func (s *CarfileOperationStruct) DeleteAllCarfiles(p0 context.Context) error {
	if s.Internal.DeleteAllCarfiles == nil {
		return ErrNotSupported
	}
	return s.Internal.DeleteAllCarfiles(p0)
}

func (s *CarfileOperationStub) DeleteAllCarfiles(p0 context.Context) error {
	return ErrNotSupported
}

func (s *CarfileOperationStruct) DeleteCarfile(p0 context.Context, p1 string) error {
	if s.Internal.DeleteCarfile == nil {
		return ErrNotSupported
	}
	return s.Internal.DeleteCarfile(p0, p1)
}

func (s *CarfileOperationStub) DeleteCarfile(p0 context.Context, p1 string) error {
	return ErrNotSupported
}

func (s *CarfileOperationStruct) QueryCacheStat(p0 context.Context) (*CacheStat, error) {
	if s.Internal.QueryCacheStat == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.QueryCacheStat(p0)
}

func (s *CarfileOperationStub) QueryCacheStat(p0 context.Context) (*CacheStat, error) {
	return nil, ErrNotSupported
}

func (s *CarfileOperationStruct) QueryCachingCarfile(p0 context.Context) (*CachingCarfile, error) {
	if s.Internal.QueryCachingCarfile == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.QueryCachingCarfile(p0)
}

func (s *CarfileOperationStub) QueryCachingCarfile(p0 context.Context) (*CachingCarfile, error) {
	return nil, ErrNotSupported
}

func (s *CommonStruct) AuthNew(p0 context.Context, p1 []auth.Permission) ([]byte, error) {
	if s.Internal.AuthNew == nil {
		return *new([]byte), ErrNotSupported
	}
	return s.Internal.AuthNew(p0, p1)
}

func (s *CommonStub) AuthNew(p0 context.Context, p1 []auth.Permission) ([]byte, error) {
	return *new([]byte), ErrNotSupported
}

func (s *CommonStruct) AuthVerify(p0 context.Context, p1 string) ([]auth.Permission, error) {
	if s.Internal.AuthVerify == nil {
		return *new([]auth.Permission), ErrNotSupported
	}
	return s.Internal.AuthVerify(p0, p1)
}

func (s *CommonStub) AuthVerify(p0 context.Context, p1 string) ([]auth.Permission, error) {
	return *new([]auth.Permission), ErrNotSupported
}

func (s *CommonStruct) Closing(p0 context.Context) (<-chan struct{}, error) {
	if s.Internal.Closing == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.Closing(p0)
}

func (s *CommonStub) Closing(p0 context.Context) (<-chan struct{}, error) {
	return nil, ErrNotSupported
}

func (s *CommonStruct) Discover(p0 context.Context) (OpenRPCDocument, error) {
	if s.Internal.Discover == nil {
		return *new(OpenRPCDocument), ErrNotSupported
	}
	return s.Internal.Discover(p0)
}

func (s *CommonStub) Discover(p0 context.Context) (OpenRPCDocument, error) {
	return *new(OpenRPCDocument), ErrNotSupported
}

func (s *CommonStruct) DownloadLogFile(p0 context.Context) ([]byte, error) {
	if s.Internal.DownloadLogFile == nil {
		return *new([]byte), ErrNotSupported
	}
	return s.Internal.DownloadLogFile(p0)
}

func (s *CommonStub) DownloadLogFile(p0 context.Context) ([]byte, error) {
	return *new([]byte), ErrNotSupported
}

func (s *CommonStruct) LogAlerts(p0 context.Context) ([]alerting.Alert, error) {
	if s.Internal.LogAlerts == nil {
		return *new([]alerting.Alert), ErrNotSupported
	}
	return s.Internal.LogAlerts(p0)
}

func (s *CommonStub) LogAlerts(p0 context.Context) ([]alerting.Alert, error) {
	return *new([]alerting.Alert), ErrNotSupported
}

func (s *CommonStruct) LogList(p0 context.Context) ([]string, error) {
	if s.Internal.LogList == nil {
		return *new([]string), ErrNotSupported
	}
	return s.Internal.LogList(p0)
}

func (s *CommonStub) LogList(p0 context.Context) ([]string, error) {
	return *new([]string), ErrNotSupported
}

func (s *CommonStruct) LogSetLevel(p0 context.Context, p1 string, p2 string) error {
	if s.Internal.LogSetLevel == nil {
		return ErrNotSupported
	}
	return s.Internal.LogSetLevel(p0, p1, p2)
}

func (s *CommonStub) LogSetLevel(p0 context.Context, p1 string, p2 string) error {
	return ErrNotSupported
}

func (s *CommonStruct) Session(p0 context.Context) (uuid.UUID, error) {
	if s.Internal.Session == nil {
		return *new(uuid.UUID), ErrNotSupported
	}
	return s.Internal.Session(p0)
}

func (s *CommonStub) Session(p0 context.Context) (uuid.UUID, error) {
	return *new(uuid.UUID), ErrNotSupported
}

func (s *CommonStruct) ShowLogFile(p0 context.Context) (*LogFile, error) {
	if s.Internal.ShowLogFile == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ShowLogFile(p0)
}

func (s *CommonStub) ShowLogFile(p0 context.Context) (*LogFile, error) {
	return nil, ErrNotSupported
}

func (s *CommonStruct) Shutdown(p0 context.Context) error {
	if s.Internal.Shutdown == nil {
		return ErrNotSupported
	}
	return s.Internal.Shutdown(p0)
}

func (s *CommonStub) Shutdown(p0 context.Context) error {
	return ErrNotSupported
}

func (s *CommonStruct) Version(p0 context.Context) (APIVersion, error) {
	if s.Internal.Version == nil {
		return *new(APIVersion), ErrNotSupported
	}
	return s.Internal.Version(p0)
}

func (s *CommonStub) Version(p0 context.Context) (APIVersion, error) {
	return *new(APIVersion), ErrNotSupported
}

func (s *DeviceStruct) DeviceID(p0 context.Context) (string, error) {
	if s.Internal.DeviceID == nil {
		return "", ErrNotSupported
	}
	return s.Internal.DeviceID(p0)
}

func (s *DeviceStub) DeviceID(p0 context.Context) (string, error) {
	return "", ErrNotSupported
}

func (s *DeviceStruct) DeviceInfo(p0 context.Context) (DevicesInfo, error) {
	if s.Internal.DeviceInfo == nil {
		return *new(DevicesInfo), ErrNotSupported
	}
	return s.Internal.DeviceInfo(p0)
}

func (s *DeviceStub) DeviceInfo(p0 context.Context) (DevicesInfo, error) {
	return *new(DevicesInfo), ErrNotSupported
}

func (s *DownloadStruct) SetDownloadSpeed(p0 context.Context, p1 int64) error {
	if s.Internal.SetDownloadSpeed == nil {
		return ErrNotSupported
	}
	return s.Internal.SetDownloadSpeed(p0, p1)
}

func (s *DownloadStub) SetDownloadSpeed(p0 context.Context, p1 int64) error {
	return ErrNotSupported
}

func (s *EdgeStruct) PingUser(p0 context.Context, p1 string) error {
	if s.Internal.PingUser == nil {
		return ErrNotSupported
	}
	return s.Internal.PingUser(p0, p1)
}

func (s *EdgeStub) PingUser(p0 context.Context, p1 string) error {
	return ErrNotSupported
}

func (s *EdgeStruct) WaitQuiet(p0 context.Context) error {
	if s.Internal.WaitQuiet == nil {
		return ErrNotSupported
	}
	return s.Internal.WaitQuiet(p0)
}

func (s *EdgeStub) WaitQuiet(p0 context.Context) error {
	return ErrNotSupported
}

func (s *LocatorStruct) AddAccessPoint(p0 context.Context, p1 string, p2 string, p3 int, p4 string) error {
	if s.Internal.AddAccessPoint == nil {
		return ErrNotSupported
	}
	return s.Internal.AddAccessPoint(p0, p1, p2, p3, p4)
}

func (s *LocatorStub) AddAccessPoint(p0 context.Context, p1 string, p2 string, p3 int, p4 string) error {
	return ErrNotSupported
}

func (s *LocatorStruct) GetAccessPoints(p0 context.Context, p1 string) ([]string, error) {
	if s.Internal.GetAccessPoints == nil {
		return *new([]string), ErrNotSupported
	}
	return s.Internal.GetAccessPoints(p0, p1)
}

func (s *LocatorStub) GetAccessPoints(p0 context.Context, p1 string) ([]string, error) {
	return *new([]string), ErrNotSupported
}

func (s *LocatorStruct) GetDownloadInfosWithCarfile(p0 context.Context, p1 string, p2 string) ([]*DownloadInfoResult, error) {
	if s.Internal.GetDownloadInfosWithCarfile == nil {
		return *new([]*DownloadInfoResult), ErrNotSupported
	}
	return s.Internal.GetDownloadInfosWithCarfile(p0, p1, p2)
}

func (s *LocatorStub) GetDownloadInfosWithCarfile(p0 context.Context, p1 string, p2 string) ([]*DownloadInfoResult, error) {
	return *new([]*DownloadInfoResult), ErrNotSupported
}

func (s *LocatorStruct) ListAreaIDs(p0 context.Context) ([]string, error) {
	if s.Internal.ListAreaIDs == nil {
		return *new([]string), ErrNotSupported
	}
	return s.Internal.ListAreaIDs(p0)
}

func (s *LocatorStub) ListAreaIDs(p0 context.Context) ([]string, error) {
	return *new([]string), ErrNotSupported
}

func (s *LocatorStruct) LoadAccessPointsForWeb(p0 context.Context) ([]AccessPoint, error) {
	if s.Internal.LoadAccessPointsForWeb == nil {
		return *new([]AccessPoint), ErrNotSupported
	}
	return s.Internal.LoadAccessPointsForWeb(p0)
}

func (s *LocatorStub) LoadAccessPointsForWeb(p0 context.Context) ([]AccessPoint, error) {
	return *new([]AccessPoint), ErrNotSupported
}

func (s *LocatorStruct) LoadUserAccessPoint(p0 context.Context, p1 string) (AccessPoint, error) {
	if s.Internal.LoadUserAccessPoint == nil {
		return *new(AccessPoint), ErrNotSupported
	}
	return s.Internal.LoadUserAccessPoint(p0, p1)
}

func (s *LocatorStub) LoadUserAccessPoint(p0 context.Context, p1 string) (AccessPoint, error) {
	return *new(AccessPoint), ErrNotSupported
}

func (s *LocatorStruct) RegisterNode(p0 context.Context, p1 string, p2 string, p3 NodeType, p4 int) ([]NodeRegisterInfo, error) {
	if s.Internal.RegisterNode == nil {
		return *new([]NodeRegisterInfo), ErrNotSupported
	}
	return s.Internal.RegisterNode(p0, p1, p2, p3, p4)
}

func (s *LocatorStub) RegisterNode(p0 context.Context, p1 string, p2 string, p3 NodeType, p4 int) ([]NodeRegisterInfo, error) {
	return *new([]NodeRegisterInfo), ErrNotSupported
}

func (s *LocatorStruct) RemoveAccessPoints(p0 context.Context, p1 string) error {
	if s.Internal.RemoveAccessPoints == nil {
		return ErrNotSupported
	}
	return s.Internal.RemoveAccessPoints(p0, p1)
}

func (s *LocatorStub) RemoveAccessPoints(p0 context.Context, p1 string) error {
	return ErrNotSupported
}

func (s *LocatorStruct) SetDeviceOnlineStatus(p0 context.Context, p1 string, p2 bool) error {
	if s.Internal.SetDeviceOnlineStatus == nil {
		return ErrNotSupported
	}
	return s.Internal.SetDeviceOnlineStatus(p0, p1, p2)
}

func (s *LocatorStub) SetDeviceOnlineStatus(p0 context.Context, p1 string, p2 bool) error {
	return ErrNotSupported
}

func (s *LocatorStruct) ShowAccessPoint(p0 context.Context, p1 string) (AccessPoint, error) {
	if s.Internal.ShowAccessPoint == nil {
		return *new(AccessPoint), ErrNotSupported
	}
	return s.Internal.ShowAccessPoint(p0, p1)
}

func (s *LocatorStub) ShowAccessPoint(p0 context.Context, p1 string) (AccessPoint, error) {
	return *new(AccessPoint), ErrNotSupported
}

func (s *LocatorStruct) UserDownloadBlockResults(p0 context.Context, p1 []UserBlockDownloadResult) error {
	if s.Internal.UserDownloadBlockResults == nil {
		return ErrNotSupported
	}
	return s.Internal.UserDownloadBlockResults(p0, p1)
}

func (s *LocatorStub) UserDownloadBlockResults(p0 context.Context, p1 []UserBlockDownloadResult) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) AuthNodeNew(p0 context.Context, p1 []auth.Permission, p2 string, p3 string) ([]byte, error) {
	if s.Internal.AuthNodeNew == nil {
		return *new([]byte), ErrNotSupported
	}
	return s.Internal.AuthNodeNew(p0, p1, p2, p3)
}

func (s *SchedulerStub) AuthNodeNew(p0 context.Context, p1 []auth.Permission, p2 string, p3 string) ([]byte, error) {
	return *new([]byte), ErrNotSupported
}

func (s *SchedulerStruct) AuthNodeVerify(p0 context.Context, p1 string) ([]auth.Permission, error) {
	if s.Internal.AuthNodeVerify == nil {
		return *new([]auth.Permission), ErrNotSupported
	}
	return s.Internal.AuthNodeVerify(p0, p1)
}

func (s *SchedulerStub) AuthNodeVerify(p0 context.Context, p1 string) ([]auth.Permission, error) {
	return *new([]auth.Permission), ErrNotSupported
}

func (s *SchedulerStruct) CacheCarfile(p0 context.Context, p1 *CacheCarfileInfo) error {
	if s.Internal.CacheCarfile == nil {
		return ErrNotSupported
	}
	return s.Internal.CacheCarfile(p0, p1)
}

func (s *SchedulerStub) CacheCarfile(p0 context.Context, p1 *CacheCarfileInfo) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) CacheResult(p0 context.Context, p1 CacheResultInfo) error {
	if s.Internal.CacheResult == nil {
		return ErrNotSupported
	}
	return s.Internal.CacheResult(p0, p1)
}

func (s *SchedulerStub) CacheResult(p0 context.Context, p1 CacheResultInfo) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) CandidateNodeConnect(p0 context.Context) error {
	if s.Internal.CandidateNodeConnect == nil {
		return ErrNotSupported
	}
	return s.Internal.CandidateNodeConnect(p0)
}

func (s *SchedulerStub) CandidateNodeConnect(p0 context.Context) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) EdgeNodeConnect(p0 context.Context) error {
	if s.Internal.EdgeNodeConnect == nil {
		return ErrNotSupported
	}
	return s.Internal.EdgeNodeConnect(p0)
}

func (s *SchedulerStub) EdgeNodeConnect(p0 context.Context) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) ElectionValidators(p0 context.Context) error {
	if s.Internal.ElectionValidators == nil {
		return ErrNotSupported
	}
	return s.Internal.ElectionValidators(p0)
}

func (s *SchedulerStub) ElectionValidators(p0 context.Context) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) ExecuteUndoneCarfilesTask(p0 context.Context) error {
	if s.Internal.ExecuteUndoneCarfilesTask == nil {
		return ErrNotSupported
	}
	return s.Internal.ExecuteUndoneCarfilesTask(p0)
}

func (s *SchedulerStub) ExecuteUndoneCarfilesTask(p0 context.Context) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) GetCarfileRecordInfo(p0 context.Context, p1 string) (CarfileRecordInfo, error) {
	if s.Internal.GetCarfileRecordInfo == nil {
		return *new(CarfileRecordInfo), ErrNotSupported
	}
	return s.Internal.GetCarfileRecordInfo(p0, p1)
}

func (s *SchedulerStub) GetCarfileRecordInfo(p0 context.Context, p1 string) (CarfileRecordInfo, error) {
	return *new(CarfileRecordInfo), ErrNotSupported
}

func (s *SchedulerStruct) GetDevicesInfo(p0 context.Context, p1 string) (DevicesInfo, error) {
	if s.Internal.GetDevicesInfo == nil {
		return *new(DevicesInfo), ErrNotSupported
	}
	return s.Internal.GetDevicesInfo(p0, p1)
}

func (s *SchedulerStub) GetDevicesInfo(p0 context.Context, p1 string) (DevicesInfo, error) {
	return *new(DevicesInfo), ErrNotSupported
}

func (s *SchedulerStruct) GetDownloadInfo(p0 context.Context, p1 string) ([]*BlockDownloadInfo, error) {
	if s.Internal.GetDownloadInfo == nil {
		return *new([]*BlockDownloadInfo), ErrNotSupported
	}
	return s.Internal.GetDownloadInfo(p0, p1)
}

func (s *SchedulerStub) GetDownloadInfo(p0 context.Context, p1 string) ([]*BlockDownloadInfo, error) {
	return *new([]*BlockDownloadInfo), ErrNotSupported
}

func (s *SchedulerStruct) GetDownloadInfosWithCarfile(p0 context.Context, p1 string, p2 string) ([]*DownloadInfoResult, error) {
	if s.Internal.GetDownloadInfosWithCarfile == nil {
		return *new([]*DownloadInfoResult), ErrNotSupported
	}
	return s.Internal.GetDownloadInfosWithCarfile(p0, p1, p2)
}

func (s *SchedulerStub) GetDownloadInfosWithCarfile(p0 context.Context, p1 string, p2 string) ([]*DownloadInfoResult, error) {
	return *new([]*DownloadInfoResult), ErrNotSupported
}

func (s *SchedulerStruct) GetExternalIP(p0 context.Context) (string, error) {
	if s.Internal.GetExternalIP == nil {
		return "", ErrNotSupported
	}
	return s.Internal.GetExternalIP(p0)
}

func (s *SchedulerStub) GetExternalIP(p0 context.Context) (string, error) {
	return "", ErrNotSupported
}

func (s *SchedulerStruct) GetNodeAppUpdateInfos(p0 context.Context) (map[int]*NodeAppUpdateInfo, error) {
	if s.Internal.GetNodeAppUpdateInfos == nil {
		return *new(map[int]*NodeAppUpdateInfo), ErrNotSupported
	}
	return s.Internal.GetNodeAppUpdateInfos(p0)
}

func (s *SchedulerStub) GetNodeAppUpdateInfos(p0 context.Context) (map[int]*NodeAppUpdateInfo, error) {
	return *new(map[int]*NodeAppUpdateInfo), ErrNotSupported
}

func (s *SchedulerStruct) GetOnlineDeviceIDs(p0 context.Context, p1 NodeTypeName) ([]string, error) {
	if s.Internal.GetOnlineDeviceIDs == nil {
		return *new([]string), ErrNotSupported
	}
	return s.Internal.GetOnlineDeviceIDs(p0, p1)
}

func (s *SchedulerStub) GetOnlineDeviceIDs(p0 context.Context, p1 NodeTypeName) ([]string, error) {
	return *new([]string), ErrNotSupported
}

func (s *SchedulerStruct) GetPublicKey(p0 context.Context) (string, error) {
	if s.Internal.GetPublicKey == nil {
		return "", ErrNotSupported
	}
	return s.Internal.GetPublicKey(p0)
}

func (s *SchedulerStub) GetPublicKey(p0 context.Context) (string, error) {
	return "", ErrNotSupported
}

func (s *SchedulerStruct) GetRunningCarfileRecords(p0 context.Context) ([]*CarfileRecordInfo, error) {
	if s.Internal.GetRunningCarfileRecords == nil {
		return *new([]*CarfileRecordInfo), ErrNotSupported
	}
	return s.Internal.GetRunningCarfileRecords(p0)
}

func (s *SchedulerStub) GetRunningCarfileRecords(p0 context.Context) ([]*CarfileRecordInfo, error) {
	return *new([]*CarfileRecordInfo), ErrNotSupported
}

func (s *SchedulerStruct) GetUndoneCarfileRecords(p0 context.Context, p1 int) (*DataListInfo, error) {
	if s.Internal.GetUndoneCarfileRecords == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.GetUndoneCarfileRecords(p0, p1)
}

func (s *SchedulerStub) GetUndoneCarfileRecords(p0 context.Context, p1 int) (*DataListInfo, error) {
	return nil, ErrNotSupported
}

func (s *SchedulerStruct) ListCarfileRecords(p0 context.Context, p1 int) (*DataListInfo, error) {
	if s.Internal.ListCarfileRecords == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ListCarfileRecords(p0, p1)
}

func (s *SchedulerStub) ListCarfileRecords(p0 context.Context, p1 int) (*DataListInfo, error) {
	return nil, ErrNotSupported
}

func (s *SchedulerStruct) LocatorConnect(p0 context.Context, p1 string, p2 string) error {
	if s.Internal.LocatorConnect == nil {
		return ErrNotSupported
	}
	return s.Internal.LocatorConnect(p0, p1, p2)
}

func (s *SchedulerStub) LocatorConnect(p0 context.Context, p1 string, p2 string) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) NodeQuit(p0 context.Context, p1 string) error {
	if s.Internal.NodeQuit == nil {
		return ErrNotSupported
	}
	return s.Internal.NodeQuit(p0, p1)
}

func (s *SchedulerStub) NodeQuit(p0 context.Context, p1 string) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) NodeResultForUserDownloadBlock(p0 context.Context, p1 NodeBlockDownloadResult) error {
	if s.Internal.NodeResultForUserDownloadBlock == nil {
		return ErrNotSupported
	}
	return s.Internal.NodeResultForUserDownloadBlock(p0, p1)
}

func (s *SchedulerStub) NodeResultForUserDownloadBlock(p0 context.Context, p1 NodeBlockDownloadResult) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) RegisterNode(p0 context.Context, p1 NodeType, p2 int) ([]NodeRegisterInfo, error) {
	if s.Internal.RegisterNode == nil {
		return *new([]NodeRegisterInfo), ErrNotSupported
	}
	return s.Internal.RegisterNode(p0, p1, p2)
}

func (s *SchedulerStub) RegisterNode(p0 context.Context, p1 NodeType, p2 int) ([]NodeRegisterInfo, error) {
	return *new([]NodeRegisterInfo), ErrNotSupported
}

func (s *SchedulerStruct) RemoveCache(p0 context.Context, p1 string, p2 string) error {
	if s.Internal.RemoveCache == nil {
		return ErrNotSupported
	}
	return s.Internal.RemoveCache(p0, p1, p2)
}

func (s *SchedulerStub) RemoveCache(p0 context.Context, p1 string, p2 string) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) RemoveCarfile(p0 context.Context, p1 string) error {
	if s.Internal.RemoveCarfile == nil {
		return ErrNotSupported
	}
	return s.Internal.RemoveCarfile(p0, p1)
}

func (s *SchedulerStub) RemoveCarfile(p0 context.Context, p1 string) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) RemoveCarfileResult(p0 context.Context, p1 RemoveCarfileResultInfo) error {
	if s.Internal.RemoveCarfileResult == nil {
		return ErrNotSupported
	}
	return s.Internal.RemoveCarfileResult(p0, p1)
}

func (s *SchedulerStub) RemoveCarfileResult(p0 context.Context, p1 RemoveCarfileResultInfo) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) ResetBackupCacheCount(p0 context.Context, p1 int) error {
	if s.Internal.ResetBackupCacheCount == nil {
		return ErrNotSupported
	}
	return s.Internal.ResetBackupCacheCount(p0, p1)
}

func (s *SchedulerStub) ResetBackupCacheCount(p0 context.Context, p1 int) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) ResetCacheExpiredTime(p0 context.Context, p1 string, p2 time.Time) error {
	if s.Internal.ResetCacheExpiredTime == nil {
		return ErrNotSupported
	}
	return s.Internal.ResetCacheExpiredTime(p0, p1, p2)
}

func (s *SchedulerStub) ResetCacheExpiredTime(p0 context.Context, p1 string, p2 time.Time) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) SetNodeAppUpdateInfo(p0 context.Context, p1 *NodeAppUpdateInfo) error {
	if s.Internal.SetNodeAppUpdateInfo == nil {
		return ErrNotSupported
	}
	return s.Internal.SetNodeAppUpdateInfo(p0, p1)
}

func (s *SchedulerStub) SetNodeAppUpdateInfo(p0 context.Context, p1 *NodeAppUpdateInfo) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) StopCacheTask(p0 context.Context, p1 string) error {
	if s.Internal.StopCacheTask == nil {
		return ErrNotSupported
	}
	return s.Internal.StopCacheTask(p0, p1)
}

func (s *SchedulerStub) StopCacheTask(p0 context.Context, p1 string) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) UserDownloadBlockResults(p0 context.Context, p1 []UserBlockDownloadResult) error {
	if s.Internal.UserDownloadBlockResults == nil {
		return ErrNotSupported
	}
	return s.Internal.UserDownloadBlockResults(p0, p1)
}

func (s *SchedulerStub) UserDownloadBlockResults(p0 context.Context, p1 []UserBlockDownloadResult) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) ValidateBlockResult(p0 context.Context, p1 ValidateResults) error {
	if s.Internal.ValidateBlockResult == nil {
		return ErrNotSupported
	}
	return s.Internal.ValidateBlockResult(p0, p1)
}

func (s *SchedulerStub) ValidateBlockResult(p0 context.Context, p1 ValidateResults) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) ValidateRunningState(p0 context.Context) (bool, error) {
	if s.Internal.ValidateRunningState == nil {
		return false, ErrNotSupported
	}
	return s.Internal.ValidateRunningState(p0)
}

func (s *SchedulerStub) ValidateRunningState(p0 context.Context) (bool, error) {
	return false, ErrNotSupported
}

func (s *SchedulerStruct) ValidateStart(p0 context.Context) error {
	if s.Internal.ValidateStart == nil {
		return ErrNotSupported
	}
	return s.Internal.ValidateStart(p0)
}

func (s *SchedulerStub) ValidateStart(p0 context.Context) error {
	return ErrNotSupported
}

func (s *SchedulerStruct) ValidateSwitch(p0 context.Context, p1 bool) error {
	if s.Internal.ValidateSwitch == nil {
		return ErrNotSupported
	}
	return s.Internal.ValidateSwitch(p0, p1)
}

func (s *SchedulerStub) ValidateSwitch(p0 context.Context, p1 bool) error {
	return ErrNotSupported
}

func (s *ValidateStruct) BeValidate(p0 context.Context, p1 ReqValidate, p2 string) error {
	if s.Internal.BeValidate == nil {
		return ErrNotSupported
	}
	return s.Internal.BeValidate(p0, p1, p2)
}

func (s *ValidateStub) BeValidate(p0 context.Context, p1 ReqValidate, p2 string) error {
	return ErrNotSupported
}

func (s *WebStruct) AddCacheTask(p0 context.Context, p1 string, p2 int, p3 time.Time) error {
	if s.Internal.AddCacheTask == nil {
		return ErrNotSupported
	}
	return s.Internal.AddCacheTask(p0, p1, p2, p3)
}

func (s *WebStub) AddCacheTask(p0 context.Context, p1 string, p2 int, p3 time.Time) error {
	return ErrNotSupported
}

func (s *WebStruct) CancelCacheTask(p0 context.Context, p1 string) error {
	if s.Internal.CancelCacheTask == nil {
		return ErrNotSupported
	}
	return s.Internal.CancelCacheTask(p0, p1)
}

func (s *WebStub) CancelCacheTask(p0 context.Context, p1 string) error {
	return ErrNotSupported
}

func (s *WebStruct) GetCacheTaskInfo(p0 context.Context, p1 string) (CarfileRecordInfo, error) {
	if s.Internal.GetCacheTaskInfo == nil {
		return *new(CarfileRecordInfo), ErrNotSupported
	}
	return s.Internal.GetCacheTaskInfo(p0, p1)
}

func (s *WebStub) GetCacheTaskInfo(p0 context.Context, p1 string) (CarfileRecordInfo, error) {
	return *new(CarfileRecordInfo), ErrNotSupported
}

func (s *WebStruct) GetCacheTaskInfos(p0 context.Context, p1 ListCacheInfosReq) (ListCacheInfosRsp, error) {
	if s.Internal.GetCacheTaskInfos == nil {
		return *new(ListCacheInfosRsp), ErrNotSupported
	}
	return s.Internal.GetCacheTaskInfos(p0, p1)
}

func (s *WebStub) GetCacheTaskInfos(p0 context.Context, p1 ListCacheInfosReq) (ListCacheInfosRsp, error) {
	return *new(ListCacheInfosRsp), ErrNotSupported
}

func (s *WebStruct) GetCarfileByCID(p0 context.Context, p1 string) (WebCarfile, error) {
	if s.Internal.GetCarfileByCID == nil {
		return *new(WebCarfile), ErrNotSupported
	}
	return s.Internal.GetCarfileByCID(p0, p1)
}

func (s *WebStub) GetCarfileByCID(p0 context.Context, p1 string) (WebCarfile, error) {
	return *new(WebCarfile), ErrNotSupported
}

func (s *WebStruct) GetNodeInfoByID(p0 context.Context, p1 string) (DevicesInfo, error) {
	if s.Internal.GetNodeInfoByID == nil {
		return *new(DevicesInfo), ErrNotSupported
	}
	return s.Internal.GetNodeInfoByID(p0, p1)
}

func (s *WebStub) GetNodeInfoByID(p0 context.Context, p1 string) (DevicesInfo, error) {
	return *new(DevicesInfo), ErrNotSupported
}

func (s *WebStruct) GetSummaryValidateMessage(p0 context.Context, p1 time.Time, p2 time.Time, p3 int, p4 int) (*SummeryValidateResult, error) {
	if s.Internal.GetSummaryValidateMessage == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.GetSummaryValidateMessage(p0, p1, p2, p3, p4)
}

func (s *WebStub) GetSummaryValidateMessage(p0 context.Context, p1 time.Time, p2 time.Time, p3 int, p4 int) (*SummeryValidateResult, error) {
	return nil, ErrNotSupported
}

func (s *WebStruct) GetSystemInfo(p0 context.Context) (SystemBaseInfo, error) {
	if s.Internal.GetSystemInfo == nil {
		return *new(SystemBaseInfo), ErrNotSupported
	}
	return s.Internal.GetSystemInfo(p0)
}

func (s *WebStub) GetSystemInfo(p0 context.Context) (SystemBaseInfo, error) {
	return *new(SystemBaseInfo), ErrNotSupported
}

func (s *WebStruct) ListBlockDownloadInfo(p0 context.Context, p1 ListBlockDownloadInfoReq) (ListBlockDownloadInfoRsp, error) {
	if s.Internal.ListBlockDownloadInfo == nil {
		return *new(ListBlockDownloadInfoRsp), ErrNotSupported
	}
	return s.Internal.ListBlockDownloadInfo(p0, p1)
}

func (s *WebStub) ListBlockDownloadInfo(p0 context.Context, p1 ListBlockDownloadInfoReq) (ListBlockDownloadInfoRsp, error) {
	return *new(ListBlockDownloadInfoRsp), ErrNotSupported
}

func (s *WebStruct) ListCacheTasks(p0 context.Context, p1 int, p2 int) (ListCacheTasksRsp, error) {
	if s.Internal.ListCacheTasks == nil {
		return *new(ListCacheTasksRsp), ErrNotSupported
	}
	return s.Internal.ListCacheTasks(p0, p1, p2)
}

func (s *WebStub) ListCacheTasks(p0 context.Context, p1 int, p2 int) (ListCacheTasksRsp, error) {
	return *new(ListCacheTasksRsp), ErrNotSupported
}

func (s *WebStruct) ListNodeConnectionLog(p0 context.Context, p1 ListNodeConnectionLogReq) (ListNodeConnectionLogRsp, error) {
	if s.Internal.ListNodeConnectionLog == nil {
		return *new(ListNodeConnectionLogRsp), ErrNotSupported
	}
	return s.Internal.ListNodeConnectionLog(p0, p1)
}

func (s *WebStub) ListNodeConnectionLog(p0 context.Context, p1 ListNodeConnectionLogReq) (ListNodeConnectionLogRsp, error) {
	return *new(ListNodeConnectionLogRsp), ErrNotSupported
}

func (s *WebStruct) ListNodes(p0 context.Context, p1 int, p2 int) (ListNodesRsp, error) {
	if s.Internal.ListNodes == nil {
		return *new(ListNodesRsp), ErrNotSupported
	}
	return s.Internal.ListNodes(p0, p1, p2)
}

func (s *WebStub) ListNodes(p0 context.Context, p1 int, p2 int) (ListNodesRsp, error) {
	return *new(ListNodesRsp), ErrNotSupported
}

func (s *WebStruct) ListValidateResult(p0 context.Context, p1 int, p2 int) (ListValidateResultRsp, error) {
	if s.Internal.ListValidateResult == nil {
		return *new(ListValidateResultRsp), ErrNotSupported
	}
	return s.Internal.ListValidateResult(p0, p1, p2)
}

func (s *WebStub) ListValidateResult(p0 context.Context, p1 int, p2 int) (ListValidateResultRsp, error) {
	return *new(ListValidateResultRsp), ErrNotSupported
}

func (s *WebStruct) RemoveCarfile(p0 context.Context, p1 string) error {
	if s.Internal.RemoveCarfile == nil {
		return ErrNotSupported
	}
	return s.Internal.RemoveCarfile(p0, p1)
}

func (s *WebStub) RemoveCarfile(p0 context.Context, p1 string) error {
	return ErrNotSupported
}

func (s *WebStruct) SetupValidation(p0 context.Context, p1 bool) error {
	if s.Internal.SetupValidation == nil {
		return ErrNotSupported
	}
	return s.Internal.SetupValidation(p0, p1)
}

func (s *WebStub) SetupValidation(p0 context.Context, p1 bool) error {
	return ErrNotSupported
}

var _ Candidate = new(CandidateStruct)
var _ CarfileOperation = new(CarfileOperationStruct)
var _ Common = new(CommonStruct)
var _ DataSync = new(DataSyncStruct)
var _ Device = new(DeviceStruct)
var _ Download = new(DownloadStruct)
var _ Edge = new(EdgeStruct)
var _ Locator = new(LocatorStruct)
var _ Scheduler = new(SchedulerStruct)
var _ Validate = new(ValidateStruct)
var _ Web = new(WebStruct)
