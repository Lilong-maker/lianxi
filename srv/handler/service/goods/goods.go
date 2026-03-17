package goods

import (
	"context"
	__ "lianxi/srv/proto/goods"

	"lianxi/srv/dasic/config"
	"lianxi/srv/handler/model"
)

type Server struct {
	__.UnimplementedGoodsServer
}

func (s *Server) GoodsAdd(_ context.Context, in *__.GoodsAddReq) (*__.GoodsAddResp, error) {

	var goods model.Goods
	err := goods.FindGoods(config.DB, in.Name)
	if err == nil {
		return &__.GoodsAddResp{
			Msg:  "商品名称已存在",
			Code: 400,
		}, nil
	}

	m := model.Goods{
		Name:  in.Name,
		Price: float64(in.Price),
		Num:   int(in.Num),
	}
	err = m.GoodsAdd(config.DB)
	if err != nil {
		return &__.GoodsAddResp{
			Msg:  "商品添加失败: " + err.Error(),
			Code: 500,
		}, nil
	}

	return &__.GoodsAddResp{
		Msg:  "商品添加成功",
		Code: 200,
	}, nil
}

func (s *Server) GoodsList(_ context.Context, in *__.GoodsListReq) (*__.GoodsListResp, error) {

	page := int(in.Page)
	pageSize := int(in.PageSize)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	goods, total, err := model.GoodsList(config.DB, page, pageSize)
	if err != nil {
		return &__.GoodsListResp{
			Msg:  "查询失败: " + err.Error(),
			Code: 500,
		}, nil
	}

	var data []*__.GoodsInfo
	for _, g := range goods {
		data = append(data, &__.GoodsInfo{
			Id:        uint64(g.ID),
			Name:      g.Name,
			Price:     float32(g.Price),
			Num:       int64(g.Num),
			CreatedAt: g.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &__.GoodsListResp{
		Msg:   "查询成功",
		Code:  200,
		Data:  data,
		Total: total,
	}, nil
}
func (s *Server) GoodsUpdate(_ context.Context, in *__.GoodsUpdateReq) (*__.GoodsUpdateResp, error) {
	var goods model.Goods
	err := goods.GetGoodsByID(config.DB, uint(in.Id))
	if err != nil {
		return &__.GoodsUpdateResp{
			Msg:  "商品不存在",
			Code: 404,
		}, nil
	}
	goods.Name = in.Name
	goods.Price = float64(in.Price)
	goods.Num = int(in.Num)

	err = goods.GoodsUpdate(config.DB)
	if err != nil {
		return &__.GoodsUpdateResp{
			Msg:  "更新失败: " + err.Error(),
			Code: 500,
		}, nil
	}

	return &__.GoodsUpdateResp{
		Msg:  "更新成功",
		Code: 200,
	}, nil
}

func (s *Server) GoodsDelete(_ context.Context, in *__.GoodsDeleteReq) (*__.GoodsDeleteResp, error) {

	var goods model.Goods
	err := goods.GetGoodsByID(config.DB, uint(in.Id))
	if err != nil {
		return &__.GoodsDeleteResp{
			Msg:  "商品不存在",
			Code: 404,
		}, nil
	}
	err = goods.GoodsDelete(config.DB)
	if err != nil {
		return &__.GoodsDeleteResp{
			Msg:  "删除失败: " + err.Error(),
			Code: 500,
		}, nil
	}
	return &__.GoodsDeleteResp{
		Msg:  "删除成功",
		Code: 200,
	}, nil
}
