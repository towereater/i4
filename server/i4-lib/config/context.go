package config

type ContextKey string

// Context
const ContextConfig ContextKey = "config"
const ContextHash ContextKey = "hash"
const ContextClientCode ContextKey = "client"
const ContextMachineCode ContextKey = "machine"

// Constants
const ClientAdminCode string = "00000"
