// COPYRIGHT OF A TEST SCHEMA DEFINITION 1
// COPYRIGHT OF A TEST SCHEMA DEFINITION 2

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

import * as wasmlib from "wasmlib";
import * as wasmtypes from "wasmlib/wasmtypes";

const schemaCommentHandlers = new Map<string, (evt: SchemaCommentEventHandlers, msg: string[]) => void>([
	["schemacomment.testEvent", (evt: SchemaCommentEventHandlers, msg: string[]) => evt.testEvent(new EventTestEvent(msg))],
	["schemacomment.testEventNoParams", (evt: SchemaCommentEventHandlers, msg: string[]) => evt.testEventNoParams(new EventTestEventNoParams(msg))],
]);

export class SchemaCommentEventHandlers implements wasmlib.IEventHandler {
/* eslint-disable @typescript-eslint/no-empty-function */
	testEvent: (evt: EventTestEvent) => void = () => {};
	testEventNoParams: (evt: EventTestEventNoParams) => void = () => {};
/* eslint-enable @typescript-eslint/no-empty-function */

	public callHandler(topic: string, params: string[]): void {
		const handler = schemaCommentHandlers.get(topic);
		if (handler) {
			handler(this, params);
		}
	}

	public onSchemaCommentTestEvent(handler: (evt: EventTestEvent) => void): void {
		this.testEvent = handler;
	}

	public onSchemaCommentTestEventNoParams(handler: (evt: EventTestEventNoParams) => void): void {
		this.testEventNoParams = handler;
	}
}

export class EventTestEvent {
	public readonly timestamp: u32;
	public readonly eventParam1: string;
	public readonly eventParam2: string;
	
	public constructor(msg: string[]) {
		const evt = new wasmlib.EventDecoder(msg);
		this.timestamp = evt.timestamp();
		this.eventParam1 = wasmtypes.stringFromString(evt.decode());
		this.eventParam2 = wasmtypes.stringFromString(evt.decode());
	}
}

export class EventTestEventNoParams {
	public readonly timestamp: u32;
	
	public constructor(msg: string[]) {
		const evt = new wasmlib.EventDecoder(msg);
		this.timestamp = evt.timestamp();
	}
}