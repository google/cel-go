package checker

import (
	"celgo/operators"
	"celgo/semantics"
	"celgo/semantics/types"
)

func AddStandard(env *Env) {
	// Some shortcuts we use when building declarations.
	paramA := types.NewTypeParam("A")
	typeParamAList := []string{"A"}
	listOfA := types.NewList(paramA)
	paramB := types.NewTypeParam("B")
	typeParamABList := []string{"A", "B"}
	mapOfAB := types.NewMap(paramA, paramB)

	// Booleans
	env.AddFunction(
		semantics.NewFunction(operators.Conditional,
			semantics.NewParameterizedOverload("conditional", false, typeParamAList, paramA, types.Bool, paramA, paramA)))
	env.AddFunction(
		semantics.NewFunction(operators.LogicalAnd,
			semantics.NewOverload("logical_and", false, types.Bool, types.Bool, types.Bool)))
	env.AddFunction(
		semantics.NewFunction(operators.LogicalOr,
			semantics.NewOverload("logical_or", false, types.Bool, types.Bool, types.Bool)))
	env.AddFunction(
		semantics.NewFunction(operators.LogicalNot,
			semantics.NewOverload("logical_not", false, types.Bool, types.Bool)))
	env.AddFunction(
		semantics.NewFunction("matches",
			semantics.NewOverload("matches", false, types.Bool, types.String, types.String)))

	// Relations
	env.AddFunction(
		semantics.NewFunction(operators.Less,
			semantics.NewOverload("less_bool", false, types.Bool, types.Bool, types.Bool),
			semantics.NewOverload("less_int64", false, types.Bool, types.Int64, types.Int64),
			semantics.NewOverload("less_uint64", false, types.Bool, types.Uint64, types.Uint64),
			semantics.NewOverload("less_double", false, types.Bool, types.Double, types.Double),
			semantics.NewOverload("less_string", false, types.Bool, types.String, types.String),
			semantics.NewOverload("less_bytes", false, types.Bool, types.Bytes, types.Bytes)))
	env.AddFunction(
		semantics.NewFunction(operators.LessEquals,
			semantics.NewOverload("less_equals_bool", false, types.Bool, types.Bool, types.Bool),
			semantics.NewOverload("less_equals_int64", false, types.Bool, types.Int64, types.Int64),
			semantics.NewOverload("less_equals_uint64", false, types.Bool, types.Uint64, types.Uint64),
			semantics.NewOverload("less_equals_double", false, types.Bool, types.Double, types.Double),
			semantics.NewOverload("less_equals_string", false, types.Bool, types.String, types.String),
			semantics.NewOverload("less_equals_bytes", false, types.Bool, types.Bytes, types.Bytes)))
	env.AddFunction(
		semantics.NewFunction(operators.Greater,
			semantics.NewOverload("greater_bool", false, types.Bool, types.Bool, types.Bool),
			semantics.NewOverload("greater_int64", false, types.Bool, types.Int64, types.Int64),
			semantics.NewOverload("greater_uint64", false, types.Bool, types.Uint64, types.Uint64),
			semantics.NewOverload("greater_double", false, types.Bool, types.Double, types.Double),
			semantics.NewOverload("greater_string", false, types.Bool, types.String, types.String),
			semantics.NewOverload("greater_bytes", false, types.Bool, types.Bytes, types.Bytes)))
	env.AddFunction(
		semantics.NewFunction(operators.GreaterEquals,
			semantics.NewOverload("greater_equals_bool", false, types.Bool, types.Bool, types.Bool),
			semantics.NewOverload("greater_equals_int64", false, types.Bool, types.Int64, types.Int64),
			semantics.NewOverload("greater_equals_uint64", false, types.Bool, types.Uint64, types.Uint64),
			semantics.NewOverload("greater_equals_double", false, types.Bool, types.Double, types.Double),
			semantics.NewOverload("greater_equals_string", false, types.Bool, types.String, types.String),
			semantics.NewOverload("greater_equals_bytes", false, types.Bool, types.Bytes, types.Bytes)))
	env.AddFunction(
		semantics.NewFunction(operators.Equals,
			semantics.NewParameterizedOverload("equals", false, typeParamAList, types.Bool, paramA, paramA)))
	env.AddFunction(
		semantics.NewFunction(operators.NotEquals,
			semantics.NewParameterizedOverload("not_equals", false, typeParamAList, types.Bool, paramA, paramA)))

	// Algebra
	env.AddFunction(
		semantics.NewFunction(operators.Subtract,
			semantics.NewOverload("subtract_int64", false, types.Int64, types.Int64, types.Int64),
			semantics.NewOverload("subtract_uint64", false, types.Uint64, types.Uint64, types.Uint64),
			semantics.NewOverload("subtract_double", false, types.Double, types.Double, types.Double)))
	env.AddFunction(
		semantics.NewFunction(operators.Multiply,
			semantics.NewOverload("multiply_int64", false, types.Int64, types.Int64, types.Int64),
			semantics.NewOverload("multiply_uint64", false, types.Uint64, types.Uint64, types.Uint64),
			semantics.NewOverload("multiply_double", false, types.Double, types.Double, types.Double)))
	env.AddFunction(
		semantics.NewFunction(operators.Divide,
			semantics.NewOverload("divide_int64", false, types.Int64, types.Int64, types.Int64),
			semantics.NewOverload("divide_uint64", false, types.Uint64, types.Uint64, types.Uint64),
			semantics.NewOverload("divide_double", false, types.Double, types.Double, types.Double)))
	env.AddFunction(
		semantics.NewFunction(operators.Modulo,
			semantics.NewOverload("modulo_int64", false, types.Int64, types.Int64, types.Int64),
			semantics.NewOverload("modulo_uint64", false, types.Uint64, types.Uint64, types.Uint64)))
	env.AddFunction(
		semantics.NewFunction(operators.Add,
			semantics.NewOverload("add_int64", false, types.Int64, types.Int64, types.Int64),
			semantics.NewOverload("add_uint64", false, types.Uint64, types.Uint64, types.Uint64),
			semantics.NewOverload("add_double", false, types.Double, types.Double, types.Double),
			semantics.NewOverload("add_string", false, types.String, types.String, types.String),
			semantics.NewOverload("add_bytes", false, types.Bytes, types.Bytes, types.Bytes),
			semantics.NewParameterizedOverload("add_list", false, typeParamAList, listOfA, listOfA, listOfA)))
	env.AddFunction(
		semantics.NewFunction(operators.Negate,
			semantics.NewOverload("negate_int64", false, types.Int64, types.Int64),
			semantics.NewOverload("negate_double", false, types.Double, types.Double)))

	// Index
	env.AddFunction(
		semantics.NewFunction(operators.Index,
			semantics.NewParameterizedOverload("index_list", false, typeParamAList, paramA, listOfA, types.Int64),
			semantics.NewParameterizedOverload("index_map", false, typeParamABList, paramB, mapOfAB, paramA)))

	// Collections
	env.AddFunction(
		semantics.NewFunction("size",
			semantics.NewOverload("size_string", false, types.Int64, types.String),
			semantics.NewOverload("size_bytes", false, types.Int64, types.Bytes),
			semantics.NewParameterizedOverload("size_list", false, typeParamAList, types.Int64, listOfA),
			semantics.NewParameterizedOverload("size_map", false, typeParamABList, types.Int64, mapOfAB)))

	env.AddFunction(
		semantics.NewFunction(operators.In,
			semantics.NewParameterizedOverload("in_list", false, typeParamAList, types.Bool, paramA, listOfA),
			semantics.NewParameterizedOverload("in_map", false, typeParamABList, types.Bool, paramA, mapOfAB)))

	for _, t := range []types.Type{types.Int64, types.Uint64, types.Bool, types.Double, types.Bytes, types.String} {
		env.AddIdent(semantics.NewIdent(t.String(), types.NewTypeType(t), nil))
	}

	env.AddFunction(
		semantics.NewFunction("list",
			semantics.NewParameterizedOverload("list_type", false, typeParamAList, types.NewTypeType(listOfA), types.NewTypeType(paramA))))

	env.AddFunction(
		semantics.NewFunction("map",
			semantics.NewParameterizedOverload(
				"map_type", false, typeParamABList, types.NewTypeType(mapOfAB), types.NewTypeType(paramA), types.NewTypeType(paramB))))

	env.AddFunction(
		semantics.NewFunction("type",
			semantics.NewParameterizedOverload("type", false, typeParamAList, types.NewTypeType(paramA), paramA)))

	// Conversions to int
	env.AddFunction(
		semantics.NewFunction("int",
			semantics.NewOverload("uint64_to_int64", false, types.Int64, types.Uint64),
			semantics.NewOverload("double_to_int64", false, types.Int64, types.Double),
			semantics.NewOverload("string_to_int64", false, types.Int64, types.String)))

	// Conversions to uint
	env.AddFunction(
		semantics.NewFunction("uint",
			semantics.NewOverload("int64_to_uint64", false, types.Uint64, types.Int64),
			semantics.NewOverload("double_to_uint64", false, types.Uint64, types.Double),
			semantics.NewOverload("string_to_uint64", false, types.Uint64, types.String)))

	// Conversions to double
	env.AddFunction(
		semantics.NewFunction("double",
			semantics.NewOverload("int64_to_double", false, types.Double, types.Int64),
			semantics.NewOverload("uint64_to_double", false, types.Double, types.Uint64),
			semantics.NewOverload("string_to_double", false, types.Double, types.String)))

	// Conversions to string
	env.AddFunction(
		semantics.NewFunction("string",
			semantics.NewOverload("int64_to_string", false, types.String, types.Int64),
			semantics.NewOverload("uint64_to_string", false, types.String, types.Uint64),
			semantics.NewOverload("double_to_string", false, types.String, types.Double),
			semantics.NewOverload("bytes_to_string", false, types.String, types.Bytes)))

	// Conversions to list
	env.AddFunction(
		semantics.NewFunction("list",
			semantics.NewParameterizedOverload(
				"to_list", false, typeParamAList, listOfA, types.NewTypeType(paramA), listOfA)))

	// Conversions to map
	env.AddFunction(
		semantics.NewFunction("map",
			semantics.NewParameterizedOverload(
				"to_map", false, typeParamABList, mapOfAB, types.NewTypeType(paramA), types.NewTypeType(paramB), mapOfAB)))

	// Conversions to bytes
	env.AddFunction(
		semantics.NewFunction("bytes",
			semantics.NewOverload("string_to_bytes", false, types.Bytes, types.String)))

	// Conversions to dyn
	env.AddFunction(
		semantics.NewFunction("dyn",
			semantics.NewParameterizedOverload("to_dyn", false, typeParamAList, types.Dynamic, paramA)))

	// Date/time functions
	env.AddFunction(
		semantics.NewFunction("getFullYear",
			semantics.NewOverload("timestamp_to_year", true, types.Int64, types.Timestamp),
			semantics.NewOverload("timestamp_to_year_with_tz", true, types.Int64, types.Timestamp, types.String)))

	env.AddFunction(
		semantics.NewFunction("getMonth",
			semantics.NewOverload("timestamp_to_month", true, types.Int64, types.Timestamp),
			semantics.NewOverload("timestamp_to_month_with_tz", true, types.Int64, types.Timestamp, types.String)))

	env.AddFunction(
		semantics.NewFunction("getDayOfYear",
			semantics.NewOverload("timestamp_to_day_of_year", true, types.Int64, types.Timestamp),
			semantics.NewOverload("timestamp_to_day_of_year_with_tz", true, types.Int64, types.Timestamp, types.String)))

	env.AddFunction(
		semantics.NewFunction("getDayOfMonth",
			semantics.NewOverload("timestamp_to_day_of_month", true, types.Int64, types.Timestamp),
			semantics.NewOverload("timestamp_to_day_of_month_with_tz", true, types.Int64, types.Timestamp, types.String)))

	env.AddFunction(
		semantics.NewFunction("getDate",
			semantics.NewOverload("timestamp_to_day_of_month_1_based", true, types.Int64, types.Timestamp),
			semantics.NewOverload("timestamp_to_day_of_month_1_based_with_tz", true, types.Int64, types.Timestamp, types.String)))

	env.AddFunction(
		semantics.NewFunction("getDayOfWeek",
			semantics.NewOverload("timestamp_to_day_of_week", true, types.Int64, types.Timestamp),
			semantics.NewOverload("timestamp_to_day_of_week_with_tz", true, types.Int64, types.Timestamp, types.String)))

	env.AddFunction(
		semantics.NewFunction("getHours",
			semantics.NewOverload("timestamp_to_hours", true, types.Int64, types.Timestamp),
			semantics.NewOverload("timestamp_to_hours_with_tz", true, types.Int64, types.Timestamp, types.String),
			semantics.NewOverload("duration_to_hours", true, types.Int64, types.Duration)))

	env.AddFunction(
		semantics.NewFunction("getMinutes",
			semantics.NewOverload("timestamp_to_minutes", true, types.Int64, types.Timestamp),
			semantics.NewOverload("timestamp_to_minutes_with_tz", true, types.Int64, types.Timestamp, types.String),
			semantics.NewOverload("duration_to_minutes", true, types.Int64, types.Duration)))

	env.AddFunction(
		semantics.NewFunction("getSeconds",
			semantics.NewOverload("timestamp_to_seconds", true, types.Int64, types.Timestamp),
			semantics.NewOverload("timestamp_to_seconds_with_tz", true, types.Int64, types.Timestamp, types.String),
			semantics.NewOverload("duration_to_seconds", true, types.Int64, types.Duration)))

	env.AddFunction(
		semantics.NewFunction("getMilliseconds",
			semantics.NewOverload("timestamp_to_milliseconds", true, types.Int64, types.Timestamp),
			semantics.NewOverload("timestamp_to_milliseconds_with_tz", true, types.Int64, types.Timestamp, types.String),
			semantics.NewOverload("duration_to_milliseconds", true, types.Int64, types.Duration)))
}
