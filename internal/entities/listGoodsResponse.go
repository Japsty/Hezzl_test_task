package entities

type meta struct {
	Total   int
	Removed int
	Limit   int
	Offset  int
}

type goods []Good

type GoodsList struct {
	meta
	goods
}
