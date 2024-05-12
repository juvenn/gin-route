/* Composable and reusable route group for gin
 *
 * Unlike gin.RouterGroup, the defined routes are not bound to gin.Engine until
 * `Dock`, they can be reused and composed among route groups.
 *
 * Use `Mount` to combine routes from multiple groups, and `Dock` to bind to
 * gin.Engine.
 */
package route

import (
	"fmt"
	urlpath "path"

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
func (group *RouteGroup) Use(handlers ...gin.HandlerFunc) *RouteGroup {
	group.handlers = append(group.handlers, handlers...)
	return group
}

// Define a handler for the route.
func (group *RouteGroup) Handle(method string, path string, handlers ...gin.HandlerFunc) {
	group.routes = append(group.routes, &routeItem{
		Method:   method,
		Path:     path,
		Handlers: handlers,
	})
}

// Define collection of route handlers with an ephemeral group.
func (group *RouteGroup) WithScope(path string, fn func(subgroup *RouteGroup)) {
	subgroup := NewGroup(path)
	fn(subgroup)
	group.Mount("", subgroup)
}

// Add routes from sub groups under path prefix.
func (group *RouteGroup) Mount(path string, groups ...*RouteGroup) {
	for _, g := range groups {
		for _, route := range g.Routes() {
			group.routes = append(group.routes, &routeItem{
				Method:   route.Method,
				Path:     urlpath.Join(path, route.Path),
				Handlers: route.Handlers,
			})
		}
	}
}

// Dock and bind routes at path to gin engine.
func (group *RouteGroup) Dock(path string, engine *gin.Engine) {
	eng := engine.Group("")
	for _, route := range group.Routes() {
		fullpath := urlpath.Join(path, route.Path)
		eng.Handle(route.Method, fullpath, route.Handlers...)
	}
}

// Return all route items, with path prefix and middleware concated.
func (group *RouteGroup) Routes() []*routeItem {
	items := make([]*routeItem, 0, len(group.routes))
	group.enumerate(func(item *routeItem) bool {
		items = append(items, item)
		return true
	})

	// for i, route := range ctx.Enumerate {
	// 	items[i] = route
	// }
	return items
}

// Enumerate route items, with path prefix and middleware concated.
// Each route is then applied to func yield. It can be ranged over.
func (group *RouteGroup) enumerate(yield func(*routeItem) bool) {
	for _, route := range group.routes {
		handlers := make([]gin.HandlerFunc, len(group.handlers)+len(route.Handlers))
		copy(handlers, group.handlers)
		copy(handlers[len(group.handlers):], route.Handlers)
		item := &routeItem{
			Method:   route.Method,
			Path:     prependSlash(urlpath.Join(group.pathPrefix, route.Path)),
			Handlers: handlers,
		}
		if !yield(item) {
			return
		}
	}
}

func (group *RouteGroup) printAll() {
	for _, route := range group.Routes() {
		fmt.Printf("%s %s %d\n", route.Method, route.Path, len(route.Handlers))
	}
}

// Prepend slash if not empty
func prependSlash(path string) string {
	if path == "" {
		return ""
	}
	if path[0:1] == "/" {
		return path
	}
	return "/" + path
}
