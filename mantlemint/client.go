package mantlemint

import (
	abcicli "github.com/cometbft/cometbft/abci/client"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/libs/service"
	tmsync "github.com/cometbft/cometbft/libs/sync"
)

// NewConcurrentQueryClient creates a local client, which will be directly calling the
// methods of the given app. + uses RWMutex for reads
func NewConcurrentQueryClient(mtx *tmsync.RWMutex, app types.Application) abcicli.Client {
	if mtx == nil {
		mtx = &tmsync.RWMutex{}
	}

	cli := &localClient{
		mtx:         mtx,
		Application: app,
	}

	cli.BaseService = *service.NewBaseService(nil, "localClient", cli)
	return cli
}

var _ abcicli.Client = (*localClient)(nil)

// NOTE: use defer to unlock mutex because Application might panic (e.g., in
// case of malicious tx or query). It only makes sense for publicly exposed
// methods like CheckTx (/broadcast_tx_* RPC endpoint) or Query (/abci_query
// RPC endpoint), but defers are used everywhere for the sake of consistency.
type localClient struct {
	service.BaseService

	mtx *tmsync.RWMutex
	types.Application
	abcicli.Callback
}

// IsRunning implements abcicli.Client.
// Subtle: this method shadows the method (BaseService).IsRunning of localClient.BaseService.
func (app *localClient) IsRunning() bool {
	return app.BaseService.IsRunning()
}

// OnReset implements abcicli.Client.
// Subtle: this method shadows the method (BaseService).OnReset of localClient.BaseService.
func (app *localClient) OnReset() error {
	return app.BaseService.OnReset()
}

// OnStart implements abcicli.Client.
// Subtle: this method shadows the method (BaseService).OnStart of localClient.BaseService.
func (app *localClient) OnStart() error {
	return app.BaseService.OnStart()
}

// OnStop implements abcicli.Client.
// Subtle: this method shadows the method (BaseService).OnStop of localClient.BaseService.
func (app *localClient) OnStop() {
	app.BaseService.OnStop()
}

// PrepareProposalAsync implements abcicli.Client.
func (app *localClient) PrepareProposalAsync(req types.RequestPrepareProposal) *abcicli.ReqRes {
	res := app.Application.PrepareProposal(req)
	return newLocalReqRes(
		types.ToRequestPrepareProposal(req),
		types.ToResponsePrepareProposal(res),
	)
}

// PrepareProposalSync implements abcicli.Client.
func (app *localClient) PrepareProposalSync(req types.RequestPrepareProposal) (*types.ResponsePrepareProposal, error) {
	res := app.Application.PrepareProposal(req)
	return &res, nil
}

// ProcessProposalAsync implements abcicli.Client.
func (app *localClient) ProcessProposalAsync(req types.RequestProcessProposal) *abcicli.ReqRes {
	res := app.Application.ProcessProposal(req)
	return newLocalReqRes(
		types.ToRequestProcessProposal(req),
		types.ToResponseProcessProposal(res),
	)
}

// ProcessProposalSync implements abcicli.Client.
func (app *localClient) ProcessProposalSync(req types.RequestProcessProposal) (*types.ResponseProcessProposal, error) {
	res := app.Application.ProcessProposal(req)
	return &res, nil
}

// Quit implements abcicli.Client.
// Subtle: this method shadows the method (BaseService).Quit of localClient.BaseService.
func (app *localClient) Quit() <-chan struct{} {
	res := app.BaseService.Quit()
	return res
}

// Reset implements abcicli.Client.
// Subtle: this method shadows the method (BaseService).Reset of localClient.BaseService.
func (app *localClient) Reset() error {
	return app.BaseService.Reset()
}

// SetLogger implements abcicli.Client.
// Subtle: this method shadows the method (BaseService).SetLogger of localClient.BaseService.
func (app *localClient) SetLogger(logger log.Logger) {
	app.BaseService.SetLogger(logger)
}

// Start implements abcicli.Client.
// Subtle: this method shadows the method (BaseService).Start of localClient.BaseService.
func (app *localClient) Start() error {
	return app.BaseService.Start()
}

// Stop implements abcicli.Client.
// Subtle: this method shadows the method (BaseService).Stop of localClient.BaseService.
func (app *localClient) Stop() error {
	return app.BaseService.Stop()
}

// String implements abcicli.Client.
// Subtle: this method shadows the method (BaseService).String of localClient.BaseService.
func (app *localClient) String() string {
	return app.BaseService.String()
}

func (app *localClient) SetResponseCallback(cb abcicli.Callback) {
	app.mtx.Lock()
	defer app.mtx.Unlock()
	app.Callback = cb
}

// TODO: change types.Application to include Error()?
func (app *localClient) Error() error {
	return nil
}

func (app *localClient) FlushAsync() *abcicli.ReqRes {
	// Do nothing
	return newLocalReqRes(types.ToRequestFlush(), nil)
}

func (app *localClient) EchoAsync(msg string) *abcicli.ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	return app.callback(
		types.ToRequestEcho(msg),
		types.ToResponseEcho(msg),
	)
}

func (app *localClient) InfoAsync(req types.RequestInfo) *abcicli.ReqRes {
	app.mtx.RLock()
	defer app.mtx.RUnlock()

	res := app.Application.Info(req)
	return app.callback(
		types.ToRequestInfo(req),
		types.ToResponseInfo(res),
	)
}

// func (app *localClient) SetOptionAsync(req types.RequestSetOption) *abcicli.ReqRes {
// 	app.mtx.Lock()
// 	defer app.mtx.Unlock()

// 	res := app.Application.SetOption(req)
// 	return app.callback(
// 		types.ToRequestSetOption(req),
// 		types.ToResponseSetOption(res),
// 	)
// }

func (app *localClient) DeliverTxAsync(params types.RequestDeliverTx) *abcicli.ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.DeliverTx(params)
	return app.callback(
		types.ToRequestDeliverTx(params),
		types.ToResponseDeliverTx(res),
	)
}

func (app *localClient) CheckTxAsync(req types.RequestCheckTx) *abcicli.ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.CheckTx(req)
	return app.callback(
		types.ToRequestCheckTx(req),
		types.ToResponseCheckTx(res),
	)
}

func (app *localClient) QueryAsync(req types.RequestQuery) *abcicli.ReqRes {
	res := app.Application.Query(req)
	return app.callback(
		types.ToRequestQuery(req),
		types.ToResponseQuery(res),
	)
}

func (app *localClient) CommitAsync() *abcicli.ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.Commit()
	return app.callback(
		types.ToRequestCommit(),
		types.ToResponseCommit(res),
	)
}

func (app *localClient) InitChainAsync(req types.RequestInitChain) *abcicli.ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.InitChain(req)
	return app.callback(
		types.ToRequestInitChain(req),
		types.ToResponseInitChain(res),
	)
}

func (app *localClient) BeginBlockAsync(req types.RequestBeginBlock) *abcicli.ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.BeginBlock(req)
	return app.callback(
		types.ToRequestBeginBlock(req),
		types.ToResponseBeginBlock(res),
	)
}

func (app *localClient) EndBlockAsync(req types.RequestEndBlock) *abcicli.ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.EndBlock(req)
	return app.callback(
		types.ToRequestEndBlock(req),
		types.ToResponseEndBlock(res),
	)
}

func (app *localClient) ListSnapshotsAsync(req types.RequestListSnapshots) *abcicli.ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.ListSnapshots(req)
	return app.callback(
		types.ToRequestListSnapshots(req),
		types.ToResponseListSnapshots(res),
	)
}

func (app *localClient) OfferSnapshotAsync(req types.RequestOfferSnapshot) *abcicli.ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.OfferSnapshot(req)
	return app.callback(
		types.ToRequestOfferSnapshot(req),
		types.ToResponseOfferSnapshot(res),
	)
}

func (app *localClient) LoadSnapshotChunkAsync(req types.RequestLoadSnapshotChunk) *abcicli.ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.LoadSnapshotChunk(req)
	return app.callback(
		types.ToRequestLoadSnapshotChunk(req),
		types.ToResponseLoadSnapshotChunk(res),
	)
}

func (app *localClient) ApplySnapshotChunkAsync(req types.RequestApplySnapshotChunk) *abcicli.ReqRes {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.ApplySnapshotChunk(req)
	return app.callback(
		types.ToRequestApplySnapshotChunk(req),
		types.ToResponseApplySnapshotChunk(res),
	)
}

//-------------------------------------------------------

func (app *localClient) FlushSync() error {
	return nil
}

func (app *localClient) EchoSync(msg string) (*types.ResponseEcho, error) {
	return &types.ResponseEcho{Message: msg}, nil
}

func (app *localClient) InfoSync(req types.RequestInfo) (*types.ResponseInfo, error) {
	app.mtx.RLock()
	defer app.mtx.RUnlock()

	res := app.Application.Info(req)
	return &res, nil
}

// func (app *localClient) SetOptionSync(req types.RequestSetOption) (*types.ResponseSetOption, error) {
// 	app.mtx.Lock()
// 	defer app.mtx.Unlock()

// 	res := app.Application.SetOption(req)
// 	return &res, nil
// }

func (app *localClient) DeliverTxSync(req types.RequestDeliverTx) (*types.ResponseDeliverTx, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.DeliverTx(req)
	return &res, nil
}

func (app *localClient) CheckTxSync(req types.RequestCheckTx) (*types.ResponseCheckTx, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.CheckTx(req)
	return &res, nil
}

func (app *localClient) QuerySync(req types.RequestQuery) (*types.ResponseQuery, error) {
	res := app.Application.Query(req)
	return &res, nil
}

func (app *localClient) CommitSync() (*types.ResponseCommit, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.Commit()
	return &res, nil
}

func (app *localClient) InitChainSync(req types.RequestInitChain) (*types.ResponseInitChain, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.InitChain(req)
	return &res, nil
}

func (app *localClient) BeginBlockSync(req types.RequestBeginBlock) (*types.ResponseBeginBlock, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.BeginBlock(req)
	return &res, nil
}

func (app *localClient) EndBlockSync(req types.RequestEndBlock) (*types.ResponseEndBlock, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.EndBlock(req)
	return &res, nil
}

func (app *localClient) ListSnapshotsSync(req types.RequestListSnapshots) (*types.ResponseListSnapshots, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.ListSnapshots(req)
	return &res, nil
}

func (app *localClient) OfferSnapshotSync(req types.RequestOfferSnapshot) (*types.ResponseOfferSnapshot, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.OfferSnapshot(req)
	return &res, nil
}

func (app *localClient) LoadSnapshotChunkSync(
	req types.RequestLoadSnapshotChunk) (*types.ResponseLoadSnapshotChunk, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.LoadSnapshotChunk(req)
	return &res, nil
}

func (app *localClient) ApplySnapshotChunkSync(
	req types.RequestApplySnapshotChunk) (*types.ResponseApplySnapshotChunk, error) {
	app.mtx.Lock()
	defer app.mtx.Unlock()

	res := app.Application.ApplySnapshotChunk(req)
	return &res, nil
}

//-------------------------------------------------------

func (app *localClient) callback(req *types.Request, res *types.Response) *abcicli.ReqRes {
	app.Callback(req, res)
	return newLocalReqRes(req, res)
}

func newLocalReqRes(req *types.Request, res *types.Response) *abcicli.ReqRes {
	reqRes := abcicli.NewReqRes(req)
	reqRes.Response = res
	return reqRes
}
