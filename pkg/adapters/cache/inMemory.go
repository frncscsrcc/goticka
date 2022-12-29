package cache

type InMemory struct {
	data map[string]map[string]Item
}

var inMemoryCache *InMemory

func init() {
	inMemoryCache = &InMemory{
		data: make(map[string]map[string]Item),
	}
}

func GetInMemoryCache() *InMemory {
	return inMemoryCache
}

func (c InMemory) Set(item Item) error {
	if _, cacheTypeExists := c.data[item.Type]; !cacheTypeExists {
		c.data[item.Type] = make(map[string]Item)
	}
	item.valid = true
	c.data[item.Type][item.Key] = item
	return nil
}

func (c InMemory) Delete(item Item) {
	if _, cacheTypeExists := c.data[item.Type]; !cacheTypeExists {
		return
	}
	delete(c.data[item.Type], item.Key)
	return
}

func (c InMemory) Get(item Item) Item {
	if _, cacheTypeExists := c.data[item.Type]; !cacheTypeExists {
		return Item{}
	}
	if value, cacheKeyExists := c.data[item.Type][item.Key]; !cacheKeyExists {
		return Item{}
	} else {
		return value
	}
}
