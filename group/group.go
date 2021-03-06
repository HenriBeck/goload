package group

import (
	"context"
	"fmt"

	"github.com/HenriBeck/goload"
)

// A group of endpoints where you can register new endpoints which should be loadtested using a simple function call.
//
// This pattern has been inspired by the `testing.T.Run` usage where you can have a function which registers sub tests.
type LoadTestGroup interface {
	Test(
		name string,
		rpm int32,
		handler func(ctx context.Context) error,
	)

	SetGroupName(groupName string)
}

type group struct {
	endpoints []goload.Endpoint

	name string
}

type groupEndpoint struct {
	name    string
	rpm     int32
	handler func(ctx context.Context) error
}

func (e *groupEndpoint) Execute(ctx context.Context) error {
	return e.handler(ctx)
}

func (e *groupEndpoint) GetRequestsPerMinute() int32 {
	return e.rpm
}

func (e *groupEndpoint) Name() string {
	return e.name
}

func (g *group) Test(
	name string,
	rpm int32,
	handler func(ctx context.Context) error,
) {
	if g.name != "" {
		g.endpoints = append(g.endpoints, &groupEndpoint{
			name:    fmt.Sprintf("%s/%s", g.name, name),
			rpm:     rpm,
			handler: handler,
		})
	} else {
		g.endpoints = append(g.endpoints, &groupEndpoint{
			name:    name,
			rpm:     rpm,
			handler: handler,
		})
	}
}

func (g *group) SetGroupName(
	groupName string,
) {
	g.name = groupName
}

// `WithGroup` allows a function to be passed
// which then can register endpoints using the `s.Register` function.
//
// This can be used to write `testing` like function which look like:
//
// func LoadTestService(s group.LoadTestGroup) {
//   s.Register("endpoint 1", 20, func(ctx context.Context) error {
//     Code...
//   })
//
//   s.Register("endpoint 2", 30, func(ctx context.Context) error {
//     Code...
//   })
// }
//
// and then registerd using `WithGroup` in the call to `goload.RunLoadtest`.
func WithGroup(s func(group LoadTestGroup)) goload.LoadTestConfig {
	g := &group{}
	s(g)

	return func(options *goload.LoadTestOptions) {
		options.Endpoints = append(options.Endpoints, g.endpoints...)
	}
}
