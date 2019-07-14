package secretSafe

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	cipher "secretSafe/chiper"
	"sync"
)

type Safe struct {
	encodingKey string
	filePath    string
	mutex       sync.Mutex
	keyValues   map[string]string
}

func (s *Safe) load() error {
	f, err := os.Open(s.filePath)
	// In case the file is empty
	if err != nil {
		s.keyValues = make(map[string]string)
		return nil
	}
	defer f.Close()

	r, err := cipher.DecryptReader(s.encodingKey, f)
	if err != nil {
		return err
	}
	return s.readKeyValues(r)
}

func (s *Safe) save() error {
	f, err := os.OpenFile(s.filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	w, err := cipher.EncryptWriter(s.encodingKey, f)
	if err != nil {
		return err
	}
	return s.writeKeyValues(w)
}

func (s *Safe) readKeyValues(r io.Reader) error {
	decoder := json.NewDecoder(r) // -->decryptReader --> {bufferedReader} --> file
	return decoder.Decode(&s.keyValues)
}

func (s *Safe) writeKeyValues(w io.Writer) error {
	encoder := json.NewEncoder(w) // --> encrypt writer --> file
	return encoder.Encode(s.keyValues)
}

// Get returns the secret from our super secret Vault
func (s *Safe) Get(key string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	err := s.load()
	if err != nil {
		return "", err
	}
	value, ok := s.keyValues[key]
	if !ok {
		return "", errors.New("Secret: no value for your key")
	}
	return value, nil
}

// Set puts new entry to our super secret Vault
func (s *Safe) Set(key, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	err := s.load()
	if err != nil {
		return err
	}
	s.keyValues[key] = value
	err = s.save()
	if err != nil {
		return err
	}
	return nil
}

// ReturnAll returns a map of secrets
func (s *Safe) ReturnAll() (map[string]string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	err := s.load()
	if err != nil {
		return nil, err
	}
	if s.keyValues == nil {
		return nil, errors.New("No values in the super secret vault!")
	}
	return s.keyValues, nil
}

// Delete removes a secret from our super secret vault
func (s *Safe) Delete(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	err := s.load()

	if err != nil {
		return err
	}

	_, ok := s.keyValues[key]
	if ok {
		delete(s.keyValues, key)
		err = s.save()
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("The is no such key in our super secret Vault! Run 'list' to check the existant keys")

}

// File takes encoding key and file path and returns pointer to Safe structure
func File(encodingKey, filepath string) *Safe {
	return &Safe{encodingKey: encodingKey,
		filePath: filepath,
	}
}
