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
	"testing"

	"github.com/stretchr/testify/assert"
)

type genericUser struct {
	Age  uint32
	Name string
}

func (u *genericUser) JavaClassName() string {
	return "org.apache.com.myUser"
}

type genericOrder struct {
	Id    int64
	Items []string
	Tags  map[string]string
	User  *genericUser
}

func (o *genericOrder) JavaClassName() string {
	return "org.apache.com.myOrder"
}

func TestGenericEncodeMapDecodeMap(t *testing.T) {
	user := &genericUser{Name: "xxx", Age: 11}

	e := NewEncoder()
	err := e.Encode(user)
	assert.NoError(t, err)
	data1 := e.Buffer()

	// decode pojo into generic map without registering genericUser
	d := NewDecoder(data1)
	v, err := d.DecodeMap()
	assert.NoError(t, err)

	m, ok := v.(GenericMap)
	assert.True(t, ok)
	assert.Equal(t, "org.apache.com.myUser", m.GetJavaClassName())
	assert.Equal(t, "xxx", m["name"])
	assert.Equal(t, int64(11), m["age"])

	// restore the pojo bytes from the generic map
	e2 := NewEncoder()
	err = e2.EncodeMap(map[string]interface{}{
		"name": m["name"],
		"age":  m["age"],
	}, "org.apache.com.myUser")
	assert.NoError(t, err)
	assert.Equal(t, data1, e2.Buffer())
}

func TestGenericNestedPojo(t *testing.T) {
	order := &genericOrder{
		Id:    100,
		User:  &genericUser{Name: "yyy", Age: 22},
		Items: []string{"a", "b"},
		Tags:  map[string]string{"k1": "v1"},
	}

	e := NewEncoder()
	assert.NoError(t, e.Encode(order))

	d := NewDecoder(e.Buffer())
	v, err := d.DecodeMap()
	assert.NoError(t, err)

	m := v.(GenericMap)
	assert.Equal(t, "org.apache.com.myOrder", m.GetJavaClassName())
	assert.Equal(t, int64(100), m["id"])

	user := m["user"].(GenericMap)
	assert.Equal(t, "org.apache.com.myUser", user.GetJavaClassName())
	assert.Equal(t, "yyy", user["name"])
	assert.Equal(t, int64(22), user["age"])

	assert.Equal(t, []interface{}{"a", "b"}, m["items"])

	tags := m["tags"].(map[interface{}]interface{})
	assert.Equal(t, "v1", tags["k1"])
}

func TestGenericRoundTripBytes(t *testing.T) {
	userMap := GenericMap{
		"name": "zzz",
		"age":  int64(33),
	}
	userMap.SetJavaClassName("org.apache.com.myUser")

	e := NewEncoder()
	assert.NoError(t, e.Encode(userMap))
	data := e.Buffer()

	d := NewDecoder(data)
	v, err := d.DecodeMap()
	assert.NoError(t, err)

	e2 := NewEncoder()
	assert.NoError(t, e2.Encode(v))
	assert.Equal(t, data, e2.Buffer())
}
