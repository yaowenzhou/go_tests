// For structure variable search matching
// Pointers and other complex member variables
// are not supported temporarily

package search

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/spf13/cast"
	"golang.org/x/exp/constraints"
)

type SearchOperator int32

const (
	SEARCH_OPERATOR_UNKNOW        SearchOperator = 0 // not use
	SEARCH_OPERATOR_CONTAIN_OR    SearchOperator = 1 // contain(fuzzy search)
	SEARCH_OPERATOR_LESS          SearchOperator = 2 // less than
	SEARCH_OPERATOR_LESS_EQUAL    SearchOperator = 3 // less than or equal
	SEARCH_OPERATOR_EQUAL         SearchOperator = 4 // equal
	SEARCH_OPERATOR_GREATER_EQUAL SearchOperator = 5 // greater than or equal
	SEARCH_OPERATOR_GREATER       SearchOperator = 6 // greater than
	SEARCH_OPERATOR_NOT_EQUAL     SearchOperator = 7 // not equal
	SEARCH_OPERATOR_NOT_CONTAIN   SearchOperator = 8 // not contain
)

var searchOperatorMap map[string]SearchOperator
var searchOperatorName []string

func init() {
	searchOperatorMap = make(map[string]SearchOperator)
	searchOperatorMap["c"] = SEARCH_OPERATOR_CONTAIN_OR
	searchOperatorMap["contain"] = SEARCH_OPERATOR_CONTAIN_OR
	searchOperatorMap["lt"] = SEARCH_OPERATOR_LESS
	searchOperatorMap["lte"] = SEARCH_OPERATOR_LESS_EQUAL
	searchOperatorMap["eq"] = SEARCH_OPERATOR_EQUAL
	searchOperatorMap["gte"] = SEARCH_OPERATOR_GREATER_EQUAL
	searchOperatorMap["gt"] = SEARCH_OPERATOR_GREATER
	searchOperatorMap["neq"] = SEARCH_OPERATOR_NOT_EQUAL
	searchOperatorMap["nc"] = SEARCH_OPERATOR_NOT_CONTAIN
	searchOperatorMap["notcontain"] = SEARCH_OPERATOR_NOT_CONTAIN
	searchOperatorName = []string{
		"unknow",
		"contain",
		"less than",
		"less than or equal",
		"equal",
		"equal or greater than",
		"greater",
		"not equal",
		"not contain",
	}
}

// searchLimit search limit
type searchLimit struct {
	// SearchOperators Supported search operations
	SearchOperators []SearchOperator
	// Error error message
	// this error will be thrown when using validCheck to check that an operator is invalid
	Error error
}

type Searcher struct {
	Field          string         // field name
	Value          string         // the value of field
	SearchOperator SearchOperator // search operator
	fieldKind      reflect.Kind   // kind of field's type
	offset         uintptr        // offset of field's in struct
	value          interface{}    // filter value for match
}

type intface struct {
	typ   unsafe.Pointer
	value unsafe.Pointer
}

type SearcherLimit struct {
	limit            map[string]*searchLimit
	structType       unsafe.Pointer // save struct's type
	fieldIndexMap    map[string]int // save field's offset in struct
	defaultStructVar interface{}
}

// getSearchOperator get search operator and check if it is valid
func getSearchOperator(
	fieldKind reflect.Kind, searchOperatorStr string, jsonTag string,
) (searchOperator SearchOperator, err error) {
	s, ok := searchOperatorMap[searchOperatorStr] // get search operator and find the unreasonable configuration
	if !ok {
		return 0, fmt.Errorf("not support search type(%s)", searchOperatorStr)
	}
	if (fieldKind >= reflect.Int && fieldKind <= reflect.Uint64) ||
		fieldKind == reflect.Float32 || fieldKind == reflect.Float64 { // Only </<=/=/>=/>/!= is allowed for numeric type
		if s != SEARCH_OPERATOR_LESS && s != SEARCH_OPERATOR_LESS_EQUAL &&
			s != SEARCH_OPERATOR_EQUAL && s != SEARCH_OPERATOR_GREATER_EQUAL &&
			s != SEARCH_OPERATOR_GREATER && s != SEARCH_OPERATOR_NOT_EQUAL {
			return 0, fmt.Errorf("field(%s) is number type, not support search type(%s)", jsonTag, searchOperatorStr)
		}
		return s, nil // Record the currently allowed search operators
	}
	if fieldKind == reflect.String { // only contain/=/!=/not contain is allowed for string type
		if s != SEARCH_OPERATOR_CONTAIN_OR && s != SEARCH_OPERATOR_EQUAL &&
			s != SEARCH_OPERATOR_NOT_EQUAL && s != SEARCH_OPERATOR_NOT_CONTAIN {
			return 0, fmt.Errorf("field(%s) is string type, not support search type(%s)", jsonTag, searchOperatorStr)
		}
		return s, nil // Record the currently allowed search operators
	}
	// Types other than numbers and strings
	// do not report errors for the time being,
	// but return 0 to make them unusable
	return SEARCH_OPERATOR_UNKNOW, nil // numeric
}

// NewSearcherLimit Construct a searcher for structure search and judgment
func NewSearcherLimit(i interface{}) (*SearcherLimit, error) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, errors.New("param i must be a struct or a pointer of struct")
	}
	structType := (*intface)(unsafe.Pointer(&i)).typ
	limit := make(map[string]*searchLimit)
	fieldIndexMap := make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		searchTag := tag.Get("search")
		if searchTag == "" {
			continue
		}
		jsonTag := tag.Get("json")
		if jsonTag == "" {
			continue
		}
		sLimit := &searchLimit{}
		searchOperatorStrs := strings.Split(searchTag, ",")
		searchOperatorDuplicateMap := make(map[string]bool)
		fKind := v.Field(i).Type().Kind()
		for _, sStr := range searchOperatorStrs {
			sStr = strings.TrimSpace(sStr)
			if _, ok := searchOperatorDuplicateMap[sStr]; ok { // duplication search operator
				continue
			}
			searchOperatorDuplicateMap[sStr] = true
			s, err := getSearchOperator(fKind, sStr, jsonTag)
			if err != nil {
				return nil, err
			}
			if s != SEARCH_OPERATOR_UNKNOW {
				sLimit.SearchOperators = append(sLimit.SearchOperators, s)
			}
		}
		sStrNew := make([]string, len(sLimit.SearchOperators))
		for k, s := range sLimit.SearchOperators {
			sStrNew[k] = searchOperatorName[s]
		}
		// full error msg
		sLimit.Error = fmt.Errorf("field(%s) only support search operate: %s",
			jsonTag, strings.Join(sStrNew, "/"))
		limit[jsonTag] = sLimit
		fieldIndexMap[jsonTag] = i
	}
	return &SearcherLimit{
		limit:            limit,
		structType:       structType,
		fieldIndexMap:    fieldIndexMap,
		defaultStructVar: i,
	}, nil
}

// getFieldOffsetAndType get field's type and offset
func (s *Searcher) getFieldOffsetAndType(
	in interface{}, limit *SearcherLimit, fieldIndex int,
) (err error) {
	intf := (*intface)(unsafe.Pointer(&in))
	if intf.typ != limit.structType {
		return fmt.Errorf("in's type is invalid")
	}
	s.fieldKind = reflect.ValueOf(in).Elem().Field(fieldIndex).Kind()
	// the offset of in.(s.Field)
	s.offset = reflect.TypeOf(in).Elem().Field(fieldIndex).Offset
	return
}

// getFilterValue Converts the value (string) used as a search
// to a value of the corresponding type
func (s *Searcher) genFilterValue() {
	switch s.fieldKind {
	case reflect.Int:
		s.value = cast.ToInt(s.Value)
	case reflect.Int8:
		s.value = cast.ToInt8(s.Value)
	case reflect.Int16:
		s.value = cast.ToInt16(s.Value)
	case reflect.Int32:
		s.value = cast.ToInt32(s.Value)
	case reflect.Int64:
		s.value = cast.ToInt64(s.Value)
	case reflect.Uint:
		s.value = cast.ToUint(s.Value)
	case reflect.Uint8:
		s.value = cast.ToUint8(s.Value)
	case reflect.Uint16:
		s.value = cast.ToUint16(s.Value)
	case reflect.Uint32:
		s.value = cast.ToUint32(s.Value)
	case reflect.Uint64:
		s.value = cast.ToUint64(s.Value)
	case reflect.Float32:
		s.value = cast.ToFloat32(s.Value)
	case reflect.Float64:
		s.value = cast.ToFloat64(s.Value)
	case reflect.String:
		s.value = s.Value
	}
}

// ValidCheck Search operator validity check
func (s *SearcherLimit) ValidCheck(infos []*Searcher) (err error) {
	for k, info := range infos {
		searchLimit, ok := s.limit[info.Field]
		if !ok {
			return fmt.Errorf("field(%s) does not support search", info.Field)
		}
		invalid := true
		for _, so := range searchLimit.SearchOperators {
			if so == info.SearchOperator { // info.SearchOperator is valid
				invalid = false
				break
			}
		}
		if invalid { // invalid message
			return fmt.Errorf("searchers[%d] is invalid, %s", k, searchLimit.Error)
		}
		err = info.getFieldOffsetAndType(s.defaultStructVar, s, s.fieldIndexMap[info.Field])
		if err != nil {
			return err
		}
		info.genFilterValue()
	}
	return nil
}

func containCompare[K1 ~string](left, right K1) bool {
	return strings.Contains(string(left), string(right))
}
func ltCompare[K1 constraints.Ordered](left, right K1) bool { // 泛型函数
	return left < right
}
func lteCompare[K1 constraints.Ordered](left, right K1) bool { // 泛型函数
	return left <= right
}
func eqCompare[K1 constraints.Ordered](left, right K1) bool { // 泛型函数
	return left == right
}
func gteCompare[K1 constraints.Ordered](left, right K1) bool { // 泛型函数
	return left >= right
}
func gtCompare[K1 constraints.Ordered](left, right K1) bool { // 泛型函数
	return left != right
}
func neqCompare[K1 constraints.Ordered](left, right K1) bool { // 泛型函数
	return left != right
}
func notContainCompare[K1 ~string](left, right K1) bool { // 泛型函数
	return !strings.Contains(string(left), string(right))
}

// doNumbericMatch numberic match
func doNumbericMatch[K1 constraints.Integer | constraints.Float](
	left, right K1, searchOperator SearchOperator) bool {
	switch searchOperator {
	case SEARCH_OPERATOR_LESS:
		return ltCompare(left, right)
	case SEARCH_OPERATOR_LESS_EQUAL:
		return lteCompare(left, right)
	case SEARCH_OPERATOR_EQUAL:
		return eqCompare(left, right)
	case SEARCH_OPERATOR_GREATER_EQUAL:
		return gteCompare(left, right)
	case SEARCH_OPERATOR_GREATER:
		return gtCompare(left, right)
	case SEARCH_OPERATOR_NOT_EQUAL:
		return neqCompare(left, right)
	}
	return false
}

// doStringMatch string match
func doStringMatch[K1 ~string](left, right K1, searchOperator SearchOperator) bool {
	switch searchOperator {
	case SEARCH_OPERATOR_CONTAIN_OR:
		return containCompare(left, right)
	case SEARCH_OPERATOR_EQUAL:
		return eqCompare(left, right)
	case SEARCH_OPERATOR_NOT_EQUAL:
		return neqCompare(left, right)
	case SEARCH_OPERATOR_NOT_CONTAIN:
		return notContainCompare(left, right)
	}
	return false
}

// match Check whether the data meets the search condition
func (s *Searcher) match(in interface{}) bool {
	dataPtr := unsafe.Pointer(uintptr((*intface)(unsafe.Pointer(&in)).value) + s.offset)
	switch s.fieldKind {
	case reflect.Int:
		return doNumbericMatch(*(*int)(dataPtr), s.value.(int), s.SearchOperator)
	case reflect.Int8:
		return doNumbericMatch(*(*int8)(dataPtr), s.value.(int8), s.SearchOperator)
	case reflect.Int16:
		return doNumbericMatch(*(*int16)(dataPtr), s.value.(int16), s.SearchOperator)
	case reflect.Int32:
		return doNumbericMatch(*(*int32)(dataPtr), s.value.(int32), s.SearchOperator)
	case reflect.Int64:
		return doNumbericMatch(*(*int64)(dataPtr), s.value.(int64), s.SearchOperator)
	case reflect.Uint:
		return doNumbericMatch(*(*uint)(dataPtr), s.value.(uint), s.SearchOperator)
	case reflect.Uint8:
		return doNumbericMatch(*(*uint8)(dataPtr), s.value.(uint8), s.SearchOperator)
	case reflect.Uint16:
		return doNumbericMatch(*(*uint16)(dataPtr), s.value.(uint16), s.SearchOperator)
	case reflect.Uint32:
		return doNumbericMatch(*(*uint32)(dataPtr), s.value.(uint32), s.SearchOperator)
	case reflect.Uint64:
		return doNumbericMatch(*(*uint64)(dataPtr), s.value.(uint64), s.SearchOperator)
	case reflect.Float32:
		return doNumbericMatch(*(*float32)(dataPtr), s.value.(float32), s.SearchOperator)
	case reflect.Float64:
		return doNumbericMatch(*(*float64)(dataPtr), s.value.(float64), s.SearchOperator)
	case reflect.String:
		return doStringMatch(*(*string)(dataPtr), s.value.(string), s.SearchOperator)
	}
	return false
}

// Filter filter datas and return filtered datas
func (s *Searcher) Filter(
	limit *SearcherLimit, datasIn []interface{},
) (datasOut []interface{}, err error) {
	if len(datasIn) == 0 {
		return nil, nil
	}
	for i := 0; i < len(datasIn); i++ {
		if datasIn[0] == nil {
			return nil, fmt.Errorf("datasIn[%d] is nil", i)
		}
		if reflect.ValueOf(datasIn[i]).Pointer() == 0 {
			return nil, fmt.Errorf("datasIn[%d] is a nil pointer", i)
		}
		typ := (*intface)(unsafe.Pointer(&datasIn[i])).typ
		// Check whether the type of dataIn [i] is the same as that recorded in the parameter limit
		if typ != limit.structType {
			return nil, fmt.Errorf("datasIn[%d]'s type is invalid", i)
		}
		if s.match(datasIn[i]) {
			datasOut = append(datasOut, datasIn[i])
		}
	}
	return
}
