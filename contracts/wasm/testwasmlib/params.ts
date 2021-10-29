// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

import * as wasmlib from "../wasmlib"
import * as sc from "./index";

export class ImmutableArrayClearParams extends wasmlib.ScMapID {

    name(): wasmlib.ScImmutableString {
        return new wasmlib.ScImmutableString(this.mapID, sc.idxMap[sc.IdxParamName]);
    }
}

export class MutableArrayClearParams extends wasmlib.ScMapID {

    name(): wasmlib.ScMutableString {
        return new wasmlib.ScMutableString(this.mapID, sc.idxMap[sc.IdxParamName]);
    }
}

export class ImmutableArrayCreateParams extends wasmlib.ScMapID {

    name(): wasmlib.ScImmutableString {
        return new wasmlib.ScImmutableString(this.mapID, sc.idxMap[sc.IdxParamName]);
    }
}

export class MutableArrayCreateParams extends wasmlib.ScMapID {

    name(): wasmlib.ScMutableString {
        return new wasmlib.ScMutableString(this.mapID, sc.idxMap[sc.IdxParamName]);
    }
}

export class ImmutableArraySetParams extends wasmlib.ScMapID {

    index(): wasmlib.ScImmutableInt32 {
        return new wasmlib.ScImmutableInt32(this.mapID, sc.idxMap[sc.IdxParamIndex]);
    }

    name(): wasmlib.ScImmutableString {
        return new wasmlib.ScImmutableString(this.mapID, sc.idxMap[sc.IdxParamName]);
    }

    value(): wasmlib.ScImmutableString {
        return new wasmlib.ScImmutableString(this.mapID, sc.idxMap[sc.IdxParamValue]);
    }
}

export class MutableArraySetParams extends wasmlib.ScMapID {

    index(): wasmlib.ScMutableInt32 {
        return new wasmlib.ScMutableInt32(this.mapID, sc.idxMap[sc.IdxParamIndex]);
    }

    name(): wasmlib.ScMutableString {
        return new wasmlib.ScMutableString(this.mapID, sc.idxMap[sc.IdxParamName]);
    }

    value(): wasmlib.ScMutableString {
        return new wasmlib.ScMutableString(this.mapID, sc.idxMap[sc.IdxParamValue]);
    }
}

export class MapStringToImmutableBytes {
    objID: i32;

    constructor(objID: i32) {
        this.objID = objID;
    }

    getBytes(key: string): wasmlib.ScImmutableBytes {
        return new wasmlib.ScImmutableBytes(this.objID, wasmlib.Key32.fromString(key).getKeyID());
    }
}

export class ImmutableParamTypesParams extends wasmlib.ScMapID {

    address(): wasmlib.ScImmutableAddress {
        return new wasmlib.ScImmutableAddress(this.mapID, sc.idxMap[sc.IdxParamAddress]);
    }

    agentID(): wasmlib.ScImmutableAgentID {
        return new wasmlib.ScImmutableAgentID(this.mapID, sc.idxMap[sc.IdxParamAgentID]);
    }

    bytes(): wasmlib.ScImmutableBytes {
        return new wasmlib.ScImmutableBytes(this.mapID, sc.idxMap[sc.IdxParamBytes]);
    }

    chainID(): wasmlib.ScImmutableChainID {
        return new wasmlib.ScImmutableChainID(this.mapID, sc.idxMap[sc.IdxParamChainID]);
    }

    color(): wasmlib.ScImmutableColor {
        return new wasmlib.ScImmutableColor(this.mapID, sc.idxMap[sc.IdxParamColor]);
    }

    hash(): wasmlib.ScImmutableHash {
        return new wasmlib.ScImmutableHash(this.mapID, sc.idxMap[sc.IdxParamHash]);
    }

    hname(): wasmlib.ScImmutableHname {
        return new wasmlib.ScImmutableHname(this.mapID, sc.idxMap[sc.IdxParamHname]);
    }

    int16(): wasmlib.ScImmutableInt16 {
        return new wasmlib.ScImmutableInt16(this.mapID, sc.idxMap[sc.IdxParamInt16]);
    }

    int32(): wasmlib.ScImmutableInt32 {
        return new wasmlib.ScImmutableInt32(this.mapID, sc.idxMap[sc.IdxParamInt32]);
    }

    int64(): wasmlib.ScImmutableInt64 {
        return new wasmlib.ScImmutableInt64(this.mapID, sc.idxMap[sc.IdxParamInt64]);
    }

    param(): sc.MapStringToImmutableBytes {
        return new sc.MapStringToImmutableBytes(this.mapID);
    }

    requestID(): wasmlib.ScImmutableRequestID {
        return new wasmlib.ScImmutableRequestID(this.mapID, sc.idxMap[sc.IdxParamRequestID]);
    }

    string(): wasmlib.ScImmutableString {
        return new wasmlib.ScImmutableString(this.mapID, sc.idxMap[sc.IdxParamString]);
    }
}

export class MapStringToMutableBytes {
    objID: i32;

    constructor(objID: i32) {
        this.objID = objID;
    }

    clear(): void {
        wasmlib.clear(this.objID)
    }

    getBytes(key: string): wasmlib.ScMutableBytes {
        return new wasmlib.ScMutableBytes(this.objID, wasmlib.Key32.fromString(key).getKeyID());
    }
}

export class MutableParamTypesParams extends wasmlib.ScMapID {

    address(): wasmlib.ScMutableAddress {
        return new wasmlib.ScMutableAddress(this.mapID, sc.idxMap[sc.IdxParamAddress]);
    }

    agentID(): wasmlib.ScMutableAgentID {
        return new wasmlib.ScMutableAgentID(this.mapID, sc.idxMap[sc.IdxParamAgentID]);
    }

    bytes(): wasmlib.ScMutableBytes {
        return new wasmlib.ScMutableBytes(this.mapID, sc.idxMap[sc.IdxParamBytes]);
    }

    chainID(): wasmlib.ScMutableChainID {
        return new wasmlib.ScMutableChainID(this.mapID, sc.idxMap[sc.IdxParamChainID]);
    }

    color(): wasmlib.ScMutableColor {
        return new wasmlib.ScMutableColor(this.mapID, sc.idxMap[sc.IdxParamColor]);
    }

    hash(): wasmlib.ScMutableHash {
        return new wasmlib.ScMutableHash(this.mapID, sc.idxMap[sc.IdxParamHash]);
    }

    hname(): wasmlib.ScMutableHname {
        return new wasmlib.ScMutableHname(this.mapID, sc.idxMap[sc.IdxParamHname]);
    }

    int16(): wasmlib.ScMutableInt16 {
        return new wasmlib.ScMutableInt16(this.mapID, sc.idxMap[sc.IdxParamInt16]);
    }

    int32(): wasmlib.ScMutableInt32 {
        return new wasmlib.ScMutableInt32(this.mapID, sc.idxMap[sc.IdxParamInt32]);
    }

    int64(): wasmlib.ScMutableInt64 {
        return new wasmlib.ScMutableInt64(this.mapID, sc.idxMap[sc.IdxParamInt64]);
    }

    param(): sc.MapStringToMutableBytes {
        return new sc.MapStringToMutableBytes(this.mapID);
    }

    requestID(): wasmlib.ScMutableRequestID {
        return new wasmlib.ScMutableRequestID(this.mapID, sc.idxMap[sc.IdxParamRequestID]);
    }

    string(): wasmlib.ScMutableString {
        return new wasmlib.ScMutableString(this.mapID, sc.idxMap[sc.IdxParamString]);
    }
}

export class ImmutableArrayLengthParams extends wasmlib.ScMapID {

    name(): wasmlib.ScImmutableString {
        return new wasmlib.ScImmutableString(this.mapID, sc.idxMap[sc.IdxParamName]);
    }
}

export class MutableArrayLengthParams extends wasmlib.ScMapID {

    name(): wasmlib.ScMutableString {
        return new wasmlib.ScMutableString(this.mapID, sc.idxMap[sc.IdxParamName]);
    }
}

export class ImmutableArrayValueParams extends wasmlib.ScMapID {

    index(): wasmlib.ScImmutableInt32 {
        return new wasmlib.ScImmutableInt32(this.mapID, sc.idxMap[sc.IdxParamIndex]);
    }

    name(): wasmlib.ScImmutableString {
        return new wasmlib.ScImmutableString(this.mapID, sc.idxMap[sc.IdxParamName]);
    }
}

export class MutableArrayValueParams extends wasmlib.ScMapID {

    index(): wasmlib.ScMutableInt32 {
        return new wasmlib.ScMutableInt32(this.mapID, sc.idxMap[sc.IdxParamIndex]);
    }

    name(): wasmlib.ScMutableString {
        return new wasmlib.ScMutableString(this.mapID, sc.idxMap[sc.IdxParamName]);
    }
}

export class ImmutableBlockRecordParams extends wasmlib.ScMapID {

    blockIndex(): wasmlib.ScImmutableInt32 {
        return new wasmlib.ScImmutableInt32(this.mapID, sc.idxMap[sc.IdxParamBlockIndex]);
    }

    recordIndex(): wasmlib.ScImmutableInt32 {
        return new wasmlib.ScImmutableInt32(this.mapID, sc.idxMap[sc.IdxParamRecordIndex]);
    }
}

export class MutableBlockRecordParams extends wasmlib.ScMapID {

    blockIndex(): wasmlib.ScMutableInt32 {
        return new wasmlib.ScMutableInt32(this.mapID, sc.idxMap[sc.IdxParamBlockIndex]);
    }

    recordIndex(): wasmlib.ScMutableInt32 {
        return new wasmlib.ScMutableInt32(this.mapID, sc.idxMap[sc.IdxParamRecordIndex]);
    }
}

export class ImmutableBlockRecordsParams extends wasmlib.ScMapID {

    blockIndex(): wasmlib.ScImmutableInt32 {
        return new wasmlib.ScImmutableInt32(this.mapID, sc.idxMap[sc.IdxParamBlockIndex]);
    }
}

export class MutableBlockRecordsParams extends wasmlib.ScMapID {

    blockIndex(): wasmlib.ScMutableInt32 {
        return new wasmlib.ScMutableInt32(this.mapID, sc.idxMap[sc.IdxParamBlockIndex]);
    }
}
