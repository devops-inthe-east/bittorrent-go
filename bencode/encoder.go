package bencode

import (
	"bytes"
	"fmt"
	"sort"
)

// Marshal takes a Go data structure and packs it
// into a bencoded byte slice

// So is it safe to assume that Marshal method can be used
// to get the package value of a bencoded dictionary/ list/
// integer or a string???

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := marshalValue(&buf, v)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// marshalValue inspects the Mystery Box and packs it
// appropriately

func marshalValue(w *bytes.Buffer, v interface{}) error {
	// we use a `switch` statement to check the type.
	// This is called a **Type Switch**. Agreed
	// However after determining the value - what is purpose of
	// having v.(type) here?

	switch val := v.(type) {
	case string:
		// What is the purpose of having the % has a prefix for
		// letter like d/s/i ???
		// I am asking this because in the fmt package - we are having this %s for
		// string , %d for integer , %v for any type ,
		// What what special significance does this % sign play here.
		w.WriteString(fmt.Sprintf("%d:%s", len(val), val))

	case int:
		w.WriteString(fmt.Sprintf("i%de", val))

		//	return nil

		// We are trying to open the Mystery box
		// if the value of the interface is of type []interface{}:
		// Then the item in the slice is also of type interface{}
		// So that item is packed by calling the marshalValue method recursively
		//
	case []interface{}:
		w.WriteString("l")

		// Range loops through the list. '_' is the index,
		// 'item' is the value

		for _, item := range val {
			err := marshalValue(w, item)
			if err != nil {
				return err
			}
		}
		w.WriteString("e")

	case map[string]interface{}:
		w.WriteString("d")

		// 1. Extract Keys
		var keys []string
		for k := range val {
			keys = append(keys, k)
			// 	Again I am struggling to building a mental model to
			// 	Assiging a types of bracket to a variable/ object/ method/ struct
			//

		}

		// Sort the keys alphabetically to ensure consistent ordering
		sort.Strings(keys)

		// 3. Lopps through sorted keys and pack the key-value pairs
		for _, k := range keys {
			// Pack the string key
			err := marshalValue(w, k)
			err = marshalValue(w, val[k]) // Pack the Value
			// We have two if statements here. One for the key and one for the value
			if err != nil {
				return err
			}
		}
		w.WriteString("e")
	default:
		return fmt.Errorf("encoder doesn't support type: %T", v)
	}
	return nil

}
