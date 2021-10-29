// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

import * as wasmlib from "../wasmlib"
import * as sc from "./index";

export class MapStringToImmutableStringArray {
    objID: i32;

    constructor(objID: i32) {
        this.objID = objID;
    }

    getStringArray(key: string): sc.ImmutableStringArray {
        let subID = wasmlib.getObjectID(this.objID, wasmlib.Key32.fromString(key).getKeyID(), wasmlib.TYPE_ARRAY|wasmlib.TYPE_STRING);
        return new sc.ImmutableStringArray(subID);
    }
}

export class ImmutableTestWasmLibState extends wasmlib.ScMapID {

    arrays(): sc.MapStringToImmutableStringArray {
        let mapID = wasmlib.getObjectID(this.mapID, sc.idxMap[sc.IdxStateArrays], wasmlib.TYPE_MAP);
        return new sc.MapStringToImmutableStringArray(mapID);
    }
}

export class MapStringToMutableStringArray {
    objID: i32;

    constructor(objID: i32) {
        this.objID = objID;
    }

    clear(): void {
        wasmlib.clear(this.objID)
    }

    getStringArray(key: string): sc.MutableStringArray {
        let subID = wasmlib.getObjectID(this.objID, wasmlib.Key32.fromString(key).getKeyID(), wasmlib.TYPE_ARRAY|wasmlib.TYPE_STRING);
        return new sc.MutableStringArray(subID);
    }
}

export class MutableTestWasmLibState extends wasmlib.ScMapID {

    arrays(): sc.MapStringToMutableStringArray {
        let mapID = wasmlib.getObjectID(this.mapID, sc.idxMap[sc.IdxStateArrays], wasmlib.TYPE_MAP);
        return new sc.MapStringToMutableStringArray(mapID);
    }
}
