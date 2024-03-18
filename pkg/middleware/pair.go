package middleware

type Pair struct {
	path   string
	method string
}

func NewPair(path, method string) Pair {
	return Pair{
		path:   path,
		method: method,
	}
}
