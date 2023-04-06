// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the schema definition file instead

import * as wasmtypes from '../wasmtypes';
import * as sc from './index';

export class ImmutableFoundryCreateNewParams extends wasmtypes.ScProxy {
    // token scheme for the new foundry
    tokenScheme(): wasmtypes.ScImmutableBytes {
        return new wasmtypes.ScImmutableBytes(this.proxy.root(sc.ParamTokenScheme));
    }
}

export class MutableFoundryCreateNewParams extends wasmtypes.ScProxy {
    // token scheme for the new foundry
    tokenScheme(): wasmtypes.ScMutableBytes {
        return new wasmtypes.ScMutableBytes(this.proxy.root(sc.ParamTokenScheme));
    }
}

export class ImmutableFoundryDestroyParams extends wasmtypes.ScProxy {
    // serial number of the foundry
    foundrySN(): wasmtypes.ScImmutableUint32 {
        return new wasmtypes.ScImmutableUint32(this.proxy.root(sc.ParamFoundrySN));
    }
}

export class MutableFoundryDestroyParams extends wasmtypes.ScProxy {
    // serial number of the foundry
    foundrySN(): wasmtypes.ScMutableUint32 {
        return new wasmtypes.ScMutableUint32(this.proxy.root(sc.ParamFoundrySN));
    }
}

export class ImmutableFoundryModifySupplyParams extends wasmtypes.ScProxy {
    // mint (default) or destroy tokens
    destroyTokens(): wasmtypes.ScImmutableBool {
        return new wasmtypes.ScImmutableBool(this.proxy.root(sc.ParamDestroyTokens));
    }

    // serial number of the foundry
    foundrySN(): wasmtypes.ScImmutableUint32 {
        return new wasmtypes.ScImmutableUint32(this.proxy.root(sc.ParamFoundrySN));
    }

    // positive nonzero amount to mint or destroy
    supplyDeltaAbs(): wasmtypes.ScImmutableBigInt {
        return new wasmtypes.ScImmutableBigInt(this.proxy.root(sc.ParamSupplyDeltaAbs));
    }
}

export class MutableFoundryModifySupplyParams extends wasmtypes.ScProxy {
    // mint (default) or destroy tokens
    destroyTokens(): wasmtypes.ScMutableBool {
        return new wasmtypes.ScMutableBool(this.proxy.root(sc.ParamDestroyTokens));
    }

    // serial number of the foundry
    foundrySN(): wasmtypes.ScMutableUint32 {
        return new wasmtypes.ScMutableUint32(this.proxy.root(sc.ParamFoundrySN));
    }

    // positive nonzero amount to mint or destroy
    supplyDeltaAbs(): wasmtypes.ScMutableBigInt {
        return new wasmtypes.ScMutableBigInt(this.proxy.root(sc.ParamSupplyDeltaAbs));
    }
}

export class ImmutableHarvestParams extends wasmtypes.ScProxy {
    // amount of base tokens to leave in the common account
    // default MinimumBaseTokensOnCommonAccount, can never be less
    forceMinimumBaseTokens(): wasmtypes.ScImmutableUint64 {
        return new wasmtypes.ScImmutableUint64(this.proxy.root(sc.ParamForceMinimumBaseTokens));
    }
}

export class MutableHarvestParams extends wasmtypes.ScProxy {
    // amount of base tokens to leave in the common account
    // default MinimumBaseTokensOnCommonAccount, can never be less
    forceMinimumBaseTokens(): wasmtypes.ScMutableUint64 {
        return new wasmtypes.ScMutableUint64(this.proxy.root(sc.ParamForceMinimumBaseTokens));
    }
}

export class ImmutableTransferAccountToChainParams extends wasmtypes.ScProxy {
    // Gas reserve in allowance for internal call to transferAllowanceTo
    // Default 100 (for now)
    gasReserve(): wasmtypes.ScImmutableUint64 {
        return new wasmtypes.ScImmutableUint64(this.proxy.root(sc.ParamGasReserve));
    }
}

export class MutableTransferAccountToChainParams extends wasmtypes.ScProxy {
    // Gas reserve in allowance for internal call to transferAllowanceTo
    // Default 100 (for now)
    gasReserve(): wasmtypes.ScMutableUint64 {
        return new wasmtypes.ScMutableUint64(this.proxy.root(sc.ParamGasReserve));
    }
}

export class ImmutableTransferAllowanceToParams extends wasmtypes.ScProxy {
    // The target L2 account
    agentID(): wasmtypes.ScImmutableAgentID {
        return new wasmtypes.ScImmutableAgentID(this.proxy.root(sc.ParamAgentID));
    }
}

export class MutableTransferAllowanceToParams extends wasmtypes.ScProxy {
    // The target L2 account
    agentID(): wasmtypes.ScMutableAgentID {
        return new wasmtypes.ScMutableAgentID(this.proxy.root(sc.ParamAgentID));
    }
}

export class ImmutableAccountFoundriesParams extends wasmtypes.ScProxy {
    // account agent ID
    agentID(): wasmtypes.ScImmutableAgentID {
        return new wasmtypes.ScImmutableAgentID(this.proxy.root(sc.ParamAgentID));
    }
}

export class MutableAccountFoundriesParams extends wasmtypes.ScProxy {
    // account agent ID
    agentID(): wasmtypes.ScMutableAgentID {
        return new wasmtypes.ScMutableAgentID(this.proxy.root(sc.ParamAgentID));
    }
}

export class ImmutableAccountNFTAmountParams extends wasmtypes.ScProxy {
    // account agent ID
    agentID(): wasmtypes.ScImmutableAgentID {
        return new wasmtypes.ScImmutableAgentID(this.proxy.root(sc.ParamAgentID));
    }
}

export class MutableAccountNFTAmountParams extends wasmtypes.ScProxy {
    // account agent ID
    agentID(): wasmtypes.ScMutableAgentID {
        return new wasmtypes.ScMutableAgentID(this.proxy.root(sc.ParamAgentID));
    }
}

export class ImmutableAccountNFTAmountInCollectionParams extends wasmtypes.ScProxy {
    // account agent ID
    agentID(): wasmtypes.ScImmutableAgentID {
        return new wasmtypes.ScImmutableAgentID(this.proxy.root(sc.ParamAgentID));
    }

    // NFT ID of collection
    collection(): wasmtypes.ScImmutableNftID {
        return new wasmtypes.ScImmutableNftID(this.proxy.root(sc.ParamCollection));
    }
}

export class MutableAccountNFTAmountInCollectionParams extends wasmtypes.ScProxy {
    // account agent ID
    agentID(): wasmtypes.ScMutableAgentID {
        return new wasmtypes.ScMutableAgentID(this.proxy.root(sc.ParamAgentID));
    }

    // NFT ID of collection
    collection(): wasmtypes.ScMutableNftID {
        return new wasmtypes.ScMutableNftID(this.proxy.root(sc.ParamCollection));
    }
}

export class ImmutableAccountNFTsParams extends wasmtypes.ScProxy {
    // account agent ID
    agentID(): wasmtypes.ScImmutableAgentID {
        return new wasmtypes.ScImmutableAgentID(this.proxy.root(sc.ParamAgentID));
    }
}

export class MutableAccountNFTsParams extends wasmtypes.ScProxy {
    // account agent ID
    agentID(): wasmtypes.ScMutableAgentID {
        return new wasmtypes.ScMutableAgentID(this.proxy.root(sc.ParamAgentID));
    }
}

export class ImmutableAccountNFTsInCollectionParams extends wasmtypes.ScProxy {
    // account agent ID
    agentID(): wasmtypes.ScImmutableAgentID {
        return new wasmtypes.ScImmutableAgentID(this.proxy.root(sc.ParamAgentID));
    }

    // NFT ID of collection
    collection(): wasmtypes.ScImmutableNftID {
        return new wasmtypes.ScImmutableNftID(this.proxy.root(sc.ParamCollection));
    }
}

export class MutableAccountNFTsInCollectionParams extends wasmtypes.ScProxy {
    // account agent ID
    agentID(): wasmtypes.ScMutableAgentID {
        return new wasmtypes.ScMutableAgentID(this.proxy.root(sc.ParamAgentID));
    }

    // NFT ID of collection
    collection(): wasmtypes.ScMutableNftID {
        return new wasmtypes.ScMutableNftID(this.proxy.root(sc.ParamCollection));
    }
}

export class ImmutableBalanceParams extends wasmtypes.ScProxy {
    // account agent ID
    agentID(): wasmtypes.ScImmutableAgentID {
        return new wasmtypes.ScImmutableAgentID(this.proxy.root(sc.ParamAgentID));
    }
}

export class MutableBalanceParams extends wasmtypes.ScProxy {
    // account agent ID
    agentID(): wasmtypes.ScMutableAgentID {
        return new wasmtypes.ScMutableAgentID(this.proxy.root(sc.ParamAgentID));
    }
}

export class ImmutableBalanceBaseTokenParams extends wasmtypes.ScProxy {
    // account agent ID
    agentID(): wasmtypes.ScImmutableAgentID {
        return new wasmtypes.ScImmutableAgentID(this.proxy.root(sc.ParamAgentID));
    }
}

export class MutableBalanceBaseTokenParams extends wasmtypes.ScProxy {
    // account agent ID
    agentID(): wasmtypes.ScMutableAgentID {
        return new wasmtypes.ScMutableAgentID(this.proxy.root(sc.ParamAgentID));
    }
}

export class ImmutableBalanceNativeTokenParams extends wasmtypes.ScProxy {
    // account agent ID
    agentID(): wasmtypes.ScImmutableAgentID {
        return new wasmtypes.ScImmutableAgentID(this.proxy.root(sc.ParamAgentID));
    }

    // native token ID
    tokenID(): wasmtypes.ScImmutableTokenID {
        return new wasmtypes.ScImmutableTokenID(this.proxy.root(sc.ParamTokenID));
    }
}

export class MutableBalanceNativeTokenParams extends wasmtypes.ScProxy {
    // account agent ID
    agentID(): wasmtypes.ScMutableAgentID {
        return new wasmtypes.ScMutableAgentID(this.proxy.root(sc.ParamAgentID));
    }

    // native token ID
    tokenID(): wasmtypes.ScMutableTokenID {
        return new wasmtypes.ScMutableTokenID(this.proxy.root(sc.ParamTokenID));
    }
}

export class ImmutableFoundryOutputParams extends wasmtypes.ScProxy {
    // serial number of the foundry
    foundrySN(): wasmtypes.ScImmutableUint32 {
        return new wasmtypes.ScImmutableUint32(this.proxy.root(sc.ParamFoundrySN));
    }
}

export class MutableFoundryOutputParams extends wasmtypes.ScProxy {
    // serial number of the foundry
    foundrySN(): wasmtypes.ScMutableUint32 {
        return new wasmtypes.ScMutableUint32(this.proxy.root(sc.ParamFoundrySN));
    }
}

export class ImmutableGetAccountNonceParams extends wasmtypes.ScProxy {
    // account agent ID
    agentID(): wasmtypes.ScImmutableAgentID {
        return new wasmtypes.ScImmutableAgentID(this.proxy.root(sc.ParamAgentID));
    }
}

export class MutableGetAccountNonceParams extends wasmtypes.ScProxy {
    // account agent ID
    agentID(): wasmtypes.ScMutableAgentID {
        return new wasmtypes.ScMutableAgentID(this.proxy.root(sc.ParamAgentID));
    }
}

export class ImmutableNftDataParams extends wasmtypes.ScProxy {
    // NFT ID
    nftID(): wasmtypes.ScImmutableNftID {
        return new wasmtypes.ScImmutableNftID(this.proxy.root(sc.ParamNftID));
    }
}

export class MutableNftDataParams extends wasmtypes.ScProxy {
    // NFT ID
    nftID(): wasmtypes.ScMutableNftID {
        return new wasmtypes.ScMutableNftID(this.proxy.root(sc.ParamNftID));
    }
}
