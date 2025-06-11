package mappath_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gekatateam/mappath"
)

func TestGet(t *testing.T) {
	tests := map[string]struct {
		p      any
		key    string
		result any
		err    any
	}{
		"from map, no slices, ok value": {
			p: map[string]any{
				"foo": "bar",
				"fizz": map[string]any{
					"buzz": 133,
				},
			},
			key:    "fizz.buzz",
			result: 133,
			err:    nil,
		},
		"from map, first level, ok value": {
			p: map[string]any{
				"foo": "bar",
				"fizz": map[string]any{
					"buzz": 133,
				},
			},
			key:    "foo",
			result: "bar",
			err:    nil,
		},
		"from map, last slice, ok value": {
			p: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz", "bazz",
				},
			},
			key:    "fizz.1",
			result: "bazz",
			err:    nil,
		},
		"from map, last slice, zero index, ok value": {
			p: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz", "bazz",
				},
			},
			key:    "fizz.0",
			result: "buzz",
			err:    nil,
		},
		"from map, through slice, ok value": {
			p: map[string]any{
				"foo": "bar",
				"fizz": []any{
					map[string]any{
						"buzz": 33,
					},
					map[string]any{
						"buzz": 33,
						"bazz": 44,
					},
				},
			},
			key:    "fizz.1.bazz",
			result: 44,
			err:    nil,
		},
		"from slice, through slice, ok value": {
			p: []any{
				map[string]any{
					"foo": "bar",
					"fizz": []any{
						map[string]any{
							"buzz": 33,
						},
						map[string]any{
							"buzz": 33,
							"bazz": 44,
						},
					},
				},
				"lorem",
			},
			key:    "0.fizz.1.bazz",
			result: 44,
			err:    nil,
		},
		"from slice, on root, ok value": {
			p: []any{
				map[string]any{
					"foo": "bar",
					"fizz": []any{
						map[string]any{
							"buzz": 33,
						},
						map[string]any{
							"buzz": 33,
							"bazz": 44,
						},
					},
				},
				"lorem",
			},
			key:    "1",
			result: "lorem",
			err:    nil,
		},
		"from slice, through slice, no such value": {
			p: []any{
				map[string]any{
					"foo": "bar",
					"fizz": []any{
						map[string]any{
							"buzz": 33,
						},
						map[string]any{
							"buzz": 33,
							"bazz": 44,
						},
					},
				},
				"lorem",
			},
			key:    "1.fizz.1.bazz",
			result: nil,
			err:    new(*mappath.NotFoundError),
		},
		"from slice, negative index, ok result": {
			p: []any{
				map[string]any{
					"foo": "bar",
					"fizz": []any{
						map[string]any{
							"buzz": 33,
						},
						map[string]any{
							"buzz": 33,
							"bazz": 44,
						},
					},
				},
				"lorem",
			},
			key:    "0.fizz.-1.bazz",
			result: 44,
			err:    nil,
		},
		"from slice, incorrect index, bad value": {
			p: []any{
				map[string]any{
					"foo": "bar",
					"fizz": []any{
						map[string]any{
							"buzz": 33,
						},
						map[string]any{
							"buzz": 33,
							"bazz": 44,
						},
					},
				},
				"lorem",
			},
			key:    "0.fizz.-100.bazz",
			result: nil,
			err:    new(*mappath.NotFoundError),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			val, err := mappath.Get(test.p, test.key)

			if err != nil {
				if !errors.As(err, &test.err) {
					t.Errorf("unexpected error - want: %v, got: %v", test.err, err)
				}
			} else {
				if test.err != nil {
					t.Errorf("unexpected error - want: nil, got: %v", err)
				}
			}

			if !reflect.DeepEqual(val, test.result) {
				t.Errorf("unexpected result - want: %v, got: %v", test.result, val)
			}
		})
	}
}

func TestPut(t *testing.T) {
	tests := map[string]struct {
		p      any
		key    string
		val    any
		result any
		err    any
	}{
		"add new key, in slice, ok result": {
			p: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"bizz",
				},
			},
			key: "fizz.3",
			val: 1337,
			result: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"bizz",
					nil,
					1337,
				},
			},
			err: nil,
		},
		"add new key, through slice, ok result": {
			p: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"bizz",
				},
			},
			key: "fizz.3",
			val: map[string]any{
				"leet": 1337,
			},
			result: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"bizz",
					nil,
					map[string]any{
						"leet": 1337,
					},
				},
			},
			err: nil,
		},
		"update current key, through slice, ok result": {
			p: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"bizz",
					nil,
					map[string]any{
						"leet": 1337,
					},
				},
			},
			key: "fizz.3.leet",
			val: "xxxx",
			result: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"bizz",
					nil,
					map[string]any{
						"leet": "xxxx",
					},
				},
			},
			err: nil,
		},
		"update current key, through slice, zero index, ok result": {
			p: map[string]any{
				"foo": "bar",
				"fizz": []any{
					map[string]any{
						"leet": 1337,
					},
					"bizz",
					nil,
					map[string]any{
						"leet": 1337,
					},
				},
			},
			key: "fizz.0.leet",
			val: "xxxx",
			result: map[string]any{
				"foo": "bar",
				"fizz": []any{
					map[string]any{
						"leet": "xxxx",
					},
					"bizz",
					nil,
					map[string]any{
						"leet": 1337,
					},
				},
			},
			err: nil,
		},
		"update current key, in slice, negative index, ok result": {
			p: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"bizz",
					nil,
					"leet",
				},
			},
			key: "fizz.-3",
			val: "xxxx",
			result: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"xxxx",
					nil,
					"leet",
				},
			},
			err: nil,
		},
		"replace current key, through slice, ok result": {
			p: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"bizz",
					nil,
					map[string]any{
						"leet": 1337,
					},
				},
			},
			key: "fizz.3",
			val: "xxxx",
			result: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"bizz",
					nil,
					"xxxx",
				},
			},
			err: nil,
		},
		"add new key, through slice, bad path": {
			p: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"bizz",
				},
			},
			key:    "fizz.buzz",
			val:    1337,
			result: nil,
			err:    new(*mappath.InvalidPathError),
		},
		"add new key, through map, bad path": {
			p: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"bizz",
				},
			},
			key:    "foo.0",
			val:    1337,
			result: nil,
			err:    new(*mappath.InvalidPathError),
		},
		"add new key, no input, ok path": {
			p:   nil,
			key: "0.fizz.3.buzz",
			val: 1337,
			result: []any{
				map[string]any{
					"fizz": []any{
						nil,
						nil,
						nil,
						map[string]any{
							"buzz": 1337,
						},
					},
				},
			},
			err: nil,
		},
		"add new key, invalid index, bad path": {
			p: []any{
				map[string]any{
					"fizz": []any{
						nil,
						nil,
						nil,
						map[string]any{
							"buzz": 1337,
						},
					},
				},
			},
			key:    "0.fizz.-300.buzz",
			val:    1337,
			result: nil,
			err:    new(*mappath.InvalidPathError),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			val, err := mappath.Put(test.p, test.key, test.val)

			if err != nil {
				if !errors.As(err, &test.err) {
					t.Errorf("unexpected error - want: %v, got: %v", test.err, err)
				}
			} else {
				if test.err != nil {
					t.Errorf("unexpected error - want: nil, got: %v", err)
				}
			}

			if !reflect.DeepEqual(val, test.result) {
				t.Errorf("unexpected result - want: %v, got: %v", test.result, val)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := map[string]struct {
		p      any
		key    string
		result any
		err    any
	}{
		"delete simple key, through slice, ok result": {
			p: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"bizz",
					nil,
					map[string]any{
						"leet": 1337,
					},
				},
			},
			key: "fizz.3.leet",
			result: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"bizz",
					nil,
					map[string]any{},
				},
			},
			err: nil,
		},
		"delete simple key on zero index, through slice, ok result": {
			p: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"bizz",
					nil,
					map[string]any{
						"leet": 1337,
					},
				},
			},
			key: "fizz.0",
			result: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"bizz",
					nil,
					map[string]any{
						"leet": 1337,
					},
				},
			},
			err: nil,
		},
		"delete complex key, in slice, ok result": {
			p: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"bizz",
					nil,
					map[string]any{
						"leet": 1337,
					},
				},
			},
			key: "fizz.3",
			result: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"bizz",
					nil,
				},
			},
			err: nil,
		},
		"delete simple key, in slice, ok path": {
			p: []any{
				map[string]any{
					"fizz": []any{
						nil,
						nil,
						nil,
						map[string]any{
							"buzz": 1337,
						},
					},
				},
			},
			key: "0.fizz.2",
			result: []any{
				map[string]any{
					"fizz": []any{
						nil,
						nil,
						map[string]any{
							"buzz": 1337,
						},
					},
				},
			},
			err: nil,
		},
		"delete key, no such key, bad result": {
			p: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"bizz",
					nil,
					map[string]any{
						"leet": 1337,
					},
				},
			},
			key:    "fizz.5.bazz",
			result: nil,
			err:    new(*mappath.NotFoundError),
		},
		"delete key, negative out-of-range key, bad result": {
			p: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"bizz",
					nil,
					map[string]any{
						"leet": 1337,
					},
				},
			},
			key:    "fizz.-5.bazz",
			result: nil,
			err:    new(*mappath.InvalidPathError),
		},
		"delete key, invalid index, bad result": {
			p: map[string]any{
				"foo": "bar",
				"fizz": []any{
					"buzz",
					"bizz",
					nil,
					map[string]any{
						"leet": 1337,
					},
				},
			},
			key:    "fizz.100.bazz",
			result: nil,
			err:    new(*mappath.NotFoundError),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			val, err := mappath.Delete(test.p, test.key)

			if err != nil {
				if !errors.As(err, &test.err) {
					t.Errorf("unexpected error - want: %v, got: %v", test.err, err)
				}
			} else {
				if test.err != nil {
					t.Errorf("unexpected error - want: nil, got: %v", err)
				}
			}

			if !reflect.DeepEqual(val, test.result) {
				t.Errorf("unexpected result - want: %v, got: %v", test.result, val)
			}
		})
	}
}
