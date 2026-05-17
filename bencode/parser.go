package bencode

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// Decode reads the bencoded data from a bufio.Reader and returns the corresponding Go data structure.

// Here is my line by line interpretation of this func we are using a method called  Reader under the bufio package .
// Peek is another method under the bufio package.
// we are saving the output of the Peek to the variable b
// we are checking if the b is equal to 1
// if b is equal to 1 then we are returning the decoded interger value of that by calling the method decodeint which is defined below.
// if b is l then we are returning the decoded list value of that by calling the method decodelist which is defined below
// if b is d then we are returning the decoded dictionary value of that by calling the method decode dict which is defined below
// otherwise if b[0] is greater or equal to '0' and less than or equal to '9' we are returning the decoded string value of that by calling the method decodeString which is defined below
// if none of the above conditions are met then we are returning an error
// Why should I include 9 within the switch statement???

func Decode(r *bufio.Reader) (interface{}, error) {
	b, err := r.Peek(1)
	if err != nil {
		return nil, err
	}

	switch b[0] {
	case 'i':
		return decodeInt(r)
	case 'l':
		return decodeList(r)
	case 'd':
		return decodeDict(r)
	default:
		if b[0] >= '0' && b[0] <= '9' {
			return decodeString(r)
		}
		return nil, fmt.Errorf("Bencode: invalid type prefix %q, b[0]")
	}

}

// Let decode these statments line by line
// This func has the responisblity to decode the bencoded data
// if it is a recusive function - how should I train my self to detect the recursive function and how to code it .
// what does '_' signigfy here??
// What is an unused variable? decodeInt is an function not a variable
// I am looking at the right variable.
// Also I want to understand the hiearchy of a [function, variable, object, class]
// Like to an untrained eye - what can be inputs and what can be outputs .

func decodeInt(r *bufio.Reader) (int, error) {
	_, err := r.ReadByte() // consume 'i'
	if err != nil {
		return 0, err
	}

	valStr, err := r.ReadString('e')
	if err != nil {
		return 0, err
	}

	valStr = valStr[:len(valStr)-1] // remove 'e'
	return strconv.Atoi(valStr)
}

// Decode func for to decode a String
func decodeString(r *bufio.Reader) (string, error) {
	lenStr, err := r.ReadString(':')
	if err != nil {
		return "", err
	}

	lenStr = lenStr[:len(lenStr)-1]
	lenght, err := strconv.Atoi(lenStr)
	if err != nil {
		return "", err
	}

	buf := make([]byte, lenght)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func decodeList(r *bufio.Reader) ([]interface{}, error) {
	_, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	var list []interface{}
	for {
		b, err := r.Peek(1)
		if err != nil {
			return nil, err
		}
		if b[0] == 'e' {
			r.ReadByte() // consume 'e'
			return list, nil
		}

		val, err := Decode(r)
		if err != nil {
			return nil, err
		}

		list = append(list, val)
	}
}

func decodeDict(r *bufio.Reader) (map[string]interface{}, error) {
	_, err := r.ReadByte() // consume 'd'

	// My interpretation of this code is that we are creating a map of strings and interface values
	// The map is called dict and it is of type map[string]interface{}
	// Now a closer look to the if return statement.
	// if err != nil {} - means that if there is an error in reading the byte
	// then we are returning the error value of that by calling the method r.ReadByte
	// why did we delcare err before the if statement.
	// so before the type asseration symbol we declare two variables - one is the value and the other is the error
	// if the error is not nil then we are return
	//  ing the error value of that by calling the method r.ReadByte

	if err != nil {
		return nil, err
	}
	dict := make(map[string]interface{})
	for {
		b, err := r.Peek(1)
		if err != nil {
			return nil, err
		}
		if b[0] == 'e' {
			r.ReadByte() // consume 'e'
			return dict, nil
		}

		keyVal, err := decodeString(r)
		if err != nil {
			return nil, err
		}

		val, err := Decode(r)
		if err != nil {
			return nil, err
		}
		dict[keyVal] = val

	}

}
