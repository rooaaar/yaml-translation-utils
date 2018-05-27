package yamlutils

import (
	"errors"
	"strings"
)

// IdenticalKeys checks if both given interface have same keys and depth
func IdenticalKeys(ref interface{}, tra interface{}, parentPath string) error {
	if _, ok := ref.(string); ok {
		if _, ok := tra.(string); ok {
			return nil
		}
	}

	refMap, ok := ref.(map[interface{}]interface{})

	if ok {
		if traMap, ok := tra.(map[interface{}]interface{}); ok {
			var errs []string
			for key, val := range refMap {
				path := parentPath + "." + key.(string)
				if _, ok := traMap[key]; !ok {
					errs = append(errs, "cannot find '"+path+"' in translations")
				} else if err := IdenticalKeys(val, traMap[key], path); err != nil {
					errs = append(errs, err.Error())
				}
			}

			for key := range traMap {
				if _, ok := refMap[key]; !ok {
					path := parentPath + "." + key.(string)
					errs = append(errs, "found redunant translation '"+path+"'")
				}
			}

			if len(errs) > 0 {
				return errors.New(strings.Join(errs, "\n"))
			}

			return nil
		}

		return errors.New("cannot convert translation to map: " + parentPath)
	}

	return errors.New("cannot convert reference to map: " + parentPath)
}
