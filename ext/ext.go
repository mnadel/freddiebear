package ext

import "strings"

const (
	// X_FREDDIEBEAR is the freddbear extension prefix.
	X_FREDDIEBEAR = `x-fb`
)

// KeyValue returns an extension-specific key and value.
func KeyValue(key, value string) string {
	b := strings.Builder{}

	WriteKeyValue(b, key, value)

	return b.String()
}

// KeyValue write an extension-specific key and value to the given Builder.
func WriteKeyValue(b strings.Builder, key string, value string) {
	b.WriteString(X_FREDDIEBEAR)
	b.WriteString(key)
	b.WriteString(":")
	b.WriteString(value)
}
