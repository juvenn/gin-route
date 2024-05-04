/* Composable and reusable route group for gin
 *
 * Unlike gin.RouterGroup, the defined routes are not bound to gin.Engine until
 * `Dock`, they can be reused and composed amoung route groups.
 *
 * Use `Mount` to combine routes from multiple groups, and `Dock` to bind to
 * gin.Engine.
 */
package ginroute

import (
	"github.com/gin-gonic/gin"
)

type routeItem struct {
	Method   string
	Path     string
	Handlers gin.HandlersChain
}

type routeGroup struct {
	routes     []*routeItem
	pathPrefix string
	handlers   gin.HandlersChain // middleware handlers for the group
}

// Build a route group with path prefix, and a chain of middleware handlers.
//
// The path and handlers are applied lazily until mount or dock time.
func NewGroup(pathPrefix string, handlers ...gin.HandlerFunc) *routeGroup {
	return &routeGroup{
		pathPrefix: pathPrefix,
		handlers:   handlers,
	}
}

// Append more middleware handlers to the group.
func (ctx *routeGroup) Use(handlers ...gin.HandlerFunc) {
	ctx.handlers = append(ctx.handlers, handlers...)
}

// Define a handler for the route.
func (ctx *routeGroup) Handle(method string, path string, handlers ...gin.HandlerFunc) {
	ctx.routes = append(ctx.routes, &routeItem{
		Method:   method,
		Path:     path,
		Handlers: handlers,
	})
}

// Define collection of route handlers with an ephemeral group.
func (ctx *routeGroup) WithContext(path string, fn func(newctx *routeGroup)) {
	newctx := NewGroup(path)
	fn(newctx)
	ctx.Mount("", newctx)
}

// Mount routes from sub groups under path prefix.
func (ctx *routeGroup) Mount(path string, groups ...*routeGroup) {
	for _, group := range groups {
		for _, route := range group.Routes() {
			route.Path = path + route.Path
			ctx.routes = append(ctx.routes, route)
		}
	}
}

// Dock and bind routes to gin engine.
func (ctx *routeGroup) Dock(engine *gin.Engine) {
	group := engine.Group("/")
	for _, route := range ctx.Routes() {
		// fmt.Printf("Define %s %s\n", route.Method, route.Path)
		group.Handle(route.Method, route.Path, route.Handlers...)
	}
}

// Enumerate route items and apply path prefix and middleware handlers.
func (ctx *routeGroup) Routes() []*routeItem {
	items := make([]*routeItem, len(ctx.routes))
	for i, route := range ctx.routes {
		handlers := make([]gin.HandlerFunc, len(ctx.handlers)+len(route.Handlers))
		copy(handlers, ctx.handlers)
		copy(handlers[len(ctx.handlers):], route.Handlers)
		items[i] = &routeItem{
			Method:   route.Method,
			Path:     ctx.pathPrefix + route.Path,
			Handlers: handlers,
		}
	}
	return items
}
