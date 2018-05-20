package golden

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type Gold struct {
	Input  interface{}
	Golden interface{}
}
type Manager struct {
	goldies []*Gold
	feature string
}

var (
	testdata = make(map[string][]*Gold)
)

type Golden interface {
	Load(interface{}, interface{})
	Update(interface{})
}

func (g *Gold) Load(input, gold interface{}) {
	if jsonBytes, err := json.Marshal(g.Input); err == nil {
		json.Unmarshal(jsonBytes, input)
	}
	if jsonBytes, err := json.Marshal(g.Golden); err == nil {
		json.Unmarshal(jsonBytes, gold)
	}
}
func (g *Gold) Update(values interface{}) {
	if *update {
		if jsonBytes, err := json.Marshal(values); err == nil {
			json.Unmarshal(jsonBytes, &g.Golden)
		}
	}
}
func NewManager(feat, given string) (m *Manager) {
	feature := fmtFeature(feat)
	if bytes, err := getBytes(feature); err == nil {
		m = &Manager{goldies: []*Gold{}, feature: feature}
		if err = json.Unmarshal(bytes, &testdata); err == nil {
			m.goldies = testdata[given]
		}
	}
	return
}
func (m *Manager) Get(i int) (g Golden) {
	g = m.goldies[i]
	return
}
func (m *Manager) NumGoldies() (n int) {
	n = len(m.goldies)
	return
}
func (m *Manager) Update() {
	if *update {
		if err := ensureDir(filepath.Dir(filename(m.feature))); err == nil {
			if jsonBytes, err := json.MarshalIndent(testdata, "", "    "); err == nil {
				ioutil.WriteFile(filename(m.feature), jsonBytes, FilePerms)
			}
		}
	}
}
