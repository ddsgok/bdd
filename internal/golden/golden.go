package golden

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/ddspog/bdd/internal/shared"
)

var (
	// testdata stores test cases about current feature being tested.
	testdata = make(map[string][]*Gold)
	// currentFeature tells the current feature bein tested.
	currentFeature = ""
)

// Gold contains information about a test case on a golden file,
// separated in Input and Golden.
type Gold struct {
	Input  interface{}
	Golden interface{}
}

// Load unmarshall the json into input and gold pointers received.
func (g *Gold) Load(input, gold interface{}) {
	if jsonBytes, err := json.Marshal(g.Input); err == nil {
		json.Unmarshal(jsonBytes, input)
	}
	if jsonBytes, err := json.Marshal(g.Golden); err == nil {
		json.Unmarshal(jsonBytes, gold)
	}
}

// Update get an struct or a map, and loads into golden part of test
// case, to update file with new values.
func (g *Gold) Update(values interface{}) {
	if *update {
		if jsonBytes, err := json.Marshal(values); err == nil {
			json.Unmarshal(jsonBytes, &g.Golden)
		}
	}
}

// Manager load a golden file for a Feature, and then separates into
// various test cases.
type Manager struct {
	goldies []*Gold
	feature string
}

// Get returns the ith test case for the feature tested in manager.
func (m *Manager) Get(i int) (g shared.Golden) {
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
		if err := ensureDir(filepath.Dir(filename(m.feature))); err == nil {
			if jsonBytes, err := json.MarshalIndent(testdata, "", "    "); err == nil {
				ioutil.WriteFile(filename(m.feature), jsonBytes, FilePerms)
			}
		}
	}
}

// NewManager creates a manager, using the feature tested and given
// context.
func NewManager(feat, given string) (m *Manager) {
	feature := fmtFeature(feat)

	if currentFeature == feature {
		if _, ok := testdata[given]; ok {
			m = &Manager{goldies: testdata[given], feature: feature}
		}
	} else {
		if bytes, err := getBytes(feature); err == nil {
			if err = json.Unmarshal(bytes, &testdata); err == nil {
				if _, ok := testdata[given]; ok {
					m = &Manager{goldies: testdata[given], feature: feature}
				}
			}
		}
	}

	return
}