## Composable route groups for Gin

Composable route groups for [Gin](https://github.com/gin-gonic/gin).

A route group defines a set of routes that share the same path prefix and middlewares. The route items defined in such group can be enumerated, and composed among other groups.

### Usage

```sh
go get github.com/juvenn/gin-route
```

```go
import (
  "github.com/gin-gonic/gin"
  route "github.com/juvenn/gin-route"
)

func main() {
  photos := route.NewGroup("/photos")
  photos.WithScope("", func(g *route.RouteGroup) {
    g.Use(logMiddleware)
    g.Handle("GET", "", queryPhotos)
    g.Handle("GET", "/:id", getPhoto)
    g.Handle("DELETE", "/:id", deletePhoto)
  })
  // Enumerate and print all routes
  for _, route := range photos.Routes() {
		fmt.Printf("%s %s %d\n", route.Method, route.Path, len(route.Handlers))
	}

  // dock to gin engine
  r := gin.Default()
  photos.Dock("/v1", r)
  photos.Dock("/v2", r)
  r.Run()
}
```
