package services

import (
	"github.com/iotaledger/wasp/packages/chain/chaintypes"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/metrics"
	"github.com/iotaledger/wasp/packages/webapi/dto"
	"github.com/iotaledger/wasp/packages/webapi/interfaces"
	"github.com/iotaledger/wasp/packages/webapi/models"
)

type MetricsService struct {
	chainProvider        chaintypes.ChainsProvider
	chainMetricsProvider *metrics.ChainMetricsProvider
}

func NewMetricsService(chainProvider chaintypes.ChainsProvider, chainMetricsProvider *metrics.ChainMetricsProvider) interfaces.MetricsService {
	return &MetricsService{
		chainProvider:        chainProvider,
		chainMetricsProvider: chainMetricsProvider,
	}
}

func (c *MetricsService) GetNodeMessageMetrics() *dto.NodeMessageMetrics {
	return &dto.NodeMessageMetrics{
		RegisteredChainIDs: c.chainMetricsProvider.RegisteredChains(),

		InStateOutput:      dto.MapMetricItem(c.chainMetricsProvider.Message.InStateOutput()),
		InAnchorOutput:     dto.MapMetricItem(c.chainMetricsProvider.Message.InAnchorOutput()),
		InOutput:           dto.MapMetricItem(c.chainMetricsProvider.Message.InOutput()),
		InOnLedgerRequest:  dto.MapMetricItem(c.chainMetricsProvider.Message.InOnLedgerRequest()),
		InTxInclusionState: dto.MapMetricItem(c.chainMetricsProvider.Message.InTxInclusionState()),

		OutPublishStateTransaction:      dto.MapMetricItem(c.chainMetricsProvider.Message.OutPublishStateTransaction()),
		OutPublishGovernanceTransaction: dto.MapMetricItem(c.chainMetricsProvider.Message.OutPublishGovernanceTransaction()),
		OutPullLatestOutput:             dto.MapMetricItem(c.chainMetricsProvider.Message.OutPullLatestOutput()),
		OutPullTxInclusionState:         dto.MapMetricItem(c.chainMetricsProvider.Message.OutPullTxInclusionState()),
		OutPullOutputByID:               dto.MapMetricItem(c.chainMetricsProvider.Message.OutPullOutputByID()),
	}
}

func (c *MetricsService) GetChainMessageMetrics(chainID isc.ChainID) *dto.ChainMessageMetrics {
	chain, err := c.chainProvider().Get(chainID)
	if err != nil {
		return nil
	}

	chainMetrics := chain.GetChainMetrics()

	return &dto.ChainMessageMetrics{
		InStateOutput:      dto.MapMetricItem(chainMetrics.Message.InStateOutput()),
		InAnchorOutput:     dto.MapMetricItem(chainMetrics.Message.InAnchorOutput()),
		InOutput:           dto.MapMetricItem(chainMetrics.Message.InOutput()),
		InOnLedgerRequest:  dto.MapMetricItem(chainMetrics.Message.InOnLedgerRequest()),
		InTxInclusionState: dto.MapMetricItem(chainMetrics.Message.InTxInclusionState()),

		OutPublishStateTransaction:      dto.MapMetricItem(chainMetrics.Message.OutPublishStateTransaction()),
		OutPublishGovernanceTransaction: dto.MapMetricItem(chainMetrics.Message.OutPublishGovernanceTransaction()),
		OutPullLatestOutput:             dto.MapMetricItem(chainMetrics.Message.OutPullLatestOutput()),
		OutPullTxInclusionState:         dto.MapMetricItem(chainMetrics.Message.OutPullTxInclusionState()),
		OutPullOutputByID:               dto.MapMetricItem(chainMetrics.Message.OutPullOutputByID()),
	}
}

func (c *MetricsService) GetChainConsensusWorkflowMetrics(chainID isc.ChainID) *models.ConsensusWorkflowMetrics {
	chain, err := c.chainProvider().Get(chainID)
	if err != nil {
		return nil
	}

	metrics := chain.GetConsensusWorkflowStatus()
	if metrics == nil {
		return nil
	}

	return models.MapConsensusWorkflowStatus(metrics)
}

func (c *MetricsService) GetChainConsensusPipeMetrics(chainID isc.ChainID) *models.ConsensusPipeMetrics {
	chain, err := c.chainProvider().Get(chainID)
	if err != nil {
		return nil
	}

	metrics := chain.GetConsensusPipeMetrics()
	if metrics == nil {
		return nil
	}

	return models.MapConsensusPipeMetrics(metrics)
}

func (c *MetricsService) GetMaxChainConfirmedStateLag() uint32 {
	return c.chainMetricsProvider.StateManager.MaxChainConfirmedStateLag()
}
