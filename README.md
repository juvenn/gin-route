## Composable route groups for Gin

Unlike gin.RouterGroup, this package allows you to compose routes from multiple
groups.

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
  r := gin.Default()
  photos := route.NewGroup("/photos")
  photos.WithContext("", func(g *route.RouteGroup) {
    g.Use()
    g.Handle("GET", "", queryPhotos)
    g.Handle("GET", "/:id", getPhoto)
    g.Handle("DELETE", "/:id", deletePhoto)
  })
  photos.Dock("/v1", r)
  photos.Dock("/v2", r)
  r.Run()
}
```
