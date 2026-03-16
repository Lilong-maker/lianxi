package request

type Goods struct {
	Name  string  `form:"user"   binding:"required"`
	Price float64 `form:"price"  binding:"required"`
	Num   int64   `form:"num"  binding:"required"`
}
