package cheapjson

import (
	"errors"
	"fmt"
	"bytes"
	"strconv"
)

func unexpected(expect string, offset, size int, data []byte) error {
	if offset < size {
		return errors.New(fmt.Sprintf("Unexptected token '%c' at: %d, expect: %s", data[offset], offset, expect))
	} else {
		return errors.New("Unexpected EOF, expect: " + expect)
	}
}

var (
	valueStart = "{, [, [0-9], -, t, f, n, \""
	bytesTrue = []byte{'r', 'u', 'e'}
	bytesFalse = []byte{'a', 'l', 's', 'e'}
	bytesNull = []byte{'u', 'l', 'l'}
)

const (
	// start of a value
	stateNone = iota
	stateString
	// after [ must be a value or ]
	stateArrayValueOrEnd
	// after a value, must be a , or ]
	stateArrayEndOrComma
	// after a {, must be a key string or }
	stateObjectKeyOrEnd
	// after a key string must be a :
	stateObjectColon
	// after a : must be a value
	// after a value, must be , or }
	stateObjectEndOrComma
	// after a , must be key string
	stateObjectKey
)

type state struct {
	value  *Value
	state  int
	parent *state
}

func addBuf(buf []byte, tempInt2, bufSize, ask int) ([]byte, int) {
	if bufSize < tempInt2 + ask {
		bufSize += 1024
		nbuf := make([]byte, bufSize)
		copy(nbuf, buf)
		return nbuf, bufSize
	}
	return buf, bufSize
}

func Unmarshal(data []byte) (value *Value, err error) {
	value = &Value{nil}
	root := &state{value, stateNone, nil}
	curr := root
	size := len(data)
	offset := 0
	bufSize := 1024
	buf := make([]byte, bufSize)
	tempUnicode := make([]int, 4)
	var tempInt int
	var tempInt2 int
	var tempInt3 int
	var tempInt4 int
	var tempByte byte
	var tempDecimal []byte
	var tempExp []byte
	for {
		// any loop start should check the whitespace
		LOOP_WHITESPACE:
		for ; offset < size; offset++ {
			switch data[offset] {
			case '\t', '\r', '\n', ' ':
				continue
			default:
				break LOOP_WHITESPACE
			}
		}
		if curr == nil {
			// must end
			if offset != size {
				err = unexpected("EOF", offset, size, data)
			}
			// NO else, check it according to the context to
			// get detailed information
			return
		}
		switch curr.state {
		case stateArrayValueOrEnd:
			if offset == size {
				err = unexpected("value or ]", offset, size, data)
				return
			}
			switch data[offset] {
			case ']':
				curr = curr.parent
				offset++
			default:
				curr.state = stateArrayEndOrComma
				curr = &state{curr.value.AddElement(), stateNone, curr}
			}
			continue
		case stateArrayEndOrComma:
			if offset == size {
				err = unexpected(", or ]", offset, size, data)
				return
			}
			switch data[offset] {
			case ']':
				offset++
				curr = curr.parent
			case ',':
				offset++
				curr.state = stateArrayEndOrComma
				curr = &state{curr.value.AddElement(), stateNone, curr}
			default:
				err = unexpected(", or ]", offset, size, data)
				return
			}
			continue
		case stateObjectColon:
			if offset == size {
				err = unexpected(":", offset, size, data)
				return
			}
			offset++
			curr.state = stateNone
			continue
		case stateObjectEndOrComma:
			if offset == size {
				err = unexpected(", or }", offset, size, data)
				return
			}
			switch data[offset] {
			case ',':
				curr.state = stateObjectKey
			case '}':
				curr = curr.parent
			default:
				err = unexpected(", or }", offset, size, data)
				return
			}
			offset++
			continue
		case stateObjectKeyOrEnd:
			if offset == size {
				err = unexpected("\" or }", offset, size, data)
				return
			}
			if data[offset] == '}' {
				offset++
				curr = curr.parent
				continue
			}
			fallthrough
		case stateObjectKey:
			if offset == size {
				err = unexpected("\"", offset, size, data)
				return
			}
			if data[offset] != '"' {
				err = unexpected("\"", offset, size, data)
				return
			}
			fallthrough
		case stateString:
			offset++
			tempInt2 = 0
			LOOP_STRING:
			for tempInt = offset; tempInt < size; tempInt++ {
				switch data[tempInt] {
				case '\\':
					tempInt++
					if tempInt == size {
						err = unexpected("escaped char", tempInt, size, data)
						return
					}
					switch data[tempInt] {
					case 'U', 'u':
						tempInt++
						if size < tempInt + 4 {
							err = unexpected("[0-F]", size, size, data)
							return
						}
						for tempInt3 = 0; tempInt3 < 4; tempInt3++ {
							switch data[tempInt] {
							case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
								tempUnicode[tempInt3] = int(data[tempInt]) - 0x30
							case 'a', 'b', 'c', 'd', 'e', 'f':
								tempUnicode[tempInt3] = int(data[tempInt]) - 0x57
							case 'A', 'B', 'C', 'D', 'E', 'F':
								tempUnicode[tempInt3] = int(data[tempInt]) - 0x37
							default:
								err = unexpected("[0-F]", tempInt, size, data)
								return
							}
							tempInt++
						}
						tempInt4 = (tempUnicode[0] << 12) | (tempUnicode[1] << 8) | (tempUnicode[2] << 4) | (tempUnicode[3])
						if tempInt4 > 0xD7FF && tempInt4 < 0xDC00 {
							// need next utf-16 part
							if size < tempInt + 6 {
								if size == tempInt || data[tempInt] != '\\' {
									err = unexpected("\\", size, size, data)
									return
								}
								if size < tempInt + 2 || (data[tempInt + 1] != 'U' && data[tempInt + 1] != 'u') {
									err = unexpected("Uu", tempInt + 1, size, data)
									return
								}
								err = unexpected("[0-F]", size, size, data)
								return
							}
							if data[tempInt] != '\\' {
								err = unexpected("\\", tempInt, size, data)
								return
							}
							tempInt++
							if data[tempInt] != 'U' && data[tempInt] != 'u' {
								err = unexpected("Uu", tempInt, size, data)
								return
							}
							tempInt++
							for tempInt3 = 0; tempInt3 < 4; tempInt3++ {
								switch data[tempInt] {
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									tempUnicode[tempInt3] = int(data[tempInt]) - 0x30
								case 'a', 'b', 'c', 'd', 'e', 'f':
									tempUnicode[tempInt3] = int(data[tempInt]) - 0x57
								case 'A', 'B', 'C', 'D', 'E', 'F':
									tempUnicode[tempInt3] = int(data[tempInt]) - 0x37
								default:
									err = unexpected("[0-F]", tempInt, size, data)
									return
								}
								tempInt++
							}
							tempInt3 = (tempUnicode[0] << 12) | (tempUnicode[1] << 8) | (tempUnicode[2] << 4) | (tempUnicode[3])
							if tempInt3 < 0xDC00 || tempInt3 > 0xDFFF {
								err = unexpected("[0xdc00 - 0xdfff]", tempInt - 4, size, data)
								return
							}
							tempInt4 = (((tempInt4 - 0xD800) << 10) | (tempInt3 - 0xDC00)) + 0x10000
						}
						tempInt--
						if tempInt4 < 0x0080 {
							buf, bufSize = addBuf(buf, tempInt2, bufSize, 1)
							buf[tempInt2] = byte(tempInt4)
							tempInt2++
						} else if (tempInt4 < 0x0800) {
							buf, bufSize = addBuf(buf, tempInt2, bufSize, 2)
							buf[tempInt2] = 0xC0 | byte(tempInt4 >> 6)
							buf[tempInt2 + 1] = 0x80 | byte(tempInt4 & 0xBF)
							tempInt2 += 2
						} else if (tempInt4 < 0x10000) {
							buf, bufSize = addBuf(buf, tempInt2, bufSize, 3)
							buf[tempInt2] = 0xE0 | byte(tempInt4 >> 12)
							buf[tempInt2 + 1] = 0x80 | byte((tempInt4 >> 6) & 0xBF)
							buf[tempInt2 + 2] = 0x80 | byte(tempInt4 & 0xBF)
							tempInt2 += 3
						} else {
							buf, bufSize = addBuf(buf, tempInt2, bufSize, 4)
							buf[tempInt2] = 0xF0 | byte(tempInt4 >> 18)
							buf[tempInt2 + 1] = 0x80 | byte((tempInt4 >> 12) & 0xBF)
							buf[tempInt2 + 2] = 0x80 | byte((tempInt4 >> 6) & 0xBF)
							buf[tempInt2 + 3] = 0x80 | byte(tempInt4 & 0xBF)
							tempInt2 += 4
						}
					case 't':
						buf, bufSize = addBuf(buf, tempInt2, bufSize, 1)
						buf[tempInt2] = '\t'
						tempInt2++
					case 'r':
						buf, bufSize = addBuf(buf, tempInt2, bufSize, 1)
						buf[tempInt2] = '\r'
						tempInt2++
					case 'n':
						buf, bufSize = addBuf(buf, tempInt2, bufSize, 1)
						buf[tempInt2] = '\n'
						tempInt2++
					case '"':
						buf, bufSize = addBuf(buf, tempInt2, bufSize, 1)
						buf[tempInt2] = '"'
						tempInt2++
					case '\\':
						buf, bufSize = addBuf(buf, tempInt2, bufSize, 1)
						buf[tempInt2] = '\\'
						tempInt2++
					case '/':
						buf, bufSize = addBuf(buf, tempInt2, bufSize, 1)
						buf[tempInt2] = '/'
						tempInt2++
					case 'b':
						buf, bufSize = addBuf(buf, tempInt2, bufSize, 1)
						buf[tempInt2] = 0x08
						tempInt2++
					case 'f':
						buf, bufSize = addBuf(buf, tempInt2, bufSize, 1)
						buf[tempInt2] = 0x0C
						tempInt2++
					default:
						err = unexpected("escape sequence", tempInt, size, data)
						return
					}
				case '"':
					break LOOP_STRING
				case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
					16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31:
					err = unexpected("unicode", tempInt, size, data)
					return
				default:
					buf, bufSize = addBuf(buf, tempInt2, bufSize, 1)
					buf[tempInt2] = data[tempInt]
					tempInt2++
				}
			}
			if tempInt == size {
				err = unexpected("\" to end string", tempInt, size, data)
				return
			}
			if curr.state == stateString {
				curr.value.value = string(buf[0:tempInt2])
				curr = curr.parent
			} else {
				curr.state = stateObjectEndOrComma
				curr = &state{curr.value.AddField(string(buf[0:tempInt2])), stateObjectColon, curr}
			}
			offset = tempInt + 1
			continue
		default:
			// the start of a value
			// stateNone
			if offset == size {
				err = unexpected("value", offset, size, data)
				return
			}
			switch data[offset] {
			case '{':
				curr.state = stateObjectKeyOrEnd
				curr.value.value = map[string]*Value{}
				offset++
				continue
			case '[':
				curr.state = stateArrayValueOrEnd
				curr.value.value = []*Value{}
				offset++
				continue
			case '"':
				curr.state = stateString
				continue
			case '0', '1', '2', '3', '4',
				'5', '6', '7', '8', '9', '-':
				tempDecimal = nil
				tempExp = nil
				tempInt4 = offset
				// get negative
				if data[offset] == '-' {
					offset++
				}
				LOOP_NUM_INT: // read number.integer part
				for tempInt = offset; tempInt < size; tempInt++ {
					switch data[tempInt] {
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						continue
					default:
						break LOOP_NUM_INT
					}
				}
				tempInt2 = tempInt - offset
				// count of integer part
				if tempInt2 == 0 {
					// this will occur when start with -
					err = unexpected("[0-9]", offset, size, data)
					return
				}
				if data[offset] == '0' {
					if tempInt2 != 1 {
						// 0 MUST only one
						offset++
						err = unexpected("[.eE]", offset, size, data)
						return
					}
				}
				offset = tempInt
				if offset < size && data[offset] == '.' {
					// has decimal
					offset++
					LOOP_NUM_DEC:
					for tempInt = offset; tempInt < size; tempInt++ {
						switch data[tempInt] {
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							continue
						default:
							break LOOP_NUM_DEC
						}
					}
					if tempInt == offset {
						// MUST contains decimal
						err = unexpected("[0-9]", offset, size, data)
						return
					}
					tempDecimal = data[offset:tempInt]
					offset = tempInt
				}
				if offset < size && (data[offset] == 'e' || data[offset] == 'E') {
					// has exponent
					offset++
					if offset == size {
						// need to check EOF for the leading +/-
						err = unexpected("[0-9]", offset, size, data)
						return
					}
					if data[offset] == '-' || data[offset] == '+' {
						offset++
					}

					LOOP_NUM_EXP:
					for tempInt = offset; tempInt < size; tempInt++ {
						switch data[tempInt] {
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							continue
						default:
							break LOOP_NUM_EXP
						}
					}
					if tempInt == offset {
						err = unexpected("[0-9]", offset, size, data)
						return
					}
					// There do not need to check the leading 0 according to the spec
					tempExp = data[offset:tempInt]
					offset = tempInt
				}
				// just simplify the cases, but we may need to confirm 1e3 is a integer
				// rather than a float value
				if tempDecimal == nil && tempExp == nil {
					curr.value.value, err = strconv.ParseInt(string(data[tempInt4:offset]), 10, 64)
				} else {
					curr.value.value, err = strconv.ParseFloat(string(data[tempInt4:offset]), 64)
				}
				if err != nil {
					return
				}
				// NORMAL to here
				curr = curr.parent
				continue
			case 'n':
				offset++
				if size < offset + 3 {
					expect := bytesNull[size - offset]
					err = unexpected(string(expect), size, size, data)
					return
				}
				if bytes.Equal(data[offset:offset + 3], bytesNull) {
					offset += 3
					curr.value.value = NULL
					curr = curr.parent
					continue
				}
				err = unexpected("null", offset, size, data)
				return
			case 't':
				offset++
				if size < offset + 3 {
					tempByte = bytesTrue[size - offset]
					err = unexpected(string(tempByte), size, size, data)
					return
				}
				if bytes.Equal(data[offset:offset + 3], bytesTrue) {
					offset += 3
					curr.value.value = true
					curr = curr.parent
					continue
				}
				err = unexpected("true", offset, size, data)
				return
			case 'f':
				offset++
				if size < offset + 4 {
					tempByte = bytesFalse[size - offset]
					err = unexpected(string(tempByte), size, size, data)
					return
				}
				if bytes.Equal(data[offset:offset + 4], bytesFalse) {
					offset += 4
					curr.value.value = false
					curr = curr.parent
					continue
				}
				err = unexpected("false", offset, size, data)
				return
			default:
				err = unexpected(valueStart, offset, size, data)
				return
			}
		}
	}
}
