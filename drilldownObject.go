package main

import "strconv"

func drilldownObject(object interface{}, element string) interface{} {
	if object != nil {
		switch object.(type) {
		default:
			object = nil
		case []interface{}:
			parsed, parseErr := strconv.ParseInt(element, 10, 32)
			if parseErr == nil && len(object.([]interface{})) > int(parsed) {
				if val := object.([]interface{})[parsed]; val != nil {
					object = val
				} else {
					object = nil
				}
			} else {
				object = nil
			}
		case map[interface{}]interface{}:
			parsed, parseErr := strconv.ParseInt(element, 10, 32)
			if val, firstOk := object.(map[interface{}]interface{})[element]; firstOk {
				object = val
			} else if val, secondOk := object.(map[interface{}]interface{})[int(parsed)]; secondOk && parseErr == nil {
				object = val
			} else {
				object = nil
			}
		}
	}

	return object
}
