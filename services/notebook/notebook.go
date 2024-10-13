package notebook

import (
	"context"

	"github.com/imabg/sync/internal/models"
	"github.com/imabg/sync/pkg/config"
)

type NotebookSrvCtx struct {
	notebookModel *models.NotebookCtx
	config config.Application
	log config.Logger
}

func NotebookServiceInit(app *config.Application) *NotebookSrvCtx {
	return &NotebookSrvCtx{
		notebookModel: models.NewNotebookModel(*app.MongoClient),
		config: *app,
		log: app.Log,
	}
}

func(n *NotebookSrvCtx) New(ctx context.Context, notebook *models.Notebook) error {
	return n.notebookModel.Create(ctx, notebook)
}