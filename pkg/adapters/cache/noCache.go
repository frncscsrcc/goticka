package cache

type NoCache struct{}

var noCache NoCache

func GetNoCache() NoCache {
	return noCache
}

func (c NoCache) Set(item Item) error {
	return nil
}

func (c NoCache) Get(item Item) Item {
	return Item{}
}
