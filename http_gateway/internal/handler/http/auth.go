package http

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
)

func forwardHeaderToClient() runtime.ServeMuxOption {

	return runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
		header := request.Header.Get("Authorization")
		// send all the headers received from the client
		return metadata.Pairs("Authorization", header)
	})
}

func withErrorHandler() runtime.ServeMuxOption {
	return runtime.WithErrorHandler(func(ctx context.Context,
		mux *runtime.ServeMux,
		marshaler runtime.Marshaler,
		writer http.ResponseWriter,
		request *http.Request, err error,
	) {
		//creating a new HTTTPStatusError with a custom status, and passing error
		newError := runtime.HTTPStatusError{
			HTTPStatus: 400,
			Err:        err,
		}
		// using default handler to do the rest of heavy lifting of marshaling error and adding headers
		runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, writer, request, &newError)
	})
}

func withMetadata() runtime.ServeMuxOption {
	return runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
		// Step 2 - Extend the context
		_ = metadata.AppendToOutgoingContext(ctx)
		// Step 3 - get the basic auth params
		username, password, ok := request.BasicAuth()
		if !ok {
			return nil
		}
		md := metadata.Pairs()
		md.Append("Username", username)
		md.Append("Password", password)
		return md
	})
}
