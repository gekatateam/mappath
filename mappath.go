package mappath

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type InvalidPathError struct {
	Path   string
	Reason string
}

func (e *InvalidPathError) Error() string { return fmt.Sprintf("%v: %v", e.Path, e.Reason) }

type NotFoundError struct {
	Path   string
	Reason string
}

func (e *NotFoundError) Error() string { return fmt.Sprintf("%v: %v", e.Path, e.Reason) }

// Get a value by specified key from provided map[string]any or []any.
func Get(p any, key string) (any, error) {
	if len(key) == 0 {
		return nil, &InvalidPathError{
			Path:   key,
			Reason: "key length cannot be zero",
		}
	}

	if key == "." {
		return p, nil
	}

	if key[0] == '.' {
		return nil, &InvalidPathError{
			Path:   key,
			Reason: "key cannot start from dot",
		}
	}

	for {
		dotIndex := strings.IndexRune(key, '.')
		if dotIndex < 0 { // no nested keys
			node, err := searchInNode(p, key)
			if err != nil {
				return nil, &NotFoundError{
					Path:   key,
					Reason: fmt.Sprintf("no such key: %v", err),
				}
			}

			return node, nil
		}

		next, err := searchInNode(p, key[:dotIndex])
		if err != nil {
			return nil, &NotFoundError{
				Path:   key,
				Reason: fmt.Sprintf("no such key: %v", err),
			}
		}

		key = key[dotIndex+1:] // shift path: foo.bar.buzz -> bar.buzz
		p = next               // shift searchable object
	}
}

// Put a passed value on a specified path in the provided map[string]any or []any and get the updated object.
func Put(p any, key string, val any) (any, error) {
	if len(key) == 0 {
		return nil, &InvalidPathError{
			Path:   key,
			Reason: "key length cannot be zero",
		}
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

		return nil, &InvalidPathError{
			Path:   key,
			Reason: "dot merge error: both root node and value must be map[string]any or []any",
		}
	}

	if key[0] == '.' {
		return nil, &InvalidPathError{
			Path:   key,
			Reason: "key cannot start from dot",
		}
	}

	return putInKey(p, key, val)
}

// Delete a value on a specified path in the provided map[string]any or []any and get the updated object.
func Delete(p any, key string) (any, error) {
	if len(key) == 0 {
		return nil, &InvalidPathError{
			Path:   key,
			Reason: "key length cannot be zero",
		}
	}

	if key == "." {
		return nil, nil
	}

	if key[0] == '.' {
		return nil, &InvalidPathError{
			Path:   key,
			Reason: "key cannot start from dot",
		}
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
	var invalidPathError *InvalidPathError
	if errors.As(err, &invalidPathError) {
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
			return nil, &NotFoundError{
				Path:   key,
				Reason: "no such key in map[string]any",
			}
		}
	case []any:
		i, err := strconv.Atoi(key)
		if err != nil {
			return nil, &NotFoundError{
				Path:   key,
				Reason: "target node is []any, but provided key cannot be converted into int",
			}
		}

		if (i >= 0) && (i < len(t)) {
			return t[i], nil
		}

		if i < 0 {
			if x := len(t) + i; (x >= 0) && (x < len(t)) {
				return t[x], nil
			} else {
				return nil, &NotFoundError{
					Path:   key,
					Reason: "node is a []any, but provided negative index is out of range",
				}
			}
		}

		return nil, &NotFoundError{
			Path:   key,
			Reason: "no such key in []any",
		}
	default:
		return nil, &InvalidPathError{
			Path:   key,
			Reason: "node must be a map[string]any or []any",
		}
	}
}

func createNode(key string) any {
	if i, err := strconv.Atoi(key); err == nil && i >= 0 {
		//lint:ignore S1019 explicitly indicates len and cap setting
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
		if err != nil {
			return nil, &InvalidPathError{
				Path:   key,
				Reason: "node is a []any, but provided key cannot be converted into int",
			}
		}

		if (i >= 0) && (i < len(t)) {
			t[i] = val
			return t, nil
		}

		if i < 0 {
			if x := len(t) + i; (x >= 0) && (x < len(t)) {
				t[x] = val
				return t, nil
			} else {
				return nil, &InvalidPathError{
					Path:   key,
					Reason: "node is a []any, but provided negative index is out of range",
				}
			}
		}

		n := slices.Grow(t, i+1-len(t))
		for j := len(n); j <= i; j++ {
			n = append(n, nil)
		}
		n[i] = val
		return n, nil
	default:
		return nil, &InvalidPathError{
			Path:   key,
			Reason: "node must be a map[string]any or []any",
		}
	}
}

func deleteFromNode(p any, key string) (any, error) {
	switch t := p.(type) {
	case map[string]any:
		if _, ok := t[key]; ok {
			delete(t, key)
			return t, nil
		}
		return nil, &NotFoundError{
			Path:   key,
			Reason: "no such key in map[string]any",
		}
	case []any:
		i, err := strconv.Atoi(key)
		if err != nil {
			return nil, &InvalidPathError{
				Path:   key,
				Reason: "node is a []any, but provided key cannot be converted into int",
			}
		}

		if (i >= 0) && (i < len(t)) {
			return slices.Delete(t, i, i+1), nil
		}

		if i < 0 {
			if x := len(t) + i; (x >= 0) && (x < len(t)) {
				return slices.Delete(t, x, x+1), nil
			} else {
				return nil, &NotFoundError{
					Path:   key,
					Reason: "node is a []any, but provided negative index is out of range",
				}
			}
		}

		return nil, &NotFoundError{
			Path:   key,
			Reason: "no such key in []any",
		}
	default:
		return nil, &InvalidPathError{
			Path:   key,
			Reason: "node must be a map[string]any or []any",
		}
	}
}
