// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

import * as wasmlib from "../wasmlib"
import * as coreblocklog from "../coreblocklog"
import * as sc from "./index";

export function funcArrayClear(ctx: wasmlib.ScFuncContext, f: sc.ArrayClearContext): void {
    let name = f.params.name().value();
    let array = f.state.arrays().getStringArray(name);
    array.clear();
}

export function funcArrayCreate(ctx: wasmlib.ScFuncContext, f: sc.ArrayCreateContext): void {
    let name = f.params.name().value();
    let array = f.state.arrays().getStringArray(name);
    array.clear();
}

export function funcArraySet(ctx: wasmlib.ScFuncContext, f: sc.ArraySetContext): void {
    let name = f.params.name().value();
    let array = f.state.arrays().getStringArray(name);
    let index = f.params.index().value();
    let value = f.params.value().value();
    array.getString(index).setValue(value);
}

export function funcParamTypes(ctx: wasmlib.ScFuncContext, f: sc.ParamTypesContext): void {
    if (f.params.address().exists()) {
        ctx.require(f.params.address().value().equals(ctx.accountID().address()), "mismatch: Address");
    }
    if (f.params.agentID().exists()) {
        ctx.require(f.params.agentID().value().equals(ctx.accountID()), "mismatch: AgentID");
    }
    if (f.params.bytes().exists()) {
        let byteData = wasmlib.Convert.fromString("these are bytes");
        ctx.require(wasmlib.Convert.equals(f.params.bytes().value(), byteData), "mismatch: Bytes");
    }
    if (f.params.chainID().exists()) {
        ctx.require(f.params.chainID().value().equals(ctx.chainID()), "mismatch: ChainID");
    }
    if (f.params.color().exists()) {
        let color = wasmlib.ScColor.fromBytes(wasmlib.Convert.fromString("RedGreenBlueYellowCyanBlackWhite"));
        ctx.require(f.params.color().value().equals(color), "mismatch: Color");
    }
    if (f.params.hash().exists()) {
        let hash = wasmlib.ScHash.fromBytes(wasmlib.Convert.fromString("0123456789abcdeffedcba9876543210"));
        ctx.require(f.params.hash().value().equals(hash), "mismatch: Hash");
    }
    if (f.params.hname().exists()) {
        ctx.require(f.params.hname().value().equals(ctx.accountID().hname()), "mismatch: Hname");
    }
    if (f.params.int16().exists()) {
        ctx.require(f.params.int16().value() == 12345, "mismatch: Int16");
    }
    if (f.params.int32().exists()) {
        ctx.require(f.params.int32().value() == 1234567890, "mismatch: Int32");
    }
    if (f.params.int64().exists()) {
        ctx.require(f.params.int64().value() == 1234567890123456789, "mismatch: Int64");
    }
    if (f.params.requestID().exists()) {
        let requestId = wasmlib.ScRequestID.fromBytes(wasmlib.Convert.fromString("abcdefghijklmnopqrstuvwxyz123456\x00\x00"));
        ctx.require(f.params.requestID().value().equals(requestId), "mismatch: RequestID");
    }
    if (f.params.string().exists()) {
        ctx.require(f.params.string().value() == "this is a string", "mismatch: String");
    }
}

export function viewArrayLength(ctx: wasmlib.ScViewContext, f: sc.ArrayLengthContext): void {
    let name = f.params.name().value();
    let array = f.state.arrays().getStringArray(name);
    let length = array.length();
    f.results.length().setValue(length);
}

export function viewArrayValue(ctx: wasmlib.ScViewContext, f: sc.ArrayValueContext): void {
    let name = f.params.name().value();
    let array = f.state.arrays().getStringArray(name);
    let index = f.params.index().value();
    let value = array.getString(index).value();
    f.results.value().setValue(value);
}

export function viewBlockRecord(ctx: wasmlib.ScViewContext, f: sc.BlockRecordContext): void {
    let records = coreblocklog.ScFuncs.getRequestReceiptsForBlock(ctx);
    records.params.blockIndex().setValue(f.params.blockIndex().value());
    records.func.call();
    let recordIndex = f.params.recordIndex().value();
    ctx.require(recordIndex < records.results.requestRecord().length(), "invalid recordIndex");
    f.results.record().setValue(records.results.requestRecord().getBytes(recordIndex).value());
}

export function viewBlockRecords(ctx: wasmlib.ScViewContext, f: sc.BlockRecordsContext): void {
    let records = coreblocklog.ScFuncs.getRequestReceiptsForBlock(ctx);
    records.params.blockIndex().setValue(f.params.blockIndex().value());
    records.func.call();
    f.results.count().setValue(records.results.requestRecord().length());
}

export function viewIotaBalance(ctx: wasmlib.ScViewContext, f: sc.IotaBalanceContext): void {
    f.results.iotas().setValue(ctx.balances().balance(wasmlib.ScColor.IOTA));
}
