package golden

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	// jsonType represents files with JSON extension.
	jsonType = fileType(".json")
	// yamlType represents files with YAML extension.
	yamlType = fileType(".yml")
)

var (
	// ErrWrongKeySequence is thrown when user tries to walk through
	// invalid sequence of maps and keys.
	ErrWrongKeySequence = errors.New("invalid key sequence for golden file")
)

// fileType tells the extension for a file.
type fileType string

// helperEncoder declare some functions used by each encoder.
type helperEncoder struct{}

// get transform input into a map, and then tries to access value using
// a key. If input isn't a map, throws an error.
func (he *helperEncoder) get(in interface{}, key string) (val interface{}, err error) {
	switch m := in.(type) {
	case map[string]interface{}:
		val = m[key]
	case map[interface{}]interface{}:
		val = m[key]
	default:
		err = ErrWrongKeySequence
	}

	return
}

// Val returns value from a map, using a sequence of keys separated by
// '.' and transform each inner value founded into a new map.
func (he *helperEncoder) Val(in interface{}, seqKeys string) (val interface{}, err error) {
	var keys []string
	val, keys = in, strings.Split(seqKeys, ".")

	for _, key := range keys {
		if err == nil {
			val, err = he.get(val, key)
		}
	}

	return
}

// jsonEncoder is responsible or parsing information inside JSON files.
type jsonEncoder struct {
	path fileHandler
	*helperEncoder
}

// Read file and store its value on destiny.
func (je *jsonEncoder) Read(dest interface{}) (err error) {
	var bytes []byte
	if bytes, err = je.path.Bytes(); err == nil {
		err = json.Unmarshal(bytes, dest)
	}

	return
}

// Write on file the contents of source.
func (je *jsonEncoder) Write(src interface{}) (err error) {
	if err = je.path.EnsureDir(); err == nil {
		var bytes []byte
		if bytes, err = json.MarshalIndent(src, "", "    "); err == nil {
			err = je.path.Write(bytes)
		}
	}

	return
}

// Load values from source to destiny.
func (je *jsonEncoder) Load(src, dest interface{}) (err error) {
	var bytes []byte
	if bytes, err = json.Marshal(src); err == nil {
		err = json.Unmarshal(bytes, dest)
	}

	return
}

// yamlEncoder is responsible or parsing information inside YAML files.
type yamlEncoder struct {
	path fileHandler
	*helperEncoder
}

// Read file and store its value on destiny.
func (ye *yamlEncoder) Read(dest interface{}) (err error) {
	var bytes []byte
	if bytes, err = ye.path.Bytes(); err == nil {
		err = yaml.Unmarshal(bytes, dest)
	}

	return
}

// Write on file the contents of source.
func (ye *yamlEncoder) Write(src interface{}) (err error) {
	if err = ye.path.EnsureDir(); err == nil {
		var bytes []byte
		if bytes, err = yaml.Marshal(src); err == nil {
			err = ye.path.Write(bytes)
		}
	}

	return
}

// Load values from source to destiny.
func (ye *yamlEncoder) Load(src, dest interface{}) (err error) {
	var bytes []byte
	if bytes, err = yaml.Marshal(src); err == nil {
		err = yaml.Unmarshal(bytes, dest)
	}

	return
}

// fileEncoder contains functions to operate on files, to encode
// and decode them to structs, maps and other values.
type fileEncoder interface {
	Read(interface{}) error
	Write(interface{}) error
	Load(interface{}, interface{}) error
	Val(interface{}, string) (interface{}, error)
}

// newEncoder creates a fileEncoder initializing with the path to the
// file the encoder will work. It selects the appropriate encoder type,
// depending of the file extension.
func newEncoder(name string) (fe fileEncoder) {
	if p, err := path(name); err == nil {
		switch p.ExtType() {
		case jsonType:
			fe = &jsonEncoder{p, &helperEncoder{}}
		case yamlType:
			fe = &yamlEncoder{p, &helperEncoder{}}
		}
	}

	return
}
