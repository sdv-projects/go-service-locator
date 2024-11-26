package sl

import (
	"context"
	"reflect"

	"github.com/google/uuid"
)

var _ ServiceLocator = (*locator)(nil)

const (
	globalScopeID = "global"
)

type contextKey string

var (
	ctxScopeID = contextKey("ScopeID")
)

type provideType byte

const (
	singleton provideType = iota
	transient
	scoped
)

type serviceProvide struct {
	Type    provideType
	Provide ProvideServiceFn
}

type locator struct {
	scopes map[string]map[reflect.Type]any

	pr map[reflect.Type]serviceProvide
}

func (l *locator) getScopeID(ctx context.Context) string {
	if id, ok := ctx.Value(ctxScopeID).(string); ok {
		return id
	}

	return globalScopeID
}

func (l *locator) register(t reflect.Type, pt provideType, pf ProvideServiceFn) {
	l.pr[t] = serviceProvide{
		Type:    pt,
		Provide: pf,
	}
}

func (l *locator) OpenScope(ctx context.Context) context.Context {
	scopeID := uuid.NewString()

	l.scopes[scopeID] = make(map[reflect.Type]any)

	return context.WithValue(ctx, ctxScopeID, scopeID)
}

func (l *locator) CloseScope(ctx context.Context) {
	scopeID := l.getScopeID(ctx)
	if scopeID == globalScopeID {
		return
	}

	delete(l.scopes, scopeID)
}

func (l *locator) AddSingleton(t reflect.Type, f ProvideServiceFn) {
	l.register(t, singleton, f)
}

func (l *locator) AddTransient(t reflect.Type, f ProvideServiceFn) {
	l.register(t, transient, f)
}

func (l *locator) AddScoped(t reflect.Type, f ProvideServiceFn) {
	l.register(t, scoped, f)
}

func (l *locator) GetService(ctx context.Context, t reflect.Type) (service any) {
	p, ok := l.pr[t]
	if !ok {
		return nil
	}

	if p.Type == transient {
		return p.Provide(ctx, l)
	}

	var scopeID string
	if p.Type == singleton {
		scopeID = globalScopeID
	} else {
		scopeID, ok = ctx.Value(ctxScopeID).(string)
		if !ok {
			return nil
		}
	}

	s, ok := l.scopes[scopeID][t]
	if !ok {
		s = p.Provide(ctx, l)
		l.scopes[scopeID][t] = s
	}

	return s
}

func NewServiceLocator() *locator {
	return &locator{
		scopes: make(map[string]map[reflect.Type]any),
		pr:     make(map[reflect.Type]serviceProvide),
	}
}
