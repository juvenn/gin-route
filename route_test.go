package route

import (
	"github.com/gin-gonic/gin"
)

func ExampleRouteGroup_Handle() {
	photos := NewGroup("photos")
	photos.Handle("GET", "", queryPhotos)
	photos.Handle("GET", "/:id", getPhoto)
	photos.printAll()
	// Output: GET /photos 1
	// GET /photos/:id 1
}

func ExampleRouteGroup_Use() {
	photos := NewGroup("photos")
	photos.Use(logMiddleware())
	photos.Handle("GET", "", queryPhotos)
	photos.Handle("GET", "/:id", getPhoto)
	photos.printAll()
	// Output: GET /photos 2
	// GET /photos/:id 2
}

func ExampleRouteGroup_Mount() {
	v1 := NewGroup("v1")
	photos := NewGroup("photos")
	photos.Handle("GET", "", queryPhotos)
	photos.Handle("GET", "/:id", getPhoto)
	v1.Mount("", photos)
	v1.printAll()
	// Output: GET /v1/photos 1
	// GET /v1/photos/:id 1
}

func ExampleRouteGroup_With() {
	v1 := NewGroup("/v1")
	v1.With("photos", func(photos *RouteGroup) {
		photos.Use(logMiddleware())
		photos.Handle("GET", "", logMiddleware(), queryPhotos)
		photos.Handle("GET", "/:id", getPhoto)
	})
	v1.printAll()
	// Output: GET /v1/photos 3
	// GET /v1/photos/:id 2
}

func logMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
func queryPhotos(ctx *gin.Context) {}
func getPhoto(ctx *gin.Context)    {}
