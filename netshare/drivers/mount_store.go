package drivers

import (
	log "github.com/Sirupsen/logrus"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type mountStore struct {
	metaDataPath string
}

func NewVolumeStore(metaDataPath string) *mountStore {
	return &mountStore{
		metaDataPath: metaDataPath,
	}
}

func (m *mountStore) GetMounts() map[string]map[string]string {
	mounts := map[string]map[string]string{}

	if _, err := os.Stat(m.metaDataPath); err != nil {
		log.Debugf("Directory '%s' not found... creating", m.metaDataPath)
		os.Mkdir(m.metaDataPath, 0755)
		return mounts
	}

	log.Debugf("Reading metadata from: %s", m.metaDataPath)
	filepath.Walk(m.metaDataPath, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		log.Debugf("Reading metadata file from: %s", path)
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Debugf("Failed to read metadata file from: %s", path)
			return nil
		}

		var opts map[string]string
		if err := json.Unmarshal(content, &opts); err == nil {
			name,_ := filepath.Rel(m.metaDataPath, path)
			log.Debugf("Mount '%s' found with options: %v", name, opts)
			mounts[name] = opts
		}

		return nil
	})

	return mounts
}

func (m *mountStore) Add(name string, opts map[string]string) {
	subpath, _ := filepath.Split(name)
	path := filepath.Join(m.metaDataPath, subpath)

	log.Debugf("Metadata directory: %s", path)
	os.MkdirAll(path, 0755)

	filePath := filepath.Join(path, filepath.Base(name))
	log.Debugf("Metadata file path: %s", filePath)

	data, _ := json.Marshal(opts);
	if err := ioutil.WriteFile(filePath, data, 0760); err != nil {
		panic(err)
	}
}

func (m *mountStore) Remove(name string) {
	path := filepath.Join(m.metaDataPath, name)
	log.Debugf("Removing metadata: %s", path)

	os.Remove(path)
}
