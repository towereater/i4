package config

type ContextKey string

// Context
const ContextConfig ContextKey = "config"
const ContextHash ContextKey = "hash"
const ContextClientCode ContextKey = "client"
const ContextMachineCode ContextKey = "machine"

const ContextDataType ContextKey = "type"
const ContextOperation ContextKey = "op"
const ContextTimestampFrom ContextKey = "tsfrom"
const ContextTimestampTo ContextKey = "tsto"
const ContextDataKey ContextKey = "key"
const ContextDataValue ContextKey = "value"

const ContextFrom ContextKey = "from"
const ContextLimit ContextKey = "limit"

// Constants
const ClientAdminCode string = "00000"
