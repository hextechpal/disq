package server

import (
	"fmt"
	"github.com/ppal31/disq/internal/storage"
	"github.com/valyala/fasthttp"
	"strconv"
)

type StorageServer struct {
	s storage.Storage
}

func (ss *StorageServer) Start() error {
	return fasthttp.ListenAndServe(":8081", ss.requestHandler)
}

func (ss *StorageServer) requestHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/send":
		ss.send(ctx)
	case "/receive":
		ss.receive(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

func (ss *StorageServer) send(ctx *fasthttp.RequestCtx) {
	err := ss.s.Send(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(err.Error())
	}
}

func (ss *StorageServer) receive(ctx *fasthttp.RequestCtx) {
	offset, err := ctx.QueryArgs().GetUint("offset")
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(fmt.Sprintf("bad `off` GET param: %v", err))
		return
	}
	maxSize, err := ctx.QueryArgs().GetUint("maxSize")
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(fmt.Sprintf("bad `off` GET param: %v", err))
		return
	}

	newOffset, err := ss.s.Receive(uint64(offset), uint64(maxSize), ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}
	ctx.Response.Header.Set("offset", strconv.FormatUint(newOffset, 10))

}
