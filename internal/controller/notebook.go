package controller

import (
	"context"
	"net/http"

	"github.com/imabg/sync/internal/models"
	"github.com/imabg/sync/pkg/config"
	"github.com/imabg/sync/pkg/errors"
	"github.com/imabg/sync/pkg/response"
	"github.com/imabg/sync/pkg/token"
	"github.com/imabg/sync/pkg/validate"
	"github.com/imabg/sync/services/notebook"
)

type INotebook interface {
	CreateNoteBook(http.ResponseWriter, *http.Request)
	UpdateNoteBook(http.ResponseWriter, *http.Request)
}

type NotebookCtx struct {
		notebookCtx context.Context
		service notebook.NotebookSrvCtx
		config config.Application
		log config.Logger
}


func NewNotebook(app *config.Application) INotebook {
	ctx := context.Background()
	return &NotebookCtx {
		notebookCtx: ctx,
		service: *notebook.NotebookServiceInit(app),
		config: *app,
		log: app.Log,
	}
}

func(nCtx *NotebookCtx) CreateNoteBook(w http.ResponseWriter, r *http.Request) {
	claim:= r.Context().Value("claims").(token.CustomClaimData)
	var notebook models.Notebook
	err := validate.GetPayload(r, &notebook)
	if err != nil {
		response.SendWithError(w, http.StatusBadRequest, *errors.BadRequestError(err.Error()))
		return
	}
	notebook.UserId = claim.UserId
	err = nCtx.service.New(nCtx.notebookCtx, &notebook)	
	if err != nil {
		response.SendWithError(w, http.StatusBadRequest, *errors.BadRequestError(err.Error()))
		return
	}
	response.Send(w, http.StatusAccepted, notebook)
}

func (nCtx *NotebookCtx) UpdateNoteBook(w http.ResponseWriter, r *http.Request) {
	
}