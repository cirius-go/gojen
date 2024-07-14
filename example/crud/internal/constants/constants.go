package constants

// ContextHookKey is the key used to set and get values from the context hook.
// ENUM(domain)
//
//go:generate go-enum -f=$GOFILE --marshal --names --values
type ContextHookKey string
