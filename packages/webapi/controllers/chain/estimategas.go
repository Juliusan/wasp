package chain

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/iota.go/v4/hexutil"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/webapi/apierrors"
	"github.com/iotaledger/wasp/packages/webapi/common"
	"github.com/iotaledger/wasp/packages/webapi/controllers/controllerutils"
	"github.com/iotaledger/wasp/packages/webapi/models"
)

func (c *Controller) estimateGasOnLedger(e echo.Context) error {
	controllerutils.SetOperation(e, "estimate_gas_onledger")
	ch, chainID, err := controllerutils.ChainFromParams(e, c.chainService, c.l1API)
	if err != nil {
		return err
	}

	var estimateGasRequest models.EstimateGasRequestOnledger
	if err = e.Bind(&estimateGasRequest); err != nil {
		return apierrors.InvalidPropertyError("body", err)
	}

	outputBytes, err := hexutil.DecodeHex(estimateGasRequest.Output)
	if err != nil {
		return apierrors.InvalidPropertyError("Request", err)
	}
	output, err := util.OutputFromBytes(outputBytes, c.l1API)
	if err != nil {
		return apierrors.InvalidPropertyError("Output", err)
	}

	req, err := isc.OnLedgerFromUTXO(
		output,
		iotago.OutputID{}, // empty outputID for estimation
	)
	if err != nil {
		return apierrors.InvalidPropertyError("Output", err)
	}
	if !req.TargetAddress().Equal(chainID.AsAddress()) {
		return apierrors.InvalidPropertyError("Request", errors.New("wrong chainID"))
	}

	rec, err := common.EstimateGas(ch, req)
	if err != nil {
		return apierrors.NewHTTPError(http.StatusBadRequest, "VM run error", err)
	}
	return e.JSON(http.StatusOK, models.MapReceiptResponse(c.l1API, rec))
}

func (c *Controller) estimateGasOffLedger(e echo.Context) error {
	controllerutils.SetOperation(e, "estimate_gas_offledger")
	ch, chainID, err := controllerutils.ChainFromParams(e, c.chainService, c.l1API)
	if err != nil {
		return err
	}

	var estimateGasRequest models.EstimateGasRequestOffledger
	if err = e.Bind(&estimateGasRequest); err != nil {
		return apierrors.InvalidPropertyError("body", err)
	}

	requestBytes, err := hexutil.DecodeHex(estimateGasRequest.Request)
	if err != nil {
		return apierrors.InvalidPropertyError("Request", err)
	}

	req, err := c.offLedgerService.ParseRequest(requestBytes)
	if err != nil {
		return apierrors.InvalidPropertyError("Request", err)
	}
	if !req.TargetAddress().Equal(chainID.AsAddress()) {
		return apierrors.InvalidPropertyError("Request", errors.New("wrong chainID"))
	}

	rec, err := common.EstimateGas(ch, req)
	if err != nil {
		return apierrors.NewHTTPError(http.StatusBadRequest, "VM run error", err)
	}
	return e.JSON(http.StatusOK, models.MapReceiptResponse(c.l1API, rec))
}
