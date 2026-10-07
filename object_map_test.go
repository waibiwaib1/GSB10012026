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
	"reflect"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

type MapUser struct {
	Name string
	Age  int32
}

func (MapUser) JavaClassName() string {
	return "com.test.map_user_registered"
}

func TestDecodeObjectToMap(t *testing.T) {
	// encode a map as an unregistered class
	m := map[string]interface{}{
		"name": "tom",
		"age":  int32(18),
	}
	e := NewEncoder()
	err := e.EncodeMapAsClass("com.test.map_user", m)
	assert.NoError(t, err)
	data1 := e.Buffer()

	// decode to map without registering the POJO
	d := NewDecoder(data1)
	res, err := d.Decode()
	assert.NoError(t, err)

	dm, ok := res.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "com.test.map_user", dm[ClassKey])
	assert.Equal(t, "tom", dm["name"])
	assert.Equal(t, int32(18), dm["age"])

	// encode the decoded map back to bytes, then decode again to compare
	e2 := NewEncoder()
	err = e2.EncodeMapClass(dm)
	assert.NoError(t, err)
	res2, err := NewDecoder(e2.Buffer()).Decode()
	assert.NoError(t, err)
	assert.Equal(t, dm, res2)

	// encode the map with the class info found in decoder,
	// the field order is kept so the result bytes should be equal
	clsDef := d.FindClassInfo("com.test.map_user")
	assert.NotNil(t, clsDef)
	e4 := NewEncoder()
	err = e4.EncodeMapAsObject(clsDef, dm)
	assert.NoError(t, err)
	assert.Equal(t, data1, e4.Buffer())
}

func TestEncodeMapAsRegisteredClass(t *testing.T) {
	RegisterPOJO(&MapUser{})

	user := &MapUser{Name: "tom", Age: 18}
	e := NewEncoder()
	err := e.Encode(user)
	assert.NoError(t, err)
	data1 := e.Buffer()

	// encode a map as the registered class, the result should be equal
	m := map[string]interface{}{
		"name": "tom",
		"age":  int32(18),
	}
	e2 := NewEncoder()
	err = e2.EncodeMapAsClass("com.test.map_user_registered", m)
	assert.NoError(t, err)
	assert.Equal(t, data1, e2.Buffer())
}

func TestDecodeObjectToMapStrict(t *testing.T) {
	m := map[string]interface{}{
		"name": "tom",
		"age":  int32(18),
	}
	e := NewEncoder()
	err := e.EncodeMapAsClass("com.test.strict_user", m)
	assert.NoError(t, err)

	// strict mode: error returned for unregistered class
	d := NewStrictDecoder(e.Buffer())
	_, err = d.Decode()
	assert.Error(t, err)

	// non-strict mode: decoded to map
	d2 := NewDecoder(e.Buffer())
	res, err := d2.Decode()
	assert.NoError(t, err)
	dm, ok := res.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "com.test.strict_user", dm[ClassKey])
}

func TestDecodeObjectRegisteredStillStruct(t *testing.T) {
	user := &MapUser{Name: "tom", Age: 18}
	e := NewEncoder()
	err := e.Encode(user)
	assert.NoError(t, err)

	d := NewDecoder(e.Buffer())
	res, err := d.Decode()
	assert.NoError(t, err)
	assert.True(t, reflect.TypeOf(res) == reflect.TypeOf(&MapUser{}))
}
