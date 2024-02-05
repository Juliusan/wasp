package corecontracts

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/vm/core/blocklog"
	"github.com/iotaledger/wasp/packages/webapi/apierrors"
	"github.com/iotaledger/wasp/packages/webapi/common"
	"github.com/iotaledger/wasp/packages/webapi/controllers/controllerutils"
	"github.com/iotaledger/wasp/packages/webapi/corecontracts"
	"github.com/iotaledger/wasp/packages/webapi/interfaces"
	"github.com/iotaledger/wasp/packages/webapi/models"
	"github.com/iotaledger/wasp/packages/webapi/params"
)

func (c *Controller) getControlAddresses(e echo.Context) error {
	ch, chainID, err := controllerutils.ChainFromParams(e, c.chainService, c.l1Api)
	if err != nil {
		return c.handleViewCallError(err, chainID)
	}
	controlAddresses, err := corecontracts.GetControlAddresses(ch)
	if err != nil {
		return c.handleViewCallError(err, chainID)
	}

	controlAddressesResponse := &models.ControlAddressesResponse{
		GoverningAddress: controlAddresses.GoverningAddress.Bech32(c.l1Api.ProtocolParameters().Bech32HRP()),
		SinceBlockIndex:  controlAddresses.SinceBlockIndex,
		StateAddress:     controlAddresses.StateAddress.Bech32(c.l1Api.ProtocolParameters().Bech32HRP()),
	}

	return e.JSON(http.StatusOK, controlAddressesResponse)
}

func (c *Controller) getBlockInfo(e echo.Context) error {
	invoker, ch, err := c.createCallViewInvoker(e)
	if err != nil {
		return err
	}

	var blockInfo *blocklog.BlockInfo
	blockIndex := e.Param(params.ParamBlockIndex)

	if blockIndex == "" {
		blockInfo, err = corecontracts.GetLatestBlockInfo(invoker, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	} else {
		var blockIndexNum uint64
		blockIndexNum, err = strconv.ParseUint(e.Param(params.ParamBlockIndex), 10, 64)
		if err != nil {
			return apierrors.InvalidPropertyError(params.ParamBlockIndex, err)
		}

		blockInfo, err = corecontracts.GetBlockInfo(invoker, uint32(blockIndexNum), e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	}
	if err != nil {
		return c.handleViewCallError(err, ch.ID())
	}

	blockInfoResponse := models.MapBlockInfoResponse(blockInfo)

	return e.JSON(http.StatusOK, blockInfoResponse)
}

func (c *Controller) getRequestIDsForBlock(e echo.Context) error {
	invoker, ch, err := c.createCallViewInvoker(e)
	if err != nil {
		return err
	}

	var requestIDs []isc.RequestID
	blockIndex := e.Param(params.ParamBlockIndex)

	if blockIndex == "" {
		requestIDs, err = corecontracts.GetRequestIDsForLatestBlock(invoker, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	} else {
		var blockIndexNum uint64
		blockIndexNum, err = params.DecodeUInt(e, params.ParamBlockIndex)
		if err != nil {
			return err
		}

		requestIDs, err = corecontracts.GetRequestIDsForBlock(invoker, uint32(blockIndexNum), e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	}

	if err != nil {
		return c.handleViewCallError(err, ch.ID())
	}

	requestIDsResponse := &models.RequestIDsResponse{
		RequestIDs: make([]string, len(requestIDs)),
	}

	for k, v := range requestIDs {
		requestIDsResponse.RequestIDs[k] = v.String()
	}

	return e.JSON(http.StatusOK, requestIDsResponse)
}

func GetRequestReceipt(e echo.Context, c interfaces.ChainService, l1API iotago.API) error {
	ch, _, err := controllerutils.ChainFromParams(e, c, l1API)
	if err != nil {
		return err
	}
	requestID, err := params.DecodeRequestID(e)
	if err != nil {
		return err
	}

	invoker := corecontracts.MakeCallViewInvoker(ch)
	receipt, ok, err := corecontracts.GetRequestReceipt(invoker, requestID, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	if err != nil {
		panic(err)
	}
	if !ok {
		return apierrors.NoRecordFoundError(errors.New("no receipt"))
	}

	resolvedReceipt, err := common.ParseReceipt(ch, receipt)
	if err != nil {
		panic(err)
	}

	return e.JSON(http.StatusOK, models.MapReceiptResponse(l1API, resolvedReceipt))
}

func (c *Controller) getRequestReceipt(e echo.Context) error {
	return GetRequestReceipt(e, c.chainService, c.l1Api)
}

func (c *Controller) getRequestReceiptsForBlock(e echo.Context) error {
	invoker, ch, err := c.createCallViewInvoker(e)
	if err != nil {
		return err
	}

	var blocklogReceipts []*blocklog.RequestReceipt
	blockIndex := e.Param(params.ParamBlockIndex)

	if blockIndex == "" {
		var blockInfo *blocklog.BlockInfo
		blockInfo, err = corecontracts.GetLatestBlockInfo(invoker, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
		if err != nil {
			return c.handleViewCallError(err, ch.ID())
		}

		blocklogReceipts, err = corecontracts.GetRequestReceiptsForBlock(invoker, blockInfo.BlockIndex(), e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	} else {
		var blockIndexNum uint64
		blockIndexNum, err = params.DecodeUInt(e, params.ParamBlockIndex)
		if err != nil {
			return err
		}

		blocklogReceipts, err = corecontracts.GetRequestReceiptsForBlock(invoker, uint32(blockIndexNum), e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	}
	if err != nil {
		return c.handleViewCallError(err, ch.ID())
	}

	receiptsResponse := make([]*models.ReceiptResponse, len(blocklogReceipts))

	for i, blocklogReceipt := range blocklogReceipts {
		parsedReceipt, err := common.ParseReceipt(ch, blocklogReceipt)
		if err != nil {
			panic(err)
		}
		receiptResp := models.MapReceiptResponse(c.l1Api, parsedReceipt)
		receiptsResponse[i] = receiptResp
	}

	return e.JSON(http.StatusOK, receiptsResponse)
}

func (c *Controller) getIsRequestProcessed(e echo.Context) error {
	invoker, ch, err := c.createCallViewInvoker(e)
	if err != nil {
		return err
	}

	requestID, err := params.DecodeRequestID(e)
	if err != nil {
		return err
	}

	requestProcessed, err := corecontracts.IsRequestProcessed(invoker, requestID, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	if err != nil {
		return c.handleViewCallError(err, ch.ID())
	}

	requestProcessedResponse := models.RequestProcessedResponse{
		ChainID:     ch.ID().Bech32(c.l1Api.ProtocolParameters().Bech32HRP()),
		RequestID:   requestID.String(),
		IsProcessed: requestProcessed,
	}

	return e.JSON(http.StatusOK, requestProcessedResponse)
}

func eventsResponse(e echo.Context, events []*isc.Event) error {
	eventsJSON := make([]*isc.EventJSON, len(events))
	for i, ev := range events {
		eventsJSON[i] = ev.ToJSONStruct()
	}
	return e.JSON(http.StatusOK, models.EventsResponse{
		Events: eventsJSON,
	})
}

func (c *Controller) getBlockEvents(e echo.Context) error {
	invoker, ch, err := c.createCallViewInvoker(e)
	if err != nil {
		return err
	}

	var events []*isc.Event
	blockIndex := e.Param(params.ParamBlockIndex)

	if blockIndex != "" {
		blockIndexNum, err := params.DecodeUInt(e, params.ParamBlockIndex)
		if err != nil {
			return err
		}

		events, err = corecontracts.GetEventsForBlock(invoker, uint32(blockIndexNum), e.QueryParam(params.ParamBlockIndexOrTrieRoot))
		if err != nil {
			return c.handleViewCallError(err, ch.ID())
		}
	} else {
		blockInfo, err := corecontracts.GetLatestBlockInfo(invoker, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
		if err != nil {
			return c.handleViewCallError(err, ch.ID())
		}

		events, err = corecontracts.GetEventsForBlock(invoker, blockInfo.BlockIndex(), e.QueryParam(params.ParamBlockIndexOrTrieRoot))
		if err != nil {
			return c.handleViewCallError(err, ch.ID())
		}
	}
	return eventsResponse(e, events)
}

func (c *Controller) getContractEvents(e echo.Context) error {
	invoker, ch, err := c.createCallViewInvoker(e)
	if err != nil {
		return err
	}

	contractHname, err := params.DecodeHNameFromHNameHexString(e, "contractHname")
	if err != nil {
		return err
	}

	events, err := corecontracts.GetEventsForContract(invoker, contractHname, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	if err != nil {
		return c.handleViewCallError(err, ch.ID())
	}
	return eventsResponse(e, events)
}

func (c *Controller) getRequestEvents(e echo.Context) error {
	invoker, ch, err := c.createCallViewInvoker(e)
	if err != nil {
		return err
	}

	requestID, err := params.DecodeRequestID(e)
	if err != nil {
		return err
	}

	events, err := corecontracts.GetEventsForRequest(invoker, requestID, e.QueryParam(params.ParamBlockIndexOrTrieRoot))
	if err != nil {
		return c.handleViewCallError(err, ch.ID())
	}
	return eventsResponse(e, events)
}
