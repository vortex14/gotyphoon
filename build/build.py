# python build stubs for package build
# File is generated by gopy. Do not edit.
# gopy pkg build

from pybindgen import retval, param, Module
import sys

mod = Module('_build')
mod.add_include('"build_go.h"')
mod.add_function('GoPyInit', None, [])
mod.add_function('DecRef', None, [param('int64_t', 'handle')])
mod.add_function('IncRef', None, [param('int64_t', 'handle')])
mod.add_function('NumHandles', retval('int'), [])
mod.add_function('Slice_bool_CTor', retval('int64_t'), [])
mod.add_function('Slice_bool_len', retval('int'), [param('int64_t', 'handle')])
mod.add_function('Slice_bool_elem', retval('bool'), [param('int64_t', 'handle'), param('int', 'idx')])
mod.add_function('Slice_bool_subslice', retval('int64_t'), [param('int64_t', 'handle'), param('int', 'st'), param('int', 'ed')])
mod.add_function('Slice_bool_set', None, [param('int64_t', 'handle'), param('int', 'idx'), param('bool', 'value')])
mod.add_function('Slice_bool_append', None, [param('int64_t', 'handle'), param('bool', 'value')])
mod.add_function('Slice_byte_CTor', retval('int64_t'), [])
mod.add_function('Slice_byte_len', retval('int'), [param('int64_t', 'handle')])
mod.add_function('Slice_byte_elem', retval('uint8_t'), [param('int64_t', 'handle'), param('int', 'idx')])
mod.add_function('Slice_byte_subslice', retval('int64_t'), [param('int64_t', 'handle'), param('int', 'st'), param('int', 'ed')])
mod.add_function('Slice_byte_set', None, [param('int64_t', 'handle'), param('int', 'idx'), param('uint8_t', 'value')])
mod.add_function('Slice_byte_append', None, [param('int64_t', 'handle'), param('uint8_t', 'value')])
mod.add_function('Slice_float32_CTor', retval('int64_t'), [])
mod.add_function('Slice_float32_len', retval('int'), [param('int64_t', 'handle')])
mod.add_function('Slice_float32_elem', retval('float'), [param('int64_t', 'handle'), param('int', 'idx')])
mod.add_function('Slice_float32_subslice', retval('int64_t'), [param('int64_t', 'handle'), param('int', 'st'), param('int', 'ed')])
mod.add_function('Slice_float32_set', None, [param('int64_t', 'handle'), param('int', 'idx'), param('float', 'value')])
mod.add_function('Slice_float32_append', None, [param('int64_t', 'handle'), param('float', 'value')])
mod.add_function('Slice_float64_CTor', retval('int64_t'), [])
mod.add_function('Slice_float64_len', retval('int'), [param('int64_t', 'handle')])
mod.add_function('Slice_float64_elem', retval('double'), [param('int64_t', 'handle'), param('int', 'idx')])
mod.add_function('Slice_float64_subslice', retval('int64_t'), [param('int64_t', 'handle'), param('int', 'st'), param('int', 'ed')])
mod.add_function('Slice_float64_set', None, [param('int64_t', 'handle'), param('int', 'idx'), param('double', 'value')])
mod.add_function('Slice_float64_append', None, [param('int64_t', 'handle'), param('double', 'value')])
mod.add_function('Slice_int_CTor', retval('int64_t'), [])
mod.add_function('Slice_int_len', retval('int'), [param('int64_t', 'handle')])
mod.add_function('Slice_int_elem', retval('int64_t'), [param('int64_t', 'handle'), param('int', 'idx')])
mod.add_function('Slice_int_subslice', retval('int64_t'), [param('int64_t', 'handle'), param('int', 'st'), param('int', 'ed')])
mod.add_function('Slice_int_set', None, [param('int64_t', 'handle'), param('int', 'idx'), param('int64_t', 'value')])
mod.add_function('Slice_int_append', None, [param('int64_t', 'handle'), param('int64_t', 'value')])
mod.add_function('Slice_int16_CTor', retval('int64_t'), [])
mod.add_function('Slice_int16_len', retval('int'), [param('int64_t', 'handle')])
mod.add_function('Slice_int16_elem', retval('int16_t'), [param('int64_t', 'handle'), param('int', 'idx')])
mod.add_function('Slice_int16_subslice', retval('int64_t'), [param('int64_t', 'handle'), param('int', 'st'), param('int', 'ed')])
mod.add_function('Slice_int16_set', None, [param('int64_t', 'handle'), param('int', 'idx'), param('int16_t', 'value')])
mod.add_function('Slice_int16_append', None, [param('int64_t', 'handle'), param('int16_t', 'value')])
mod.add_function('Slice_int32_CTor', retval('int64_t'), [])
mod.add_function('Slice_int32_len', retval('int'), [param('int64_t', 'handle')])
mod.add_function('Slice_int32_elem', retval('int32_t'), [param('int64_t', 'handle'), param('int', 'idx')])
mod.add_function('Slice_int32_subslice', retval('int64_t'), [param('int64_t', 'handle'), param('int', 'st'), param('int', 'ed')])
mod.add_function('Slice_int32_set', None, [param('int64_t', 'handle'), param('int', 'idx'), param('int32_t', 'value')])
mod.add_function('Slice_int32_append', None, [param('int64_t', 'handle'), param('int32_t', 'value')])
mod.add_function('Slice_int64_CTor', retval('int64_t'), [])
mod.add_function('Slice_int64_len', retval('int'), [param('int64_t', 'handle')])
mod.add_function('Slice_int64_elem', retval('int64_t'), [param('int64_t', 'handle'), param('int', 'idx')])
mod.add_function('Slice_int64_subslice', retval('int64_t'), [param('int64_t', 'handle'), param('int', 'st'), param('int', 'ed')])
mod.add_function('Slice_int64_set', None, [param('int64_t', 'handle'), param('int', 'idx'), param('int64_t', 'value')])
mod.add_function('Slice_int64_append', None, [param('int64_t', 'handle'), param('int64_t', 'value')])
mod.add_function('Slice_int8_CTor', retval('int64_t'), [])
mod.add_function('Slice_int8_len', retval('int'), [param('int64_t', 'handle')])
mod.add_function('Slice_int8_elem', retval('int8_t'), [param('int64_t', 'handle'), param('int', 'idx')])
mod.add_function('Slice_int8_subslice', retval('int64_t'), [param('int64_t', 'handle'), param('int', 'st'), param('int', 'ed')])
mod.add_function('Slice_int8_set', None, [param('int64_t', 'handle'), param('int', 'idx'), param('int8_t', 'value')])
mod.add_function('Slice_int8_append', None, [param('int64_t', 'handle'), param('int8_t', 'value')])
mod.add_function('Slice_rune_CTor', retval('int64_t'), [])
mod.add_function('Slice_rune_len', retval('int'), [param('int64_t', 'handle')])
mod.add_function('Slice_rune_elem', retval('int32_t'), [param('int64_t', 'handle'), param('int', 'idx')])
mod.add_function('Slice_rune_subslice', retval('int64_t'), [param('int64_t', 'handle'), param('int', 'st'), param('int', 'ed')])
mod.add_function('Slice_rune_set', None, [param('int64_t', 'handle'), param('int', 'idx'), param('int32_t', 'value')])
mod.add_function('Slice_rune_append', None, [param('int64_t', 'handle'), param('int32_t', 'value')])
mod.add_function('Slice_string_CTor', retval('int64_t'), [])
mod.add_function('Slice_string_len', retval('int'), [param('int64_t', 'handle')])
mod.add_function('Slice_string_elem', retval('char*'), [param('int64_t', 'handle'), param('int', 'idx')])
mod.add_function('Slice_string_subslice', retval('int64_t'), [param('int64_t', 'handle'), param('int', 'st'), param('int', 'ed')])
mod.add_function('Slice_string_set', None, [param('int64_t', 'handle'), param('int', 'idx'), param('char*', 'value')])
mod.add_function('Slice_string_append', None, [param('int64_t', 'handle'), param('char*', 'value')])
mod.add_function('Slice_uint_CTor', retval('int64_t'), [])
mod.add_function('Slice_uint_len', retval('int'), [param('int64_t', 'handle')])
mod.add_function('Slice_uint_elem', retval('uint64_t'), [param('int64_t', 'handle'), param('int', 'idx')])
mod.add_function('Slice_uint_subslice', retval('int64_t'), [param('int64_t', 'handle'), param('int', 'st'), param('int', 'ed')])
mod.add_function('Slice_uint_set', None, [param('int64_t', 'handle'), param('int', 'idx'), param('uint64_t', 'value')])
mod.add_function('Slice_uint_append', None, [param('int64_t', 'handle'), param('uint64_t', 'value')])
mod.add_function('Slice_uint16_CTor', retval('int64_t'), [])
mod.add_function('Slice_uint16_len', retval('int'), [param('int64_t', 'handle')])
mod.add_function('Slice_uint16_elem', retval('uint16_t'), [param('int64_t', 'handle'), param('int', 'idx')])
mod.add_function('Slice_uint16_subslice', retval('int64_t'), [param('int64_t', 'handle'), param('int', 'st'), param('int', 'ed')])
mod.add_function('Slice_uint16_set', None, [param('int64_t', 'handle'), param('int', 'idx'), param('uint16_t', 'value')])
mod.add_function('Slice_uint16_append', None, [param('int64_t', 'handle'), param('uint16_t', 'value')])
mod.add_function('Slice_uint32_CTor', retval('int64_t'), [])
mod.add_function('Slice_uint32_len', retval('int'), [param('int64_t', 'handle')])
mod.add_function('Slice_uint32_elem', retval('uint32_t'), [param('int64_t', 'handle'), param('int', 'idx')])
mod.add_function('Slice_uint32_subslice', retval('int64_t'), [param('int64_t', 'handle'), param('int', 'st'), param('int', 'ed')])
mod.add_function('Slice_uint32_set', None, [param('int64_t', 'handle'), param('int', 'idx'), param('uint32_t', 'value')])
mod.add_function('Slice_uint32_append', None, [param('int64_t', 'handle'), param('uint32_t', 'value')])
mod.add_function('Slice_uint64_CTor', retval('int64_t'), [])
mod.add_function('Slice_uint64_len', retval('int'), [param('int64_t', 'handle')])
mod.add_function('Slice_uint64_elem', retval('uint64_t'), [param('int64_t', 'handle'), param('int', 'idx')])
mod.add_function('Slice_uint64_subslice', retval('int64_t'), [param('int64_t', 'handle'), param('int', 'st'), param('int', 'ed')])
mod.add_function('Slice_uint64_set', None, [param('int64_t', 'handle'), param('int', 'idx'), param('uint64_t', 'value')])
mod.add_function('Slice_uint64_append', None, [param('int64_t', 'handle'), param('uint64_t', 'value')])
mod.add_function('Slice_uint8_CTor', retval('int64_t'), [])
mod.add_function('Slice_uint8_len', retval('int'), [param('int64_t', 'handle')])
mod.add_function('Slice_uint8_elem', retval('uint8_t'), [param('int64_t', 'handle'), param('int', 'idx')])
mod.add_function('Slice_uint8_subslice', retval('int64_t'), [param('int64_t', 'handle'), param('int', 'st'), param('int', 'ed')])
mod.add_function('Slice_uint8_set', None, [param('int64_t', 'handle'), param('int', 'idx'), param('uint8_t', 'value')])
mod.add_function('Slice_uint8_append', None, [param('int64_t', 'handle'), param('uint8_t', 'value')])

mod.generate(open('build.c', 'w'))

