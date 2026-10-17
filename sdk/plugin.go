// Package sdk is the public plugin API. It does not import the application or its internals.
package sdk

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/events"
	"github.com/dingding229/mediatool-plugin/sdk/host"
)

const ContractVersion = contracts.ContractVersion

type Config = contracts.Config
type Form = contracts.Form
type Page = contracts.Page
type PageRequest = contracts.PageRequest
type Dashboard = contracts.Dashboard
type DashboardRequest = contracts.DashboardRequest
type DashboardMeta = contracts.DashboardMeta
type API = contracts.Descriptor
type Service = contracts.Descriptor
type Command = contracts.Descriptor
type Module = contracts.Descriptor
type Action = contracts.Descriptor

type Plugin interface {
	InitPlugin(context.Context, Config) error
	GetState(context.Context) (bool, error)
	StopService(context.Context) error
}
type NameProvider interface {
	GetName(context.Context) (string, error)
}
type FormProvider interface {
	GetForm(context.Context) (Form, error)
}
type PageProvider interface {
	GetPage(context.Context, PageRequest) (Page, error)
}
type APIProvider interface {
	GetAPI(context.Context) ([]API, error)
}
type ServiceProvider interface {
	GetService(context.Context) ([]Service, error)
}
type CommandProvider interface {
	GetCommand(context.Context) ([]Command, error)
}
type DashboardProvider interface {
	GetDashboard(context.Context, DashboardRequest) (Dashboard, error)
}
type DashboardMetaProvider interface {
	GetDashboardMeta(context.Context) ([]DashboardMeta, error)
}
type ModuleProvider interface {
	GetModule(context.Context) ([]Module, error)
}
type ActionProvider interface {
	GetActions(context.Context) ([]Action, error)
}
type Services struct {
	Host       *host.Client
	Events     *events.Registry
	InstanceID string
}

// Factory must not be called by catalog enumeration or package initialization.
type Factory func(Services) Plugin
