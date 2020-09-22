package interfaces

const (
	KeyError       = int32(-1)
	KeyLength      = KeyError - 1
	KeyLog         = KeyLength - 1
	KeyTrace       = KeyLog - 1
	KeyTraceHost   = KeyTrace - 1
	KeyUserDefined = KeyTraceHost - 1
)

type HostInterface interface {
	GetInt(objId int32, keyId int32) int64
	GetKeyId(key string) int32
	GetObjectId(objId int32, keyId int32, typeId int32) int32
	GetString(objId int32, keyId int32) string
	SetInt(objId int32, keyId int32, value int64)
	SetString(objId int32, keyId int32, value string)
}
