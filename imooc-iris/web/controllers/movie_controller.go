package controllers

import (
	"github.com/kataras/iris/mvc"
	"web/GoWeb商品秒杀系统/imooc-iris/repositories"
	"web/GoWeb商品秒杀系统/imooc-iris/serives"
)

type MovieController struct {}

func (m *MovieController) Get() mvc.View {
	movieRepository := repositories.NewMovieManager()
	movieService := serives.NewMovieSeriviceManager(movieRepository)
	movieResult := movieService.ShowMovieName()
	return mvc.View{
		Name:"movie/index.html",
		Data:movieResult,
	}
}