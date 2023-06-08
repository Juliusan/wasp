// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package coreaccounts

import (
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"
)

var coreAccountsHandlers = map[string]func(*CoreAccountsEventHandlers, *wasmtypes.WasmDecoder){
	"coreaccounts.foundryCreated": func(evt *CoreAccountsEventHandlers, dec *wasmtypes.WasmDecoder) { evt.onCoreAccountsFoundryCreatedThunk(dec) },
	"coreaccounts.foundryDestroyed": func(evt *CoreAccountsEventHandlers, dec *wasmtypes.WasmDecoder) { evt.onCoreAccountsFoundryDestroyedThunk(dec) },
	"coreaccounts.foundryModified": func(evt *CoreAccountsEventHandlers, dec *wasmtypes.WasmDecoder) { evt.onCoreAccountsFoundryModifiedThunk(dec) },
}

type CoreAccountsEventHandlers struct {
	myID uint32
	foundryCreated func(e *EventFoundryCreated)
	foundryDestroyed func(e *EventFoundryDestroyed)
	foundryModified func(e *EventFoundryModified)
}

var _ wasmlib.IEventHandlers = new(CoreAccountsEventHandlers)

func NewCoreAccountsEventHandlers() *CoreAccountsEventHandlers {
	return &CoreAccountsEventHandlers{ myID: wasmlib.EventHandlersGenerateID() }
}

func (h *CoreAccountsEventHandlers) CallHandler(topic string, dec *wasmtypes.WasmDecoder) {
	handler := coreAccountsHandlers[topic]
	if handler != nil {
		handler(h, dec)
	}
}

func (h *CoreAccountsEventHandlers) ID() uint32 {
	return h.myID
}

func (h *CoreAccountsEventHandlers) OnCoreAccountsFoundryCreated(handler func(e *EventFoundryCreated)) {
	h.foundryCreated = handler
}

func (h *CoreAccountsEventHandlers) OnCoreAccountsFoundryDestroyed(handler func(e *EventFoundryDestroyed)) {
	h.foundryDestroyed = handler
}

func (h *CoreAccountsEventHandlers) OnCoreAccountsFoundryModified(handler func(e *EventFoundryModified)) {
	h.foundryModified = handler
}

type EventFoundryCreated struct {
	Timestamp uint64
	FoundrySN uint32
}

func (h *CoreAccountsEventHandlers) onCoreAccountsFoundryCreatedThunk(dec *wasmtypes.WasmDecoder) {
	if h.foundryCreated == nil {
		return
	}
	e := &EventFoundryCreated{}
	e.Timestamp = wasmtypes.Uint64Decode(dec)
	e.FoundrySN = wasmtypes.Uint32Decode(dec)
	dec.Close()
	h.foundryCreated(e)
}

type EventFoundryDestroyed struct {
	Timestamp uint64
	FoundrySN uint32
}

func (h *CoreAccountsEventHandlers) onCoreAccountsFoundryDestroyedThunk(dec *wasmtypes.WasmDecoder) {
	if h.foundryDestroyed == nil {
		return
	}
	e := &EventFoundryDestroyed{}
	e.Timestamp = wasmtypes.Uint64Decode(dec)
	e.FoundrySN = wasmtypes.Uint32Decode(dec)
	dec.Close()
	h.foundryDestroyed(e)
}

type EventFoundryModified struct {
	Timestamp uint64
	FoundrySN uint32
}

func (h *CoreAccountsEventHandlers) onCoreAccountsFoundryModifiedThunk(dec *wasmtypes.WasmDecoder) {
	if h.foundryModified == nil {
		return
	}
	e := &EventFoundryModified{}
	e.Timestamp = wasmtypes.Uint64Decode(dec)
	e.FoundrySN = wasmtypes.Uint32Decode(dec)
	dec.Close()
	h.foundryModified(e)
}