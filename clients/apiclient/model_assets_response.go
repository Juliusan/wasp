/*
Wasp API

REST API for the Wasp node

API version: 0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package apiclient

import (
	"encoding/json"
)

// checks if the AssetsResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AssetsResponse{}

// AssetsResponse struct for AssetsResponse
type AssetsResponse struct {
	// The base tokens (uint64 as string)
	BaseTokens string `json:"baseTokens"`
	NativeTokens []NativeTokenJSON `json:"nativeTokens"`
}

// NewAssetsResponse instantiates a new AssetsResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAssetsResponse(baseTokens string, nativeTokens []NativeTokenJSON) *AssetsResponse {
	this := AssetsResponse{}
	this.BaseTokens = baseTokens
	this.NativeTokens = nativeTokens
	return &this
}

// NewAssetsResponseWithDefaults instantiates a new AssetsResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAssetsResponseWithDefaults() *AssetsResponse {
	this := AssetsResponse{}
	return &this
}

// GetBaseTokens returns the BaseTokens field value
func (o *AssetsResponse) GetBaseTokens() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.BaseTokens
}

// GetBaseTokensOk returns a tuple with the BaseTokens field value
// and a boolean to check if the value has been set.
func (o *AssetsResponse) GetBaseTokensOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.BaseTokens, true
}

// SetBaseTokens sets field value
func (o *AssetsResponse) SetBaseTokens(v string) {
	o.BaseTokens = v
}

// GetNativeTokens returns the NativeTokens field value
func (o *AssetsResponse) GetNativeTokens() []NativeTokenJSON {
	if o == nil {
		var ret []NativeTokenJSON
		return ret
	}

	return o.NativeTokens
}

// GetNativeTokensOk returns a tuple with the NativeTokens field value
// and a boolean to check if the value has been set.
func (o *AssetsResponse) GetNativeTokensOk() ([]NativeTokenJSON, bool) {
	if o == nil {
		return nil, false
	}
	return o.NativeTokens, true
}

// SetNativeTokens sets field value
func (o *AssetsResponse) SetNativeTokens(v []NativeTokenJSON) {
	o.NativeTokens = v
}

func (o AssetsResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AssetsResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["baseTokens"] = o.BaseTokens
	toSerialize["nativeTokens"] = o.NativeTokens
	return toSerialize, nil
}

type NullableAssetsResponse struct {
	value *AssetsResponse
	isSet bool
}

func (v NullableAssetsResponse) Get() *AssetsResponse {
	return v.value
}

func (v *NullableAssetsResponse) Set(val *AssetsResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableAssetsResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableAssetsResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAssetsResponse(val *AssetsResponse) *NullableAssetsResponse {
	return &NullableAssetsResponse{value: val, isSet: true}
}

func (v NullableAssetsResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAssetsResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


