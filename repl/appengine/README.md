# CEL REPL Web
Provides a simple web interface for working with the CEL REPL.

The two main components are 1) a Go application that implements a service
wrapping the REPL and serving static web content 2) an angular based web page
that provides an interface around the API.

The REPL service is stateless -- it initializes a new REPL instance and applies
the requested commands in order on each request.

## Development

To run the application in development, run the angular CLI builder in watch mode
and run the Go server as follows:

```
# from the `repl/appengine/web` directory:
ng build --watch

# from the repl/appengine directory:
go run ./main --serve_static ./web/dist/web
```

## Deploy on google cloud appengine

1. Build the angular application with `ng build`.

1. Follow the instructions here:
(https://cloud.google.com/appengine/docs/standard/go/building-app). Make sure to
follow the instructions for setting up your gcloud cli, and the appengine
support in cloud console for your project.