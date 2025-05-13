package category

import (
	"github.com/nogavadu/articles-service/internal/service"
)

type Implementation struct {
	categoryServ service.CategoryService
}

func New(categoryService service.CategoryService) *Implementation {
	return &Implementation{
		categoryServ: categoryService,
	}
}
