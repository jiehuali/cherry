package cherryHandler

import (
	"github.com/cherry-game/cherry/const"
	"github.com/cherry-game/cherry/extend/reflect"
	"github.com/cherry-game/cherry/interfaces"
	"github.com/cherry-game/cherry/logger"
	"github.com/cherry-game/cherry/net/message"
	"github.com/cherry-game/cherry/net/route"
)

type (
	//handlerComponent Handler component
	HandlerComponent struct {
		cherryInterfaces.BaseComponent                                      // base component
		HandlerComponentOptions                                             // opts
		handlers                       map[string]cherryInterfaces.IHandler // key:handlerName, value: Handler
	}

	HandlerComponentOptions struct {
		beforeFilters []FilterFn
		afterFilters  []FilterFn
		nameFn        func(string) string
	}

	UnhandledMessage struct {
		Session cherryInterfaces.ISession
		Route   *cherryRoute.Route
		Msg     *cherryMessage.Message
	}

	FilterFn func(msg *UnhandledMessage) bool
)

func NewComponent() *HandlerComponent {
	return &HandlerComponent{
		handlers: make(map[string]cherryInterfaces.IHandler),
		HandlerComponentOptions: HandlerComponentOptions{
			beforeFilters: make([]FilterFn, 0),
			afterFilters:  make([]FilterFn, 0),
			nameFn: func(s string) string {
				return s
			},
		},
	}
}

func (h *HandlerComponent) Name() string {
	return cherryConst.HandlerComponent
}

func (h *HandlerComponent) Init() {
	for _, handler := range h.handlers {
		handler.Set(h.App())
		handler.PreInit()
		handler.Init()
		handler.AfterInit()
	}
}

func (h *HandlerComponent) Stop() {
	for _, handler := range h.handlers {
		handler.Stop()
	}
}

func (h *HandlerComponent) Registers(handlers ...cherryInterfaces.IHandler) {
	for _, handler := range handlers {
		name := handler.Name()
		if name == "" {
			name = cherryReflect.GetStructName(handler)
		}
		h.RegisterWithName(name, handler)
	}
}

func (h *HandlerComponent) RegisterWithName(name string, handler cherryInterfaces.IHandler) {
	if name == "" {
		cherryLogger.Warnf("[Handler= %h] name is empty. skipped.", cherryReflect.GetStructName(handler))
		return
	}

	if handler == nil {
		cherryLogger.Warnf("[Handler= %s] is empty. skipped.", name)
		return
	}

	name = h.nameFn(name)
	if name == "" {
		cherryLogger.Warnf("[Handler= %h] name is empty. skipped.", cherryReflect.GetStructName(handler))
		return
	}

	handler.SetName(name)

	if _, found := h.handlers[name]; found {
		cherryLogger.Errorf("[Handler name = %s] is duplicate!", handler.Name())
		return
	}

	h.handlers[name] = handler
}

func (h *HandlerComponent) DoHandle(msg *UnhandledMessage) {
	if msg == nil || msg.Route == nil {
		return
	}

	if msg.Route.NodeType() != h.App().NodeType() {
		return
	}

	if !h.App().Running() {
		//ignore message
		return
	}

	handler := h.GetHandler(msg.Route)
	if handler == nil {
		cherryLogger.Errorf("[Route = %h] not found handler.", msg.Route)
		return
	}

	handler.PutMessage(msg)
}

func (h *HandlerComponent) GetHandler(route *cherryRoute.Route) cherryInterfaces.IHandler {
	handler := h.handlers[h.nameFn(route.HandlerName())]
	if handler == nil {
		cherryLogger.Warnf("could not find handle worker for Route = %v", route)
		return nil
	}
	return handler
}

// PostEvent 发布事件
func (h *HandlerComponent) PostEvent(event cherryInterfaces.IEvent) {
	if event == nil {
		return
	}

	for _, handler := range h.handlers {
		if _, found := handler.GetEvent(event.EventName()); found {
			handler.PutMessage(event)
		}
	}
}

func (c *HandlerComponentOptions) GetBeforeFilter() []FilterFn {
	return c.beforeFilters
}

func (c *HandlerComponentOptions) BeforeFilter(beforeFilters ...FilterFn) {
	if len(beforeFilters) < 1 {
		return
	}
	c.beforeFilters = append(c.beforeFilters, beforeFilters...)
}

func (c *HandlerComponentOptions) GetAfterFilter() []FilterFn {
	return c.afterFilters
}

func (c *HandlerComponentOptions) AfterFilter(afterFilters ...FilterFn) {
	if len(afterFilters) < 1 {
		return
	}
	c.afterFilters = append(c.afterFilters, afterFilters...)
}

func (c *HandlerComponentOptions) SetNameFn(fn func(string) string) {
	if fn == nil {
		return
	}
	c.nameFn = fn
}

// NodeRoute  结点路由规则 nodeType:结点类型,routeFunc 路由规则
func (*HandlerComponentOptions) NodeRoute(nodeType string, routeFunc cherryInterfaces.RouteFunction) {
	cherryLogger.Panic(nodeType, routeFunc)
}
