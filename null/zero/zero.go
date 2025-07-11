package zero

import "github.com/guregu/null/v6/zero"

var (
	BoolFrom    = zero.BoolFrom
	BoolFromPtr = zero.BoolFromPtr
	NewBool     = zero.NewBool

	ByteFrom    = zero.ByteFrom
	ByteFromPtr = zero.ByteFromPtr
	NewByte     = zero.NewByte

	FloatFrom    = zero.FloatFrom
	FloatFromPtr = zero.FloatFromPtr
	NewFloat     = zero.NewFloat

	IntFrom    = zero.IntFrom
	IntFromPtr = zero.IntFromPtr
	NewInt     = zero.NewInt

	Int16From    = zero.Int16From
	Int16FromPtr = zero.Int16FromPtr
	NewInt16     = zero.NewInt16

	Int32From    = zero.Int32From
	Int32FromPtr = zero.Int32FromPtr
	NewInt32     = zero.NewInt32

	StringFrom    = zero.StringFrom
	StringFromPtr = zero.StringFromPtr
	NewString     = zero.NewString

	TimeFrom    = zero.TimeFrom
	TimeFromPtr = zero.TimeFromPtr
	NewTime     = zero.NewTime

	ValueFrom    = zero.ValueFrom
	ValueFromPtr = zero.ValueFromPtr
	NewValue     = zero.NewValue
)

type (
	Bool                = zero.Bool
	Byte                = zero.Byte
	Float               = zero.Float
	Int                 = zero.Int
	Int16               = zero.Int16
	Int32               = zero.Int32
	Int64               = zero.Int64
	String              = zero.String
	Time                = zero.Time
	Value[T comparable] = zero.Value[T]
)
