/* Composable and reusable route group for gin
 *
 * Unlike gin.RouterGroup, the defined routes are not bound to gin.Engine until
 * `Dock`, they can be reused and composed amoung route groups.
 *
 * Use `Mount` to combine routes from multiple groups, and `Dock` to bind to
 * gin.Engine.
 */
package route

import (
	"fmt"
	xpath "path"

	"github.com/gin-gonic/gin"
)

type routeItem struct {
	Method   string
	Path     string
	Handlers gin.HandlersChain
}

type RouteGroup struct {
	routes     []*routeItem
	pathPrefix string
	handlers   gin.HandlersChain // middleware handlers for the group
}

// Build a route group with path prefix, and a chain of middleware handlers.
//
// The path and handlers are applied lazily until mount or dock time.
func NewGroup(pathPrefix string, handlers ...gin.HandlerFunc) *RouteGroup {
	return &RouteGroup{
		pathPrefix: pathPrefix,
		handlers:   handlers,
	}
}

// Append more middleware handlers to the group.
func (ctx *RouteGroup) Use(handlers ...gin.HandlerFunc) *RouteGroup {
	ctx.handlers = append(ctx.handlers, handlers...)
	return ctx
}

// Define a handler for the route.
func (ctx *RouteGroup) Handle(method string, path string, handlers ...gin.HandlerFunc) {
	ctx.routes = append(ctx.routes, &routeItem{
		Method:   method,
		Path:     path,
		Handlers: handlers,
	})
}

// Define collection of route handlers with an ephemeral group.
func (ctx *RouteGroup) With(path string, fn func(group *RouteGroup)) {
	group := NewGroup(path)
	fn(group)
	ctx.Mount("", group)
}

// Add routes from sub groups under path prefix.
func (ctx *RouteGroup) Mount(path string, groups ...*RouteGroup) {
	for _, group := range groups {
		for _, route := range group.Routes() {
			ctx.routes = append(ctx.routes, &routeItem{
				Method:   route.Method,
				Path:     path + route.Path,
				Handlers: route.Handlers,
			})
		}
	}
}

// Dock and bind routes at path to gin engine.
func (ctx *RouteGroup) Dock(path string, engine *gin.Engine) {
	eng := engine.Group("")
	for _, route := range ctx.Routes() {
		fullpath := prependSlash(xpath.Join(path, route.Path))
		eng.Handle(route.Method, fullpath, route.Handlers...)
	}
}

// Enumerate route items and concat path prefix and middleware handlers.
func (ctx *RouteGroup) Routes() []*routeItem {
	items := make([]*routeItem, 0, len(ctx.routes))
	ctx.enumerate(func(item *routeItem) bool {
		items = append(items, item)
		return true
	})

	// for i, route := range ctx.Enumerate {
	// 	items[i] = route
	// }
	return items
}

// enumerate route items and concat path prefix and middleware handlers.
// Each route is then applied to func yield. It can be ranged over.
func (ctx *RouteGroup) enumerate(yield func(*routeItem) bool) {
	for _, route := range ctx.routes {
		handlers := make([]gin.HandlerFunc, len(ctx.handlers)+len(route.Handlers))
		copy(handlers, ctx.handlers)
		copy(handlers[len(ctx.handlers):], route.Handlers)
		item := &routeItem{
			Method:   route.Method,
			Path:     prependSlash(xpath.Join(ctx.pathPrefix, route.Path)),
			Handlers: handlers,
		}
		if !yield(item) {
			return
		}
	}
}

func (ctx *RouteGroup) printAll() {
	for _, route := range ctx.Routes() {
		fmt.Printf("%s %s %d\n", route.Method, route.Path, len(route.Handlers))
	}
}

func prependSlash(path string) string {
	if path == "" || path == "/" {
		return ""
	}
	if path[0:1] == "/" {
		return path
	}
	return "/" + path
}
