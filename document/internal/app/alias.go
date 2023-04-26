package app

import (
	documentApp "github.com/0B1t322/Documents-Service/document/internal/app/documents"
	elementsApp "github.com/0B1t322/Documents-Service/document/internal/app/elements"
	stylesApp "github.com/0B1t322/Documents-Service/document/internal/app/styles"
	"github.com/0B1t322/Documents-Service/internal/core/events/amqp"
)

// Alias apps
type (
	Documents = documentApp.App
	Elements  = elementsApp.App
	Styles    = stylesApp.App
)

// Alias apps constructor
var (
	NewDocumentApp = documentApp.New
	NewElementsApp = elementsApp.New
	NewStylesApp   = stylesApp.New
)

// Alias repository
type (
	DocumentsRepository = documentApp.Repository
	ElementsRepository  = elementsApp.Repository
	StylesRepository    = stylesApp.Repository
)

// Alias event publisher
var (
	NewAMQPEventPublisher = amqp.NewEventsPublisher
)
