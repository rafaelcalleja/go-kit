package ports

import (
	"net/http"

	"github.com/rafaelcalleja/go-kit/internal/common/server/httperr"
	"github.com/rafaelcalleja/go-kit/internal/store/application"
	"github.com/rafaelcalleja/go-kit/internal/store/application/command"
)

type HttpServer struct {
	app application.Application
}

func NewHttpServer(application application.Application) HttpServer {
	return HttpServer{
		app: application,
	}
}

func (g HttpServer) CreateProduct(response http.ResponseWriter, request *http.Request, productId string) {
	if err := g.app.CommandBus.Dispatch(request.Context(), command.NewCreateProductCommand(productId)); err != nil {
		httperr.InternalError(err.Error(), err, response, request)
		return
	}

	response.WriteHeader(http.StatusCreated)
}
