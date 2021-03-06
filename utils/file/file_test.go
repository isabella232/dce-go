/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/paypal/dce-go/config"
	"github.com/paypal/dce-go/types"
	"github.com/stretchr/testify/assert"
)

func TestPrefixTaskId(t *testing.T) {
	config.GetConfig().SetDefault(types.NO_FOLDER, true)
	prefix := "taskid"
	session := "session"
	res := PrefixTaskId(prefix, session)
	a := strings.Split(res, "_")
	if len(a) != 2 {
		t.Fatalf("expected len to be 2, got %v", len(a))
	}
	if a[0] != prefix {
		t.Fatalf("expected prefix to be 'taskid', got %s", a[0])
	}
	if a[1] != session {
		t.Fatalf("expected session to be 'session', got %s", a[1])
	}
}

func TestParseYamls(t *testing.T) {
	config.GetConfig().SetDefault(types.NO_FOLDER, true)
	yamls := []string{"testdata/docker-adhoc.yml", "testdata/docker-long.yml", "testdata/docker-empty.yml"}
	res, err := ParseYamls(&yamls)
	if err != nil {
		t.Fatalf("Got error to parseyamls %v", err)
	}
	//servs := res["testdata/docker-adhoc.yml"]["services"].(map[interface{}]interface{})
	fmt.Println(res)
}

func TestWriteToGeneratedFile(t *testing.T) {
	config.GetConfig().SetDefault(types.NO_FOLDER, true)
	file, err := WriteToFile("wirtetofiletest.txt", []byte("hello,world"))
	if err != nil {
		t.Fatalf("Got error to write to file %v", err)
	}
	if file != "wirtetofiletest.txt" {
		t.Fatalf("expected file name to be 'wirtetofiletest.txt', got %s", file)
	}
	exist := CheckFileExist(file)
	if !exist {
		t.Fatal("file isn't generated")
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		t.Errorf("Error to read file,%s", err.Error())
	}
	if string(data) != "hello,world" {
		t.Fatalf("expected content to be 'hello,world', got %s", string(data))
	}
	os.Remove(file)
}

func TestIndexArray(t *testing.T) {
	array := make([]string, 3)
	array[0] = "pen"
	array[1] = "apple"
	array[2] = "peach"
	i := IndexArray(array, "apple")
	assert.Equal(t, 1, i, "")
	i = IndexArray(array, "test")
	assert.Equal(t, -1, i, "")
}

func TestReplaceArrayElement(t *testing.T) {
	//Test array
	array := make([]interface{}, 3)
	array[0] = "pen"
	array[1] = "apple"
	array[2] = "peach"
	res := ReplaceElement(array, "pen", "pencil").([]interface{})
	if len(res) != len(array) || res[0] != "pencil" {
		t.Fatalf("expected first element to be 'pencil', but got %s", res[0])
	}

	//Replace element not exist in array
	res1 := ReplaceElement(array, "not_exist", "not_exist").([]interface{})
	if len(res1) != len(res) {
		t.Fatalf("Expected array doesn't change, but got %v \n", res1)
	}

	array[0] = "fruit=banana"
	res2 := ReplaceElement(array, "^fruit=", "fruit=apple").([]interface{})
	if res2[0] != "fruit=apple" {
		t.Fatalf("Expected fruit=apple replace fruit=banana, but got %v \n", res2)
	}

	array[0] = "fruit.banana"
	res5 := ReplaceElement(array, "^fruit$", "fruit=apple").([]interface{})
	if res5[0] != "fruit.banana" {
		t.Fatalf("Expected no changes, but got %v \n", res5)
	}

	array[0] = "fruitbanana"
	res6 := ReplaceElement(array, "^fruit$", "fruit=apple").([]interface{})
	if res6[0] != "fruitbanana" {
		t.Fatalf("Expected no changes, but got %v \n", res6)
	}

	//Test map
	m := make(map[interface{}]interface{})
	m["key1"] = "val1"
	m["key2"] = "val2"
	res3 := ReplaceElement(m, "key2", "val3").(map[interface{}]interface{})
	if res3["key2"].(string) != "val3" {
		t.Fatalf("expected first element to be 'val3', but got %s", res3["key3"])
	}

	res4 := ReplaceElement(m, "key3", "val3").(map[interface{}]interface{})
	_, ok := res4["key3"]
	if ok {
		t.Fatalf("Expected new element not added, but got %s", res4["key3"])
	}
}

func TestAppendElement(t *testing.T) {
	//Test array
	array := make([]interface{}, 3)
	array[0] = "pen"
	array[1] = "apple"
	array[2] = "peach"
	res := AppendElement(array, "monkey", "monkey").([]interface{})
	if len(res) != len(array)+1 || res[3] != "monkey" {
		t.Fatalf("expected first element to be 'monkey', but got %s", res[3])
	}

	//Test duplicate element won't be appended
	res = AppendElement(array, "monkey", "monkey").([]interface{})
	if len(res) != len(array)+1 || res[3] != "monkey" {
		t.Fatalf("expected first element to be 'monkey', but got %s", res[3])
	}

	array[0] = "fruit=banana"
	res2 := AppendElement(array, "^fruit=", "fruit=apple").([]interface{})
	if res2[0] != "fruit=apple" {
		t.Fatalf("Expected fruit=apple replace fruit=banana, but got %v \n", res2)
	}

	array[0] = "fruit.banana"
	res5 := AppendElement(array, "^fruit$", "fruit=apple1").([]interface{})
	if res5[len(res5)-1] != "fruit=apple1" {
		t.Fatalf("Expected fruit=apple will be appended, but got %v \n", res5)
	}

	//Test map
	m := make(map[interface{}]interface{})
	m["key1"] = "val1"
	m["key2"] = "val2"
	res1 := AppendElement(m, "key3", "val3").(map[interface{}]interface{})
	if res1["key3"].(string) != "val3" {
		t.Fatalf("expected first element to be 'val3', but got %s", res1["key3"])
	}
}

func TestSplitYAML(t *testing.T) {
	config.GetConfig().Set(types.NO_FOLDER, true)

	expectedFile1 := []string{"testdata/docker-compose-base.yml", "testdata/docker-compose-qa.yml", "testdata/docker-compose-production.yml", "testdata/docker-compose-sandbox.yml", "testdata/docker-compose-debug.yml", "testdata/docker-compose-healthcheck.yml"}
	files1, err := SplitYAML("testdata/yaml")
	if err != nil {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(files1, expectedFile1) {
		t.Errorf("expected files to be %v, but got %v", expectedFile1, files1)
	}

	for _, file := range files1 {
		os.Remove(file)
	}

	expectedFile2 := []string{"testdata/docker-adhoc.yml"}
	files2, err := SplitYAML("testdata/docker-adhoc.yml")
	if err != nil {
		t.Errorf("Error split file %s\n", err.Error())
	}
	if files2[0] != expectedFile2[0] {
		t.Errorf("expected files to be %s , but got %v", expectedFile2, files2)
	}

	for _, file := range files2 {
		os.Remove(file)
	}
}

func TestConvertArrayToMap(t *testing.T) {
	a := []interface{}{"a=b", "c", "d="}
	m := ConvertArrayToMap(a)
	assert.Equal(t, m["a"], "b", "a=b")
	assert.Equal(t, m["c"], "", "c")
	assert.Equal(t, m["d"], "", "d=")

	b := []interface{}{}
	m = ConvertArrayToMap(b)
	assert.Equal(t, len(m), 0, "empty map")
}

func TestGetDirFilesRecv(t *testing.T) {
	var files []string
	GetDirFilesRecv("testdata/docker-adhoc.yml", &files)
	assert.Equal(t, []string{"testdata/docker-adhoc.yml"}, files, "non archive format files should work")
	files = []string{}
	GetDirFilesRecv("testdata/config.zip", &files)
	assert.Equal(t, []string{"testdata/config/docker-adhoc.yml"}, files, "archive format files should work")
}
