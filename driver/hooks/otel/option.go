package otel

import (
	"github.com/XSAM/otelsql"
)

var (
	WithAttributes                 = otelsql.WithAttributes
	WithAttributesGetter           = otelsql.WithAttributesGetter
	WithInstrumentAttributesGetter = otelsql.WithInstrumentAttributesGetter
	WithMeterProvider              = otelsql.WithMeterProvider
	WithSQLCommenter               = otelsql.WithSQLCommenter
	WithSpanNameFormatter          = otelsql.WithSpanNameFormatter
	WithSpanOptions                = otelsql.WithSpanOptions
	WithTracerProvider             = otelsql.WithTracerProvider
)

type (
	Option = otelsql.Option
)
