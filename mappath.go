package mappath

import (
	"errors"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrInvalidPath  = errors.New("invalid path")
	ErrNoSuchField  = errors.New("no such field on the provided path")
)

// Get a value by specified key from provided map[string]any or []any.
func Get(p any, key string) (any, error) {
	if len(key) == 0 {
		return nil, ErrInvalidPath
	}

	if key == "." {
		return p, nil
	}

	if key[0] == '.' {
		return nil, ErrInvalidPath
	}

	for {
		dotIndex := strings.IndexRune(key, '.')
		if dotIndex < 0 { // no nested keys
			node, err := searchInNode(p, key)
			if err != nil {
				return nil, ErrNoSuchField
			}

			return node, nil
		}

		next, err := searchInNode(p, key[:dotIndex])
		if err != nil {
			return nil, ErrNoSuchField
		}

		key = key[dotIndex+1:] // shift path: foo.bar.buzz -> bar.buzz
		p = next // shift searchable object
	}
}

// Put a passed value on a specified path in the provided map[string]any or []any and get the updated object.
func Put(p any, key string, val any) (any, error) {
	if len(key) == 0 {
		return nil, ErrInvalidPath
	}

	if key == "." {
		if p == nil {
			return val, nil
		}

		if pl, ok := p.(map[string]any); ok {
			if vl, ok := val.(map[string]any); ok {
				for k, v := range vl {
					pl[k] = v
				}
				return pl, nil
			}
		}

		if pl, ok := p.([]any); ok {
			if vl, ok := val.([]any); ok {
				return append(pl, vl...), nil
			}
		}

		return nil, ErrInvalidPath
	}

	if key[0] == '.' {
		return nil, ErrInvalidPath
	}

	return putInKey(p, key, val)
}

// Delete a value on a specified path in the provided map[string]any or []any and get the updated object.
func Delete(p any, key string) (any, error) {
	if len(key) == 0 {
		return nil, ErrInvalidPath
	}

	if key == "." {
		return nil, nil
	}

	if key[0] == '.' {
		return nil, ErrInvalidPath
	}

	return deleteFromKey(p, key)
}

// Clone passed map[string]any or []any.
func Clone(p any) any {
	switch t := p.(type) {
	case map[string]any:
		m := make(map[string]any)
		for k := range t {
			m[k] = Clone(t[k])
		}
		return m
	case []any:
		s := make([]any, len(t))
		for i := range t {
			s[i] = Clone(t[i])
		}
		return s
	default:
		return p
	}
}

func putInKey(p any, key string, val any) (any, error) {
	dotIndex := strings.IndexRune(key, '.')
	if dotIndex < 0 { // no nested keys
		if p == nil {
			p = createNode(key)
		}
		return putInNode(p, key, val)
	}

	currNode := p
	currKey := key[:dotIndex]

	if currNode == nil {
		currNode = createNode(currKey)
	}

	nextKey := key[dotIndex+1:]
	nextNode, err := searchInNode(currNode, currKey)
	if err == ErrInvalidPath {
		return nil, err
	}

	nextNode, err = putInKey(nextNode, nextKey, val)
	if err != nil {
		return nil, err
	}

	return putInNode(currNode, currKey, nextNode)
}

func deleteFromKey(p any, key string) (any, error) {
	dotIndex := strings.IndexRune(key, '.')
	if dotIndex < 0 { // no nested keys
		if p == nil {
			p = createNode(key)
		}
		return deleteFromNode(p, key)
	}

	currNode := p
	currKey := key[:dotIndex]

	nextKey := key[dotIndex+1:]
	nextNode, err := searchInNode(currNode, currKey)
	if err != nil {
		return nil, err
	}

	nextNode, err = deleteFromKey(nextNode, nextKey)
	if err != nil {
		return nil, err
	}

	return putInNode(currNode, currKey, nextNode)
}

func searchInNode(p any, key string) (any, error) {
	switch t := p.(type) {
	case map[string]any:
		if val, ok := t[key]; ok {
			return val, nil
		} else {
			return nil, ErrNoSuchField
		}
	case []any:
		i, err := strconv.Atoi(key)
		if err != nil || i < 0 {
			return nil, ErrNoSuchField
		}

		if i < len(t) {
			return t[i], nil
		}

		return nil, ErrNoSuchField
	default:
		return nil, ErrInvalidPath
	}
}

func createNode(key string) (any) {
	if i, err := strconv.Atoi(key); err == nil && i >= 0 {
		s := make([]any, i+1, i+1)
		return s
	}

	m := make(map[string]any)
	return m
}

func putInNode(p any, key string, val any) (any, error) {
	switch t := p.(type) {
	case map[string]any:
		t[key] = val
		return t, nil
	case []any:
		i, err := strconv.Atoi(key)
		if err != nil || i < 0 {
			return nil, ErrInvalidPath
		}

		if i < len(t) {
			t[i] = val
			return t, nil
		}

		n := slices.Grow(t, i+1-len(t))
		for j := len(n); j <= i; j++ {
			n = append(n, nil)
		}
		n[i] = val
		return n, nil
	default:
		return nil, ErrInvalidPath
	}
}

func deleteFromNode(p any, key string) (any, error) {
	switch t := p.(type) {
	case map[string]any:
		if _, ok := t[key]; ok {
			delete(t, key)
			return t, nil
		}
		return nil, ErrNoSuchField
	case []any:
		i, err := strconv.Atoi(key)
		if err != nil || i < 0 {
			return nil, ErrInvalidPath
		}

		if i < len(t) {
			return append(t[:i], t[i+1:]...), nil
		}

		return nil, ErrNoSuchField
	default:
		return nil, ErrInvalidPath
	}
}
