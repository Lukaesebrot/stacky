package v1

import (
	"github.com/Lukaesebrot/stacky/stacks"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// endPutStackHost handles the PUT '/stacks/{name}/hosts' endpoint
func endPutStackHost(ctx *fasthttp.RequestCtx) {
	// Retrieve the stack
	stack, err := stacks.Retrieve(ctx.UserValue("name").(string))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetBody(errorResponse(fasthttp.StatusNotFound, "the requested stack couldn't be found", nil))
			return
		}
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody(errorResponse(fasthttp.StatusInternalServerError, err.Error(), nil))
		return
	}

	// Add the given host to the current stack
	host := string(ctx.QueryArgs().Peek("host"))
	if host == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody(errorResponse(fasthttp.StatusBadRequest, "no empty fields allowed", nil))
		return
	}
	err = stack.AddHost(host)
	if err != nil {
		if err == stacks.ErrHostAlreadyExists {
			ctx.SetStatusCode(fasthttp.StatusUnprocessableEntity)
			ctx.SetBody(errorResponse(fasthttp.StatusUnprocessableEntity, err.Error(), nil))
			return
		}
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody(errorResponse(fasthttp.StatusInternalServerError, err.Error(), nil))
		return
	}

	// Respond with a success message
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(successResponse(fasthttp.StatusOK, "the given host was added to the current stack", nil))
}
