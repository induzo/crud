package crud

// ListModifiers will modify the query, it can be url.Values for example
// Mostly used for Get and GetList
type ListModifiers map[string][]string
