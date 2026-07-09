package registry

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
)

const modelsBucketName = "models"

// BboltStore persists model definitions in a bbolt database file.
type BboltStore struct {
	db       *bolt.DB
	tempPath string
}

// NewBboltStore opens or creates a bbolt database at path.
func NewBboltStore(path string) (*BboltStore, error) {
	db, err := bolt.Open(path, 0o600, &bolt.Options{Timeout: 500 * time.Millisecond})
	if err == nil {
		if err := initBucket(db); err != nil {
			_ = db.Close()
			return nil, err
		}
		return &BboltStore{db: db}, nil
	}

	if errors.Is(err, bolt.ErrTimeout) {
		// Try to create a temporary copy to allow read-only operations
		tempDir := os.TempDir()
		tempPath := filepath.Join(tempDir, fmt.Sprintf("tinybrain-models-%d.db", time.Now().UnixNano()))
		if copyErr := copyFile(path, tempPath); copyErr == nil {
			tempDB, openErr := bolt.Open(tempPath, 0o600, nil)
			if openErr == nil {
				if err := initBucket(tempDB); err != nil {
					_ = tempDB.Close()
					_ = os.Remove(tempPath)
					return nil, err
				}
				return &BboltStore{db: tempDB, tempPath: tempPath}, nil
			}
			_ = os.Remove(tempPath)
		}
	}

	return nil, fmt.Errorf("open bbolt store %q: %w", path, err)
}

func initBucket(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(modelsBucketName))
		return err
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// RegisterModel inserts a model definition. Duplicate IDs return ErrDuplicateID.
func (s *BboltStore) RegisterModel(def ModelDefinition) error {
	if s.tempPath != "" {
		return fmt.Errorf("cannot write to model registry: database is locked by another process (is the daemon running?)")
	}
	if def.ID == "" {
		return fmt.Errorf("model ID is required")
	}

	value, err := encodeModelDefinition(def)
	if err != nil {
		return fmt.Errorf("encode model %q: %w", def.ID, err)
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(modelsBucketName))
		if bucket == nil {
			return fmt.Errorf("models bucket missing")
		}
		if bucket.Get([]byte(def.ID)) != nil {
			return ErrDuplicateID
		}
		return bucket.Put([]byte(def.ID), value)
	})
}

// GetModel returns the model definition for id.
func (s *BboltStore) GetModel(id string) (ModelDefinition, error) {
	var def ModelDefinition

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(modelsBucketName))
		if bucket == nil {
			return fmt.Errorf("models bucket missing")
		}
		raw := bucket.Get([]byte(id))
		if raw == nil {
			return ErrNotFound
		}
		decoded, err := decodeModelDefinition(raw)
		if err != nil {
			return fmt.Errorf("decode model %q: %w", id, err)
		}
		def = decoded
		return nil
	})
	if err != nil {
		return ModelDefinition{}, err
	}
	return def, nil
}

// ListModels returns all registered model definitions.
func (s *BboltStore) ListModels() []ModelDefinition {
	var out []ModelDefinition

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(modelsBucketName))
		if bucket == nil {
			return nil
		}
		return bucket.ForEach(func(_, value []byte) error {
			def, err := decodeModelDefinition(value)
			if err != nil {
				return err
			}
			out = append(out, def)
			return nil
		})
	})
	if err != nil {
		return nil
	}

	return out
}

// IsEmpty reports whether the store contains no model definitions.
func (s *BboltStore) IsEmpty() (bool, error) {
	empty := true
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(modelsBucketName))
		if bucket == nil {
			return nil
		}
		empty = bucket.Stats().KeyN == 0
		return nil
	})
	return empty, err
}

// Close closes the underlying bbolt database.
func (s *BboltStore) Close() error {
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	if s.tempPath != "" {
		_ = os.Remove(s.tempPath)
	}
	return err
}

func encodeModelDefinition(def ModelDefinition) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(def); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decodeModelDefinition(raw []byte) (ModelDefinition, error) {
	var def ModelDefinition
	if err := gob.NewDecoder(bytes.NewReader(raw)).Decode(&def); err != nil {
		return ModelDefinition{}, err
	}
	return def, nil
}
