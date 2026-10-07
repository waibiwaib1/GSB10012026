/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hessian

import (
	"io"
	"reflect"
	"sort"
)

import (
	perrors "github.com/pkg/errors"
)

// genericClassKey carries the java class name of a decoded/encoded pojo.
// It starts with '#' so that it can never collide with a normal hessian field name.
const genericClassKey = "#class"

// GenericMap is the generic representation of a hessian pojo.
// It is a map whose keys are the field names of the java object and whose
// values are the decoded field values.
//
// The java class name of the pojo is stored under the special key
// genericClassKey ("#class"), which makes it possible to restore the exact
// hessian bytes without registering any go struct.
type GenericMap map[string]interface{}

// JavaClassName implements the POJO interface.
func (m GenericMap) JavaClassName() string {
	if name, ok := m[genericClassKey].(string); ok {
		return name
	}
	return ""
}

// SetJavaClassName sets the java class name of the generic map.
func (m GenericMap) SetJavaClassName(name string) {
	m[genericClassKey] = name
}

// GetJavaClassName returns the java class name of the generic map.
func (m GenericMap) GetJavaClassName() string {
	return m.JavaClassName()
}

/////////////////////////////////////////
// encode
/////////////////////////////////////////

// EncodeMap encodes m as a hessian pojo whose java class name is javaClassName.
// The map keys are used as field names and sorted lexicographically, so the
// produced bytes are deterministic.
//
// Nested pojos may be passed as GenericMap (carrying their own class names);
// plain maps are encoded as hessian untyped maps.
func (e *Encoder) EncodeMap(m map[string]interface{}, javaClassName string) error {
	gm := make(GenericMap, len(m)+1)
	for k, v := range m {
		gm[k] = v
	}
	gm[genericClassKey] = javaClassName
	return e.Encode(gm)
}

func (e *Encoder) encGenericMap(m GenericMap) error {
	javaClassName := m.JavaClassName()
	if javaClassName == "" {
		return perrors.New("GenericMap must carry a non-empty java class name")
	}

	// check ref
	if n, ok := e.checkRefMap(reflect.ValueOf(m)); ok {
		e.buffer = encRef(e.buffer, n)
		return nil
	}

	fieldNames := make([]string, 0, len(m))
	for k := range m {
		if k == genericClassKey {
			continue
		}
		fieldNames = append(fieldNames, k)
	}
	sort.Strings(fieldNames)

	// write object definition
	idx := -1
	for i := range e.classInfoList {
		if javaClassName == e.classInfoList[i].javaName {
			idx = i
			break
		}
	}
	if idx == -1 {
		var bDef []byte
		bDef = encByte(bDef, BC_OBJECT_DEF)
		bDef = encString(bDef, javaClassName)
		bDef = encInt32(bDef, int32(len(fieldNames)))
		for _, name := range fieldNames {
			bDef = encString(bDef, name)
		}

		clsDef := &classInfo{javaName: javaClassName, fieldNameList: fieldNames, buffer: bDef}
		idx = len(e.classInfoList)
		e.classInfoList = append(e.classInfoList, clsDef)
		e.buffer = append(e.buffer, bDef...)
	}

	// write object instance
	if byte(idx) <= OBJECT_DIRECT_MAX {
		e.buffer = encByte(e.buffer, byte(idx)+BC_OBJECT_DIRECT)
	} else {
		e.buffer = encByte(e.buffer, BC_OBJECT)
		e.buffer = encInt32(e.buffer, int32(idx))
	}

	for _, name := range fieldNames {
		if err := e.Encode(m[name]); err != nil {
			return perrors.Wrapf(err, "failed to encode generic field: %s", name)
		}
	}

	return nil
}

/////////////////////////////////////////
// decode
/////////////////////////////////////////

// DecodeMap decodes a hessian pojo into a GenericMap without needing the go
// struct to be registered beforehand. Every nested pojo is decoded into a
// GenericMap as well. The java class name is stored under the "#class" key.
func (d *Decoder) DecodeMap() (interface{}, error) {
	return d.decodeGeneric(TAG_READ)
}

func (d *Decoder) decodeGeneric(flag int32) (interface{}, error) {
	var (
		tag byte
		err error
	)

	if flag != TAG_READ {
		tag = byte(flag)
	} else {
		tag, err = d.ReadByte()
		if err != nil {
			return nil, perrors.WithStack(err)
		}
	}

	switch {
	case tag == BC_NULL:
		return nil, nil
	case tag == BC_REF:
		return d.decGenericRef(TAG_READ)
	case tag == BC_TRUE:
		return true, nil
	case tag == BC_FALSE:
		return false, nil

	case (0x80 <= tag && tag <= 0xbf) || (0xc0 <= tag && tag <= 0xcf) ||
		(0xd0 <= tag && tag <= 0xd7) || tag == BC_INT:
		return d.decInt32(int32(tag))

	case (tag >= 0xd8 && tag <= 0xef) || (tag >= 0xf0 && tag <= 0xff) ||
		(tag >= 0x38 && tag <= 0x3f) || (tag == BC_LONG_INT) || (tag == BC_LONG):
		return d.decInt64(int32(tag))

	case tag == BC_DATE_MINUTE || tag == BC_DATE:
		return d.decDate(int32(tag))

	case tag == BC_DOUBLE_ZERO || tag == BC_DOUBLE_ONE || tag == BC_DOUBLE_BYTE ||
		tag == BC_DOUBLE_SHORT || tag == BC_DOUBLE_MILL || tag == BC_DOUBLE:
		return d.decDouble(int32(tag))

	case tag == BC_STRING_CHUNK || tag == BC_STRING ||
		(tag >= BC_STRING_DIRECT && tag <= STRING_DIRECT_MAX) ||
		(tag >= 0x30 && tag <= 0x33):
		return d.decString(int32(tag))

	case tag == BC_BINARY || tag == BC_BINARY_CHUNK || (tag >= 0x20 && tag <= 0x2f) ||
		(tag >= BC_BINARY_SHORT && tag <= 0x3f):
		return d.decBinary(int32(tag))

	case (tag >= BC_LIST_DIRECT && tag <= 0x77) || tag == BC_LIST_FIXED || tag == BC_LIST_VARIABLE ||
		(tag >= BC_LIST_DIRECT_UNTYPED && tag <= 0x7f) ||
		tag == BC_LIST_FIXED_UNTYPED || tag == BC_LIST_VARIABLE_UNTYPED:
		return d.decGenericList(int32(tag))

	case tag == BC_MAP || tag == BC_MAP_UNTYPED:
		return d.decGenericMap(int32(tag))

	case tag == BC_OBJECT_DEF:
		clsDef, decErr := d.decClassDef()
		if decErr != nil {
			return nil, perrors.Wrap(decErr, "decodeGeneric->decClassDef")
		}
		d.appendClsDef(clsDef.(*classInfo))
		return d.decodeGeneric(TAG_READ)

	case tag == BC_OBJECT:
		idx, decErr := d.decInt32(TAG_READ)
		if decErr != nil {
			return nil, perrors.WithStack(decErr)
		}
		return d.decGenericInstance(int(idx))

	case BC_OBJECT_DIRECT <= tag && tag <= (BC_OBJECT_DIRECT+OBJECT_DIRECT_MAX):
		return d.decGenericInstance(int(tag - BC_OBJECT_DIRECT))

	default:
		return nil, perrors.Errorf("decodeGeneric illegal tag: %+v", tag)
	}
}

func (d *Decoder) decGenericRef(flag int32) (interface{}, error) {
	idx, err := d.decInt32(flag)
	if err != nil {
		return nil, perrors.WithStack(err)
	}
	if int(idx) < 0 || int(idx) >= len(d.refs) {
		return nil, ErrIllegalRefIndex
	}
	return d.refs[idx], nil
}

func (d *Decoder) decGenericInstance(idx int) (interface{}, error) {
	if idx < 0 || idx >= len(d.classInfoList) {
		return nil, perrors.Errorf("illegal class index @idx %d", idx)
	}
	cls := d.classInfoList[idx]

	m := make(GenericMap, len(cls.fieldNameList)+1)
	m[genericClassKey] = cls.javaName
	d.appendRefs(m)

	for _, fieldName := range cls.fieldNameList {
		v, err := d.decodeGeneric(TAG_READ)
		if err != nil {
			return nil, perrors.Wrapf(err, "decGenericInstance field name:%s", fieldName)
		}
		m[fieldName] = v
	}

	return m, nil
}

func (d *Decoder) decGenericMap(flag int32) (interface{}, error) {
	var (
		tag byte
		err error
		k   interface{}
		v   interface{}
	)

	if flag != TAG_READ {
		tag = byte(flag)
	} else {
		tag, _ = d.ReadByte()
	}

	switch {
	case tag == BC_NULL:
		return nil, nil
	case tag == BC_REF:
		return d.decGenericRef(TAG_READ)
	case tag == BC_MAP:
		// read map type, ignored: maps stay generic.
		if _, err = d.decMapType(); err != nil {
			return nil, err
		}
	case tag == BC_MAP_UNTYPED:
		// do nothing
	default:
		return nil, perrors.Errorf("illegal generic map tag:%+v", tag)
	}

	m := make(map[interface{}]interface{})
	d.appendRefs(m)
	for d.peekByte() != BC_END {
		if k, err = d.decodeGeneric(TAG_READ); err != nil {
			return nil, err
		}
		if v, err = d.decodeGeneric(TAG_READ); err != nil {
			return nil, err
		}
		m[k] = v
	}
	if _, err = d.ReadByte(); err != nil {
		return nil, perrors.WithStack(err)
	}
	return m, nil
}

func (d *Decoder) decGenericList(flag int32) (interface{}, error) {
	var (
		tag    byte
		err    error
		list   []interface{}
		length int32
	)

	if flag != TAG_READ {
		tag = byte(flag)
	} else {
		tag, _ = d.ReadByte()
	}

	switch {
	case tag == BC_NULL:
		return nil, nil
	case tag == BC_REF:
		return d.decGenericRef(TAG_READ)
	}

	// read list type / length
	switch {
	case tag >= BC_LIST_DIRECT && tag <= 0x77:
		length = int32(tag - BC_LIST_DIRECT)
		if _, err = d.decMapType(); err != nil {
			return nil, err
		}
	case tag == BC_LIST_FIXED:
		if _, err = d.decMapType(); err != nil {
			return nil, err
		}
		if length, err = d.decInt32(TAG_READ); err != nil {
			return nil, err
		}
	case tag == BC_LIST_VARIABLE:
		if _, err = d.decMapType(); err != nil {
			return nil, err
		}
		length = -1
	case tag >= BC_LIST_DIRECT_UNTYPED && tag <= 0x7f:
		length = int32(tag - BC_LIST_DIRECT_UNTYPED)
	case tag == BC_LIST_FIXED_UNTYPED:
		if length, err = d.decInt32(TAG_READ); err != nil {
			return nil, err
		}
	case tag == BC_LIST_VARIABLE_UNTYPED:
		length = -1
	default:
		return nil, perrors.Errorf("illegal generic list tag:%+v", tag)
	}

	if length >= 0 {
		list = make([]interface{}, 0, length)
	} else {
		list = make([]interface{}, 0, 8)
	}
	holder := d.appendRefs(list)

	for {
		if length >= 0 && int32(len(list)) >= length {
			break
		}
		if length < 0 && d.peekByte() == BC_END {
			if _, err = d.ReadByte(); err != nil {
				return nil, perrors.WithStack(err)
			}
			break
		}

		v, decErr := d.decodeGeneric(TAG_READ)
		if decErr != nil {
			if perrors.Is(decErr, io.EOF) {
				break
			}
			return nil, perrors.WithStack(decErr)
		}
		list = append(list, v)
		holder.change(reflect.ValueOf(list))
	}

	return list, nil
}
