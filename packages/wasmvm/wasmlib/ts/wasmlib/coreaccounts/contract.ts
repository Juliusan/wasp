// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the schema definition file instead

import * as wasmlib from '../index';
import * as sc from './index';

export class DepositCall {
    func: wasmlib.ScFunc;

    public constructor(ctx: wasmlib.ScFuncCallContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncDeposit);
    }
}

export class FoundryCreateNewCall {
    func:    wasmlib.ScFunc;
    params:  sc.MutableFoundryCreateNewParams = new sc.MutableFoundryCreateNewParams(wasmlib.ScView.nilProxy);
    results: sc.ImmutableFoundryCreateNewResults = new sc.ImmutableFoundryCreateNewResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScFuncCallContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncFoundryCreateNew);
    }
}

export class FoundryDestroyCall {
    func:   wasmlib.ScFunc;
    params: sc.MutableFoundryDestroyParams = new sc.MutableFoundryDestroyParams(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScFuncCallContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncFoundryDestroy);
    }
}

export class FoundryModifySupplyCall {
    func:   wasmlib.ScFunc;
    params: sc.MutableFoundryModifySupplyParams = new sc.MutableFoundryModifySupplyParams(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScFuncCallContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncFoundryModifySupply);
    }
}

export class HarvestCall {
    func:   wasmlib.ScFunc;
    params: sc.MutableHarvestParams = new sc.MutableHarvestParams(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScFuncCallContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncHarvest);
    }
}

export class TransferAccountToChainCall {
    func:   wasmlib.ScFunc;
    params: sc.MutableTransferAccountToChainParams = new sc.MutableTransferAccountToChainParams(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScFuncCallContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncTransferAccountToChain);
    }
}

export class TransferAllowanceToCall {
    func:   wasmlib.ScFunc;
    params: sc.MutableTransferAllowanceToParams = new sc.MutableTransferAllowanceToParams(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScFuncCallContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncTransferAllowanceTo);
    }
}

export class WithdrawCall {
    func: wasmlib.ScFunc;

    public constructor(ctx: wasmlib.ScFuncCallContext) {
        this.func = new wasmlib.ScFunc(ctx, sc.HScName, sc.HFuncWithdraw);
    }
}

export class AccountFoundriesCall {
    func:    wasmlib.ScView;
    params:  sc.MutableAccountFoundriesParams = new sc.MutableAccountFoundriesParams(wasmlib.ScView.nilProxy);
    results: sc.ImmutableAccountFoundriesResults = new sc.ImmutableAccountFoundriesResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewCallContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewAccountFoundries);
    }
}

export class AccountNFTAmountCall {
    func:    wasmlib.ScView;
    params:  sc.MutableAccountNFTAmountParams = new sc.MutableAccountNFTAmountParams(wasmlib.ScView.nilProxy);
    results: sc.ImmutableAccountNFTAmountResults = new sc.ImmutableAccountNFTAmountResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewCallContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewAccountNFTAmount);
    }
}

export class AccountNFTAmountInCollectionCall {
    func:    wasmlib.ScView;
    params:  sc.MutableAccountNFTAmountInCollectionParams = new sc.MutableAccountNFTAmountInCollectionParams(wasmlib.ScView.nilProxy);
    results: sc.ImmutableAccountNFTAmountInCollectionResults = new sc.ImmutableAccountNFTAmountInCollectionResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewCallContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewAccountNFTAmountInCollection);
    }
}

export class AccountNFTsCall {
    func:    wasmlib.ScView;
    params:  sc.MutableAccountNFTsParams = new sc.MutableAccountNFTsParams(wasmlib.ScView.nilProxy);
    results: sc.ImmutableAccountNFTsResults = new sc.ImmutableAccountNFTsResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewCallContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewAccountNFTs);
    }
}

export class AccountNFTsInCollectionCall {
    func:    wasmlib.ScView;
    params:  sc.MutableAccountNFTsInCollectionParams = new sc.MutableAccountNFTsInCollectionParams(wasmlib.ScView.nilProxy);
    results: sc.ImmutableAccountNFTsInCollectionResults = new sc.ImmutableAccountNFTsInCollectionResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewCallContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewAccountNFTsInCollection);
    }
}

export class AccountsCall {
    func:    wasmlib.ScView;
    results: sc.ImmutableAccountsResults = new sc.ImmutableAccountsResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewCallContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewAccounts);
    }
}

export class BalanceCall {
    func:    wasmlib.ScView;
    params:  sc.MutableBalanceParams = new sc.MutableBalanceParams(wasmlib.ScView.nilProxy);
    results: sc.ImmutableBalanceResults = new sc.ImmutableBalanceResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewCallContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewBalance);
    }
}

export class BalanceBaseTokenCall {
    func:    wasmlib.ScView;
    params:  sc.MutableBalanceBaseTokenParams = new sc.MutableBalanceBaseTokenParams(wasmlib.ScView.nilProxy);
    results: sc.ImmutableBalanceBaseTokenResults = new sc.ImmutableBalanceBaseTokenResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewCallContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewBalanceBaseToken);
    }
}

export class BalanceNativeTokenCall {
    func:    wasmlib.ScView;
    params:  sc.MutableBalanceNativeTokenParams = new sc.MutableBalanceNativeTokenParams(wasmlib.ScView.nilProxy);
    results: sc.ImmutableBalanceNativeTokenResults = new sc.ImmutableBalanceNativeTokenResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewCallContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewBalanceNativeToken);
    }
}

export class FoundryOutputCall {
    func:    wasmlib.ScView;
    params:  sc.MutableFoundryOutputParams = new sc.MutableFoundryOutputParams(wasmlib.ScView.nilProxy);
    results: sc.ImmutableFoundryOutputResults = new sc.ImmutableFoundryOutputResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewCallContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewFoundryOutput);
    }
}

export class GetAccountNonceCall {
    func:    wasmlib.ScView;
    params:  sc.MutableGetAccountNonceParams = new sc.MutableGetAccountNonceParams(wasmlib.ScView.nilProxy);
    results: sc.ImmutableGetAccountNonceResults = new sc.ImmutableGetAccountNonceResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewCallContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewGetAccountNonce);
    }
}

export class GetNativeTokenIDRegistryCall {
    func:    wasmlib.ScView;
    results: sc.ImmutableGetNativeTokenIDRegistryResults = new sc.ImmutableGetNativeTokenIDRegistryResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewCallContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewGetNativeTokenIDRegistry);
    }
}

export class NftDataCall {
    func:    wasmlib.ScView;
    params:  sc.MutableNftDataParams = new sc.MutableNftDataParams(wasmlib.ScView.nilProxy);
    results: sc.ImmutableNftDataResults = new sc.ImmutableNftDataResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewCallContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewNftData);
    }
}

export class TotalAssetsCall {
    func:    wasmlib.ScView;
    results: sc.ImmutableTotalAssetsResults = new sc.ImmutableTotalAssetsResults(wasmlib.ScView.nilProxy);

    public constructor(ctx: wasmlib.ScViewCallContext) {
        this.func = new wasmlib.ScView(ctx, sc.HScName, sc.HViewTotalAssets);
    }
}

export class ScFuncs {
    // A no-op that has the side effect of crediting any transferred tokens to the sender's account.
    static deposit(ctx: wasmlib.ScFuncCallContext): DepositCall {
        return new DepositCall(ctx);
    }

    // Creates a new foundry with the specified token scheme, and assigns the foundry to the sender.
    static foundryCreateNew(ctx: wasmlib.ScFuncCallContext): FoundryCreateNewCall {
        const f = new FoundryCreateNewCall(ctx);
        f.params = new sc.MutableFoundryCreateNewParams(wasmlib.newCallParamsProxy(f.func));
        f.results = new sc.ImmutableFoundryCreateNewResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Destroys a given foundry output on L1, reimbursing the storage deposit to the caller.
    // The foundry must be owned by the caller.
    static foundryDestroy(ctx: wasmlib.ScFuncCallContext): FoundryDestroyCall {
        const f = new FoundryDestroyCall(ctx);
        f.params = new sc.MutableFoundryDestroyParams(wasmlib.newCallParamsProxy(f.func));
        return f;
    }

    // Mints or destroys tokens for the given foundry, which must be owned by the caller.
    static foundryModifySupply(ctx: wasmlib.ScFuncCallContext): FoundryModifySupplyCall {
        const f = new FoundryModifySupplyCall(ctx);
        f.params = new sc.MutableFoundryModifySupplyParams(wasmlib.newCallParamsProxy(f.func));
        return f;
    }

    // Moves all tokens from the chain common account to the sender's L2 account.
    // The chain owner is the only one who can call this entry point.
    static harvest(ctx: wasmlib.ScFuncCallContext): HarvestCall {
        const f = new HarvestCall(ctx);
        f.params = new sc.MutableHarvestParams(wasmlib.newCallParamsProxy(f.func));
        return f;
    }

    // Moves the specified allowance from the sender SC's L2 account on the target chain
    // to sender SC's L2 account on the origin chain.
    static transferAccountToChain(ctx: wasmlib.ScFuncCallContext): TransferAccountToChainCall {
        const f = new TransferAccountToChainCall(ctx);
        f.params = new sc.MutableTransferAccountToChainParams(wasmlib.newCallParamsProxy(f.func));
        return f;
    }

    // Moves the specified allowance from the sender's L2 account to the given L2 account on the chain.
    static transferAllowanceTo(ctx: wasmlib.ScFuncCallContext): TransferAllowanceToCall {
        const f = new TransferAllowanceToCall(ctx);
        f.params = new sc.MutableTransferAllowanceToParams(wasmlib.newCallParamsProxy(f.func));
        return f;
    }

    // Moves tokens from the caller's on-chain account to the caller's L1 address.
    // The number of tokens to be withdrawn must be specified via the allowance of the request.
    static withdraw(ctx: wasmlib.ScFuncCallContext): WithdrawCall {
        return new WithdrawCall(ctx);
    }

    // Returns a set of all foundries owned by the given account.
    static accountFoundries(ctx: wasmlib.ScViewCallContext): AccountFoundriesCall {
        const f = new AccountFoundriesCall(ctx);
        f.params = new sc.MutableAccountFoundriesParams(wasmlib.newCallParamsProxy(f.func));
        f.results = new sc.ImmutableAccountFoundriesResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns the amount of NFTs owned by the given account.
    static accountNFTAmount(ctx: wasmlib.ScViewCallContext): AccountNFTAmountCall {
        const f = new AccountNFTAmountCall(ctx);
        f.params = new sc.MutableAccountNFTAmountParams(wasmlib.newCallParamsProxy(f.func));
        f.results = new sc.ImmutableAccountNFTAmountResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns the amount of NFTs in the specified collection owned by the given account.
    static accountNFTAmountInCollection(ctx: wasmlib.ScViewCallContext): AccountNFTAmountInCollectionCall {
        const f = new AccountNFTAmountInCollectionCall(ctx);
        f.params = new sc.MutableAccountNFTAmountInCollectionParams(wasmlib.newCallParamsProxy(f.func));
        f.results = new sc.ImmutableAccountNFTAmountInCollectionResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns the NFT IDs for all NFTs owned by the given account.
    static accountNFTs(ctx: wasmlib.ScViewCallContext): AccountNFTsCall {
        const f = new AccountNFTsCall(ctx);
        f.params = new sc.MutableAccountNFTsParams(wasmlib.newCallParamsProxy(f.func));
        f.results = new sc.ImmutableAccountNFTsResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns the NFT IDs for all NFTs in the specified collection owned by the given account.
    static accountNFTsInCollection(ctx: wasmlib.ScViewCallContext): AccountNFTsInCollectionCall {
        const f = new AccountNFTsInCollectionCall(ctx);
        f.params = new sc.MutableAccountNFTsInCollectionParams(wasmlib.newCallParamsProxy(f.func));
        f.results = new sc.ImmutableAccountNFTsInCollectionResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns a set of all agent IDs that own assets on the chain.
    static accounts(ctx: wasmlib.ScViewCallContext): AccountsCall {
        const f = new AccountsCall(ctx);
        f.results = new sc.ImmutableAccountsResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns the fungible tokens owned by the given Agent ID on the chain.
    static balance(ctx: wasmlib.ScViewCallContext): BalanceCall {
        const f = new BalanceCall(ctx);
        f.params = new sc.MutableBalanceParams(wasmlib.newCallParamsProxy(f.func));
        f.results = new sc.ImmutableBalanceResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns the amount of base tokens owned by an agent on the chain
    static balanceBaseToken(ctx: wasmlib.ScViewCallContext): BalanceBaseTokenCall {
        const f = new BalanceBaseTokenCall(ctx);
        f.params = new sc.MutableBalanceBaseTokenParams(wasmlib.newCallParamsProxy(f.func));
        f.results = new sc.ImmutableBalanceBaseTokenResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns the amount of specific native tokens owned by an agent on the chain
    static balanceNativeToken(ctx: wasmlib.ScViewCallContext): BalanceNativeTokenCall {
        const f = new BalanceNativeTokenCall(ctx);
        f.params = new sc.MutableBalanceNativeTokenParams(wasmlib.newCallParamsProxy(f.func));
        f.results = new sc.ImmutableBalanceNativeTokenResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns specified foundry output in serialized form.
    static foundryOutput(ctx: wasmlib.ScViewCallContext): FoundryOutputCall {
        const f = new FoundryOutputCall(ctx);
        f.params = new sc.MutableFoundryOutputParams(wasmlib.newCallParamsProxy(f.func));
        f.results = new sc.ImmutableFoundryOutputResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns the current account nonce for an Agent.
    // The account nonce is used to issue unique off-ledger requests.
    static getAccountNonce(ctx: wasmlib.ScViewCallContext): GetAccountNonceCall {
        const f = new GetAccountNonceCall(ctx);
        f.params = new sc.MutableGetAccountNonceParams(wasmlib.newCallParamsProxy(f.func));
        f.results = new sc.ImmutableGetAccountNonceResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns a set of all native tokenIDs that are owned by the chain.
    static getNativeTokenIDRegistry(ctx: wasmlib.ScViewCallContext): GetNativeTokenIDRegistryCall {
        const f = new GetNativeTokenIDRegistryCall(ctx);
        f.results = new sc.ImmutableGetNativeTokenIDRegistryResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns the data for a given NFT that is on the chain.
    static nftData(ctx: wasmlib.ScViewCallContext): NftDataCall {
        const f = new NftDataCall(ctx);
        f.params = new sc.MutableNftDataParams(wasmlib.newCallParamsProxy(f.func));
        f.results = new sc.ImmutableNftDataResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }

    // Returns the balances of all fungible tokens controlled by the chain.
    static totalAssets(ctx: wasmlib.ScViewCallContext): TotalAssetsCall {
        const f = new TotalAssetsCall(ctx);
        f.results = new sc.ImmutableTotalAssetsResults(wasmlib.newCallResultsProxy(f.func));
        return f;
    }
}
