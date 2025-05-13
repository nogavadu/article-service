package article

import (
	"github.com/nogavadu/articles-service/internal/service"
)

type Implementation struct {
	articleServ service.ArticleService
}

func New(articleService service.ArticleService) *Implementation {
	return &Implementation{
		articleServ: articleService,
	}
}
