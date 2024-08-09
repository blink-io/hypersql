package scan

import "github.com/stephenafamo/scan"

var (
	CtxKeyAllowUnknownColumns = scan.CtxKeyAllowUnknownColumns
)

type (
	AfterMod            = scan.AfterMod
	BeforeFunc          = scan.BeforeFunc
	MapperMod           = scan.MapperMod
	MappingError        = scan.MappingError
	MappingOption       = scan.MappingOption
	MappingSourceOption = scan.MappingSourceOption
	Queryer             = scan.Queryer
	Row                 = scan.Row
	RowValidator        = scan.RowValidator
	Rows                = scan.Rows
	StructMapperSource  = scan.StructMapperSource
	TypeConverter       = scan.TypeConverter

	ICursor[T any] interface {
		scan.ICursor[T]
	}
)
