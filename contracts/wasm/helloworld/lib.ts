// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

import * as wasmlib from "../wasmlib"
import * as sc from "./index";

export function on_call(index: i32): void {
    return wasmlib.onCall(index);
}

export function on_load(): void {
    let exports = new wasmlib.ScExports();
    exports.addFunc(sc.FuncHelloWorld, funcHelloWorldThunk);
    exports.addView(sc.ViewGetHelloWorld, viewGetHelloWorldThunk);

    for (let i = 0; i < sc.keyMap.length; i++) {
        sc.idxMap[i] = wasmlib.Key32.fromString(sc.keyMap[i]);
    }
}

function funcHelloWorldThunk(ctx: wasmlib.ScFuncContext): void {
    ctx.log("helloworld.funcHelloWorld");
    let f = new sc.HelloWorldContext();
    f.state.mapID = wasmlib.OBJ_ID_STATE;
    sc.funcHelloWorld(ctx, f);
    ctx.log("helloworld.funcHelloWorld ok");
}

function viewGetHelloWorldThunk(ctx: wasmlib.ScViewContext): void {
    ctx.log("helloworld.viewGetHelloWorld");
    let f = new sc.GetHelloWorldContext();
    f.results.mapID = wasmlib.OBJ_ID_RESULTS;
    f.state.mapID = wasmlib.OBJ_ID_STATE;
    sc.viewGetHelloWorld(ctx, f);
    ctx.log("helloworld.viewGetHelloWorld ok");
}
