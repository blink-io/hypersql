package null

import "github.com/guregu/null/v6"

var (
	BoolFrom    = null.BoolFrom
	BoolFromPtr = null.BoolFromPtr
	NewBool     = null.NewBool

	ByteFrom    = null.ByteFrom
	ByteFromPtr = null.ByteFromPtr
	NewByte     = null.NewByte

	FloatFrom    = null.FloatFrom
	FloatFromPtr = null.FloatFromPtr
	NewFloat     = null.NewFloat

	IntFrom    = null.IntFrom
	IntFromPtr = null.IntFromPtr
	NewInt     = null.NewInt

	Int16From    = null.Int16From
	Int16FromPtr = null.Int16FromPtr
	NewInt16     = null.NewInt16

	Int32From    = null.Int32From
	Int32FromPtr = null.Int32FromPtr
	NewInt32     = null.NewInt32

	StringFrom    = null.StringFrom
	StringFromPtr = null.StringFromPtr
	NewString     = null.NewString

	TimeFrom    = null.TimeFrom
	TimeFromPtr = null.TimeFromPtr
	NewTime     = null.NewTime

	ValueFrom    = null.ValueFrom
	ValueFromPtr = null.ValueFromPtr
	NewValue     = null.NewValue

	NewUint     = null.NewValue[uint]
	UintFrom    = null.ValueFrom[uint]
	UintFromPtr = null.ValueFromPtr[uint]
)

type (
	Bool                = null.Bool
	Byte                = null.Byte
	String              = null.String
	Time                = null.Time
	Float               = null.Float
	Value[T comparable] = null.Value[T]
	Int                 = null.Int
	Int16               = null.Int16
	Int32               = null.Int32
	Int64               = null.Int64
)
