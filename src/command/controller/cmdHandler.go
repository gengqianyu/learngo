package controller

import (
	"command/common"
	"command/model"
	"log"
	"net/http"
)

type CmdHandler struct {
	model *model.Command
}

func CreateCmdHandler() *CmdHandler {
	return &CmdHandler{
		model: model.CreateCommand(),
	}
}

func (c *CmdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(c.model)
	err := common.View("command").Render(w, c.model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
