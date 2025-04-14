package repositories

import "web/GoWeb商品秒杀系统/imooc-iris/datamodels"

// 接口必须实现
type MovieRepository interface {
	GetMovieName() string
}

type MovieManager struct {}

func NewMovieManager() MovieRepository  {
	return &MovieManager{}
}

func (m *MovieManager) GetMovieName() string {
	// 模拟赋值给模型
	movie := &datamodels.Movie{"来慕课，看视频，学技术！"}
	return movie.Name
}

