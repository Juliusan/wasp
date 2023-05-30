// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

import * as wasmlib from '../index';
import * as sc from './index';

export class AddAllowedStateControllerAddressCall {
    func:   wasmlib.ScFunc;
    params: sc.MutableAddAllowedStateControllerAddressParams = new sc.MutableAddAllowedStateControllerAddressParams(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScFuncClientContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncAddAllowedStateControllerAddress);
    }
}

export class AddCandidateNodeCall {
    func:   wasmlib.ScFunc;
    params: sc.MutableAddCandidateNodeParams = new sc.MutableAddCandidateNodeParams(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScFuncClientContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncAddCandidateNode);
    }
}

export class ChangeAccessNodesCall {
    func:   wasmlib.ScFunc;
    params: sc.MutableChangeAccessNodesParams = new sc.MutableChangeAccessNodesParams(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScFuncClientContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncChangeAccessNodes);
    }
}

export class ClaimChainOwnershipCall {
    func: wasmlib.ScFunc;

    public constructor(ctx: wasmlib.ScFuncClientContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncClaimChainOwnership);
    }
}

export class DelegateChainOwnershipCall {
    func:   wasmlib.ScFunc;
    params: sc.MutableDelegateChainOwnershipParams = new sc.MutableDelegateChainOwnershipParams(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScFuncClientContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncDelegateChainOwnership);
    }
}

export class RemoveAllowedStateControllerAddressCall {
    func:   wasmlib.ScFunc;
    params: sc.MutableRemoveAllowedStateControllerAddressParams = new sc.MutableRemoveAllowedStateControllerAddressParams(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScFuncClientContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncRemoveAllowedStateControllerAddress);
    }
}

export class RevokeAccessNodeCall {
    func:   wasmlib.ScFunc;
    params: sc.MutableRevokeAccessNodeParams = new sc.MutableRevokeAccessNodeParams(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScFuncClientContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncRevokeAccessNode);
    }
}

export class RotateStateControllerCall {
    func:   wasmlib.ScFunc;
    params: sc.MutableRotateStateControllerParams = new sc.MutableRotateStateControllerParams(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScFuncClientContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncRotateStateController);
    }
}

export class SetEVMGasRatioCall {
    func:   wasmlib.ScFunc;
    params: sc.MutableSetEVMGasRatioParams = new sc.MutableSetEVMGasRatioParams(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScFuncClientContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncSetEVMGasRatio);
    }
}

export class SetFeePolicyCall {
    func:   wasmlib.ScFunc;
    params: sc.MutableSetFeePolicyParams = new sc.MutableSetFeePolicyParams(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScFuncClientContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncSetFeePolicy);
    }
}

export class SetGasLimitsCall {
    func:   wasmlib.ScFunc;
    params: sc.MutableSetGasLimitsParams = new sc.MutableSetGasLimitsParams(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScFuncClientContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncSetGasLimits);
    }
}

export class SetMetadataCall {
    func:   wasmlib.ScFunc;
    params: sc.MutableSetMetadataParams = new sc.MutableSetMetadataParams(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScFuncClientContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncSetMetadata);
    }
}

export class StartMaintenanceCall {
    func: wasmlib.ScFunc;

    public constructor(ctx: wasmlib.ScFuncClientContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncStartMaintenance);
    }
}

export class StopMaintenanceCall {
    func: wasmlib.ScFunc;

    public constructor(ctx: wasmlib.ScFuncClientContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncStopMaintenance);
    }
}

export class GetAllowedStateControllerAddressesCall {
    func:    wasmlib.ScView;
    results: sc.ImmutableGetAllowedStateControllerAddressesResults = new sc.ImmutableGetAllowedStateControllerAddressesResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewClientContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewGetAllowedStateControllerAddresses);
    }
}

export class GetChainInfoCall {
    func:    wasmlib.ScView;
    results: sc.ImmutableGetChainInfoResults = new sc.ImmutableGetChainInfoResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewClientContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewGetChainInfo);
    }
}

export class GetChainNodesCall {
    func:    wasmlib.ScView;
    results: sc.ImmutableGetChainNodesResults = new sc.ImmutableGetChainNodesResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewClientContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewGetChainNodes);
    }
}

export class GetChainOwnerCall {
    func:    wasmlib.ScView;
    results: sc.ImmutableGetChainOwnerResults = new sc.ImmutableGetChainOwnerResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewClientContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewGetChainOwner);
    }
}

export class GetEVMGasRatioCall {
    func:    wasmlib.ScView;
    results: sc.ImmutableGetEVMGasRatioResults = new sc.ImmutableGetEVMGasRatioResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewClientContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewGetEVMGasRatio);
    }
}

export class GetFeePolicyCall {
    func:    wasmlib.ScView;
    results: sc.ImmutableGetFeePolicyResults = new sc.ImmutableGetFeePolicyResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewClientContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewGetFeePolicy);
    }
}

export class GetGasLimitsCall {
    func:    wasmlib.ScView;
    results: sc.ImmutableGetGasLimitsResults = new sc.ImmutableGetGasLimitsResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewClientContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewGetGasLimits);
    }
}

export class GetMaintenanceStatusCall {
    func:    wasmlib.ScView;
    results: sc.ImmutableGetMaintenanceStatusResults = new sc.ImmutableGetMaintenanceStatusResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewClientContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewGetMaintenanceStatus);
    }
}

export class GetMetadataCall {
    func:    wasmlib.ScView;
    results: sc.ImmutableGetMetadataResults = new sc.ImmutableGetMetadataResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewClientContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewGetMetadata);
    }
}

export class ScFuncs {
    // Adds the given address to the list of identities that constitute the state controller.
    static addAllowedStateControllerAddress(ctx: wasmlib.ScFuncClientContext): AddAllowedStateControllerAddressCall {
        const f = new AddAllowedStateControllerAddressCall(ctx);
        f.params = new sc.MutableAddAllowedStateControllerAddressParams(wasmlib.newCallParamsProxy(f.func));
        return f;
    }

    // Adds a node to the list of candidates.
    static addCandidateNode(ctx: wasmlib.ScFuncClientContext): AddCandidateNodeCall {
        const f = new AddCandidateNodeCall(ctx);
        f.params = new sc.MutableAddCandidateNodeParams(wasmlib.newCallParamsProxy(f.func));
        return f;
    }

    // Iterates through the given map of actions and applies them.
    static changeAccessNodes(ctx: wasmlib.ScFuncClientContext): ChangeAccessNodesCall {
        const f = new ChangeAccessNodesCall(ctx);
        f.params = new sc.MutableChangeAccessNodesParams(wasmlib.newCallParamsProxy(f.func));
        return f;
    }

    // Claims the ownership of the chain if the caller matches the identity
    // that was set in delegateChainOwnership().
    static claimChainOwnership(ctx: wasmlib.ScFuncClientContext): ClaimChainOwnershipCall {
        return new ClaimChainOwnershipCall(ctx);
    }

    // Sets the Agent ID o as the new owner for the chain.
    // This change will only be effective once claimChainOwnership() is called by o.
    static delegateChainOwnership(ctx: wasmlib.ScFuncClientContext): DelegateChainOwnershipCall {
        const f = new DelegateChainOwnershipCall(ctx);
        f.params = new sc.MutableDelegateChainOwnershipParams(wasmlib.newCallParamsProxy(f.func));
        return f;
    }

    // Removes the given address from the list of identities that constitute the state controller.
    static removeAllowedStateControllerAddress(ctx: wasmlib.ScFuncClientContext): RemoveAllowedStateControllerAddressCall {
        const f = new RemoveAllowedStateControllerAddressCall(ctx);
        f.params = new sc.MutableRemoveAllowedStateControllerAddressParams(wasmlib.newCallParamsProxy(f.func));
        return f;
    }

    // Removes a node from the list of candidates.
    static revokeAccessNode(ctx: wasmlib.ScFuncClientContext): RevokeAccessNodeCall {
        const f = new RevokeAccessNodeCall(ctx);
        f.params = new sc.MutableRevokeAccessNodeParams(wasmlib.newCallParamsProxy(f.func));
        return f;
    }

    // Called when the committee is about to be rotated to the given address.
    // If it succeeds, the next state transition will become a governance transition,
    // thus updating the state controller in the chain's Alias Output.
    // If it fails, nothing happens.
    static rotateStateController(ctx: wasmlib.ScFuncClientContext): RotateStateControllerCall {
        const f = new RotateStateControllerCall(ctx);
        f.params = new sc.MutableRotateStateControllerParams(wasmlib.newCallParamsProxy(f.func));
        return f;
    }

    // Sets the EVM gas ratio for the chain.
    static setEVMGasRatio(ctx: wasmlib.ScFuncClientContext): SetEVMGasRatioCall {
        const f = new SetEVMGasRatioCall(ctx);
        f.params = new sc.MutableSetEVMGasRatioParams(wasmlib.newCallParamsProxy(f.func));
        return f;
    }

    // Sets the fee policy for the chain.
    static setFeePolicy(ctx: wasmlib.ScFuncClientContext): SetFeePolicyCall {
        const f = new SetFeePolicyCall(ctx);
        f.params = new sc.MutableSetFeePolicyParams(wasmlib.newCallParamsProxy(f.func));
        return f;
    }

    // Sets the gas limits for the chain.
    static setGasLimits(ctx: wasmlib.ScFuncClientContext): SetGasLimitsCall {
        const f = new SetGasLimitsCall(ctx);
        f.params = new sc.MutableSetGasLimitsParams(wasmlib.newCallParamsProxy(f.func));
        return f;
    }

    // Changes optional extra metadata that is appended to the L1 AliasOutput.
    static setMetadata(ctx: wasmlib.ScFuncClientContext): SetMetadataCall {
        const f = new SetMetadataCall(ctx);
        f.params = new sc.MutableSetMetadataParams(wasmlib.newCallParamsProxy(f.func));
        return f;
    }

    // Starts the chain maintenance mode, meaning no further requests
    // will be processed except calls to the governance contract.
    static startMaintenance(ctx: wasmlib.ScFuncClientContext): StartMaintenanceCall {
        return new StartMaintenanceCall(ctx);
    }

    // Stops the maintenance mode.
    static stopMaintenance(ctx: wasmlib.ScFuncClientContext): StopMaintenanceCall {
        return new StopMaintenanceCall(ctx);
    }

    // Returns the list of allowed state controllers.
    static getAllowedStateControllerAddresses(ctx: wasmlib.ScViewClientContext): GetAllowedStateControllerAddressesCall {
        const f = new GetAllowedStateControllerAddressesCall(ctx);
        f.results = new sc.ImmutableGetAllowedStateControllerAddressesResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns information about the chain.
    static getChainInfo(ctx: wasmlib.ScViewClientContext): GetChainInfoCall {
        const f = new GetChainInfoCall(ctx);
        f.results = new sc.ImmutableGetChainInfoResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns the current access nodes and candidates.
    static getChainNodes(ctx: wasmlib.ScViewClientContext): GetChainNodesCall {
        const f = new GetChainNodesCall(ctx);
        f.results = new sc.ImmutableGetChainNodesResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns the AgentID of the chain owner.
    static getChainOwner(ctx: wasmlib.ScViewClientContext): GetChainOwnerCall {
        const f = new GetChainOwnerCall(ctx);
        f.results = new sc.ImmutableGetChainOwnerResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns the EVM gas ratio.
    static getEVMGasRatio(ctx: wasmlib.ScViewClientContext): GetEVMGasRatioCall {
        const f = new GetEVMGasRatioCall(ctx);
        f.results = new sc.ImmutableGetEVMGasRatioResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns the fee policy.
    static getFeePolicy(ctx: wasmlib.ScViewClientContext): GetFeePolicyCall {
        const f = new GetFeePolicyCall(ctx);
        f.results = new sc.ImmutableGetFeePolicyResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns the gas limits.
    static getGasLimits(ctx: wasmlib.ScViewClientContext): GetGasLimitsCall {
        const f = new GetGasLimitsCall(ctx);
        f.results = new sc.ImmutableGetGasLimitsResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns whether the chain is undergoing maintenance.
    static getMaintenanceStatus(ctx: wasmlib.ScViewClientContext): GetMaintenanceStatusCall {
        const f = new GetMaintenanceStatusCall(ctx);
        f.results = new sc.ImmutableGetMaintenanceStatusResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns the extra metadata that is added to the chain AliasOutput.
    static getMetadata(ctx: wasmlib.ScViewClientContext): GetMetadataCall {
        const f = new GetMetadataCall(ctx);
        f.results = new sc.ImmutableGetMetadataResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }
}
