package main

import "testing"

func TestCheckKeys(t *testing.T) {
	interfaces := []struct {
		one interface{}
		two interface{}
	}{
		{
			map[interface{}]interface{}{
				"key":  "val1",
				"key2": "val2",
			},
			map[interface{}]interface{}{
				"key":  "val2",
				"key2": "val3",
			},
		},
		{
			map[interface{}]interface{}{
				"key2": "val2",
				"key":  "val1",
			},
			map[interface{}]interface{}{
				"key":  "val2",
				"key2": "val3",
			},
		},
		{
			map[interface{}]interface{}{
				"key": map[interface{}]interface{}{
					"sub key 1": "sub val 1",
				},
				"key2": "val1",
			},
			map[interface{}]interface{}{
				"key": map[interface{}]interface{}{
					"sub key 1": "sub val 1",
				},
				"key2": "val3",
			},
		},
		{
			map[interface{}]interface{}{
				"key": map[interface{}]interface{}{
					"sub key 1": "sub val 1",
					"sub key 2": "sub val 2",
				},
				"key2": "val1",
			},
			map[interface{}]interface{}{
				"key": map[interface{}]interface{}{
					"sub key 1": "sub val 1",
					"sub key 2": "sub val 2",
				},
				"key2": "val3",
			},
		},
		{
			map[interface{}]interface{}{
				"key2": "val1",
				"key": map[interface{}]interface{}{
					"sub key 1": "sub val 1",
					"sub key 2": "sub val 2",
				},
			},
			map[interface{}]interface{}{
				"key2": "val3",
				"key": map[interface{}]interface{}{
					"sub key 1": "sub val 1",
					"sub key 2": "sub val 2",
				},
			},
		},
		{
			map[interface{}]interface{}{
				"key": map[interface{}]interface{}{
					"sub key 1": "sub val 1",
					"sub key 2": map[interface{}]interface{}{
						"sub sub key 1": "sub sub val 1",
					},
				},
				"key2": "val1",
			},
			map[interface{}]interface{}{
				"key": map[interface{}]interface{}{
					"sub key 1": "sub val 1",
					"sub key 2": map[interface{}]interface{}{
						"sub sub key 1": "sub sub val 1",
					},
				},
				"key2": "val1",
			},
		},
		{
			map[interface{}]interface{}{
				"key": map[interface{}]interface{}{
					"sub key 1": "sub val 1",
					"sub key 2": map[interface{}]interface{}{
						"sub sub key 1": "sub sub val 1",
						"sub sub key 2": "sub sub val 1",
					},
				},
				"key2": "val",
			},
			map[interface{}]interface{}{
				"key": map[interface{}]interface{}{
					"sub key 1": "sub val 1",
					"sub key 2": map[interface{}]interface{}{
						"sub sub key 1": "sub sub val 1",
						"sub sub key 2": "sub sub val 1",
					},
				},
				"key2": "val1",
			},
		},
	}

	for _, in := range interfaces {
		err := checkKeys(in.one, in.two, "")
		if err != nil {
			t.Errorf("checkKeys(%v,%v,'') was incorrect, got: %s, want: nil.", in.one, in.two, err)
		}
	}
}

func TestCheckReduantsKeys(t *testing.T) {
	interfaces := []struct {
		one interface{}
		two interface{}
		err string
	}{
		{
			map[interface{}]interface{}{
				"key": "val",
			},
			map[interface{}]interface{}{
				"key":  "val",
				"key2": "val",
			},
			"find redunant translation '.key2'",
		},
		{
			map[interface{}]interface{}{
				"key":  "val",
				"key2": "val",
			},
			map[interface{}]interface{}{
				"key3": "val",
				"key":  "val",
				"key2": "val",
			},
			"find redunant translation '.key3'",
		},
		{
			map[interface{}]interface{}{
				"key2": "val1",
				"key": map[interface{}]interface{}{
					"sub key 1": "sub val 1",
					"sub key 2": "sub val 2",
				},
			},
			map[interface{}]interface{}{
				"key2": "val3",
				"key": map[interface{}]interface{}{
					"sub key 1": "sub val 1",
					"sub key 2": "sub val 2",
					"sub key 3": "sub val 2",
				},
			},
			"find redunant translation '.key.sub key 3'",
		},
		{
			map[interface{}]interface{}{
				"key": map[interface{}]interface{}{
					"sub key 1": "sub val 1",
					"sub key 2": map[interface{}]interface{}{
						"sub sub key 1": "sub sub val 1",
						"sub sub key 2": "sub sub val 1",
					},
				},
				"key2": "val",
			},
			map[interface{}]interface{}{
				"key": map[interface{}]interface{}{
					"sub key 1": "sub val 1",
					"sub key 2": map[interface{}]interface{}{
						"sub sub key 1": "sub sub val 1",
						"sub sub key 2": "sub sub val 1",
						"sub sub key 3": "sub sub val 1",
					},
				},
				"key2": "val1",
			},
			"find redunant translation '.key.sub key 2.sub sub key 3'",
		},
	}

	for _, in := range interfaces {
		err := checkKeys(in.one, in.two, "")
		if err == nil || err.Error() != in.err {
			t.Errorf("checkKeys(%v,%v,'') was incorrect, got: %v, want: %s.", in.one, in.two, err, in.err)
		}
	}
}

func TestCheckMissedKeys(t *testing.T) {
	interfaces := []struct {
		one interface{}
		two interface{}
		err string
	}{
		{
			map[interface{}]interface{}{
				"key":  "val",
				"key2": "val",
			},
			map[interface{}]interface{}{
				"key": "val",
			},
			"cannot find '.key2' in translations",
		},
		{
			map[interface{}]interface{}{
				"key":  "val",
				"key2": "val",
				"key3": map[interface{}]interface{}{
					"sub key 1": "sub val 1",
					"sub key 2": "sub val 2",
				},
			},
			map[interface{}]interface{}{
				"key":  "val",
				"key2": "val",
				"key3": map[interface{}]interface{}{
					"sub key 1": "sub val 1",
				},
			},
			"cannot find '.key3.sub key 2' in translations",
		},
		{
			map[interface{}]interface{}{
				"key": map[interface{}]interface{}{
					"sub key 1": "sub val 1",
					"sub key 2": map[interface{}]interface{}{
						"sub sub key 1": "sub sub val 1",
						"sub sub key 2": "sub sub val 1",
						"sub sub key 3": "sub sub val 1",
					},
				},
				"key2": "val",
			},
			map[interface{}]interface{}{
				"key": map[interface{}]interface{}{
					"sub key 1": "sub val 1",
					"sub key 2": map[interface{}]interface{}{
						"sub sub key 1": "sub sub val 1",
						"sub sub key 2": "sub sub val 1",
					},
				},
				"key2": "val1",
			},
			"cannot find '.key.sub key 2.sub sub key 3' in translations",
		},
	}

	for _, in := range interfaces {
		err := checkKeys(in.one, in.two, "")
		if err == nil || err.Error() != in.err {
			t.Errorf("checkKeys(%v,%v,'') was incorrect, got: %v, want: %s.", in.one, in.two, err, in.err)
		}
	}
}
