package golden

import (
	"strings"

	"github.com/ddspog/bdd/internal/common"
	"github.com/pkg/errors"
)

var (
	// testdata stores test cases about current feature being tested.
	testdata = make(map[string][]*Gold)
	// currentFeature tells the current feature being tested.
	currentFeature = ""
	// encoder do operations on the golden file, whatever its type.
	encoder fileEncoder
	// ErrInvalidKeyPrefix it's an error returned when gold.Get is called with key starting with wrong format.
	ErrInvalidKeyPrefix = errors.New("the golden key must be prefixed by 'input.' or 'golden.'")
)

// Gold contains information about a test case on a golden file,
// separated in Input and Golden.
type Gold struct {
	Input  interface{} `json:"input" yaml:"input"`
	Golden interface{} `json:"golden" yaml:"golden"`
}

// Get returns value from golden file, using a json sequence of keys.
func (g *Gold) Get(key string) (val interface{}) {
	var err error
	if strings.HasPrefix(key, "input.") {
		val, err = encoder.Val(g.Input, strings.TrimPrefix(key, "input."))
	} else if strings.HasPrefix(key, "golden.") {
		val, err = encoder.Val(g.Golden, strings.TrimPrefix(key, "golden."))
	} else {
		err = ErrInvalidKeyPrefix
	}

	if err != nil {
		panic(err)
	}

	return
}

// Load unmarshall the json into input and gold pointers received.
func (g *Gold) Load(input, gold interface{}) {
	_ = encoder.Load(g.Input, input)
	_ = encoder.Load(g.Golden, gold)
}

// Update get an struct or a map, and loads into golden part of test
// case, to update file with new values.
func (g *Gold) Update(values func() interface{}) {
	if *update {
		_ = encoder.Load(values(), &g.Golden)
	}
}

// Manager load a golden file for a Feature, and then separates into
// various test cases.
type Manager struct {
	goldies []*Gold
	feature string
}

// Get returns the i-th test case for the feature tested in manager.
func (m *Manager) Get(i int) (g common.Golden) {
	g = m.goldies[i]
	return
}

// NumGoldies return number of test cases loaded with manager to a
// single feature.
func (m *Manager) NumGoldies() (n int) {
	n = len(m.goldies)
	return
}

// Update uses the new values received from each test case, and then
// write into golden file for the feature tested.
func (m *Manager) Update() {
	if *update {
		if err := encoder.Write(testdata); err != nil {
			panic(err)
		}
	}
}

// NewManager creates a manager, using the feature tested and given
// context.
func NewManager(feat, given string) (m *Manager) {
	feature := strings.Replace(
		strings.Title(feat),
		" ", "", -1,
	)

	var err error
	if encoder, err = newEncoder(feature); err == nil {
		panic(err)
	}

	if currentFeature != feature {
		if err = encoder.Read(&testdata); err != nil {
			panic(err)
		}
	}

	if _, ok := testdata[given]; ok {
		m = &Manager{goldies: testdata[given], feature: feature}
	}

	return
}
