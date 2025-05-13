package crop

import (
	"github.com/nogavadu/articles-service/internal/service"
)

type Implementation struct {
	cropServ service.CropService
}

func New(cropService service.CropService) *Implementation {
	return &Implementation{
		cropServ: cropService,
	}
}
