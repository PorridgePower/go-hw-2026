package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const (
	IN     = "in"
	MIN    = "min"
	MAX    = "max"
	LENGTH = "len"
	REGEXP = "regexp"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var complexErr strings.Builder
	for _, e := range v {
		complexErr.WriteString(e.Field)
		complexErr.WriteString(": ")
		complexErr.WriteString(e.Err.Error())
		complexErr.WriteString("\n")
	}
	return complexErr.String()
}

type validationRule struct {
	ruleType  string
	ruleValue ruleValue
}

type ruleValue struct {
	minVal  *int64
	maxVal  *int64
	allowed []string
	regexp  *regexp.Regexp
	len     *int
}

var (
	errParser    = fmt.Errorf("invalid validate tag")
	errValidator = fmt.Errorf("value does not satisfy validation rule")
)

func isAllowedRule(ruleType string) bool {
	switch ruleType {
	case IN, MIN, MAX, LENGTH, REGEXP:
		return true
	default:
		return false
	}
}

func parseTag(stag reflect.StructTag) ([]validationRule, error) {
	validateTag, ok := stag.Lookup("validate")
	if !ok {
		return nil, nil
	}
	strRules := strings.Split(validateTag, "|")
	fmt.Println(strRules)

	rules := make([]validationRule, 0, len(strRules))

	for _, r := range strRules {
		ruleType, ruleValueStr, found := strings.Cut(r, ":")

		if !found || ruleValueStr == "" {
			// TODO: program tag err
			return nil, fmt.Errorf("%w: empty tag '%s'", errParser, validateTag)
		}

		var rv ruleValue
		switch ruleType {
		case MIN:
			ruleMin, err := strconv.ParseInt(ruleValueStr, 10, 64)
			if err != nil {
				// TODO: program err
				return nil, fmt.Errorf("%w: '%s' is not a number", errParser, ruleValueStr)
			}
			rv.minVal = &ruleMin
		case MAX:
			ruleMax, err := strconv.ParseInt(ruleValueStr, 10, 64)
			if err != nil {
				// TODO: program err
				return nil, fmt.Errorf("%w: '%s' is not a number", errParser, ruleValueStr)
			}
			rv.maxVal = &ruleMax

		case IN:
			var allowedValues []string
			for _, part := range strings.Split(ruleValueStr, ",") {
				part = strings.TrimSpace(part)
				allowedValues = append(allowedValues, part)
			}
			rv.allowed = allowedValues
		case REGEXP:
			re, err := regexp.Compile(ruleValueStr)
			if err != nil {
				// TODO: program err
				return nil, fmt.Errorf("%w: invalid regexp '%s'", errParser, ruleValueStr)
			}
			rv.regexp = re

		case LENGTH:
			requiredLen, err := strconv.Atoi(ruleValueStr)
			if err != nil {
				// TODO: program err
				return nil, fmt.Errorf("%w: '%s' is not a number", errParser, ruleValueStr)
			}
			rv.len = &requiredLen

		default:
			// TODO: program err
			return nil, fmt.Errorf("%w: unknown tag value '%s'", errParser, ruleType)

		}

		rule := validationRule{ruleType: ruleType, ruleValue: rv}
		fmt.Println("Parsed rule:", rule)

		rules = append(rules, rule)
	}
	return rules, nil
}

func validateItem(input reflect.Value, rules []validationRule) error {
	switch input.Kind() {
	case reflect.Int:
		inputVal := input.Int()
		for _, r := range rules {
			switch r.ruleType {
			case MIN:
				if inputVal < *r.ruleValue.minVal {
					return fmt.Errorf("%w: %d < %d", errValidator, inputVal, r.ruleValue.minVal)
				}
			case MAX:
				if inputVal > *r.ruleValue.maxVal {
					return fmt.Errorf("%w: %d > %d", errValidator, inputVal, r.ruleValue.maxVal)
				}
			case IN:
				var allowedValues []int64
				for _, part := range r.ruleValue.allowed {
					n, err := strconv.Atoi(part)
					if err != nil {
						// TODO: program err
						return fmt.Errorf("%w: '%s' is not a number", errParser, part)
					}
					allowedValues = append(allowedValues, int64(n))
				}
				if !slices.Contains(allowedValues, inputVal) {
					return fmt.Errorf("%w: value not in %v", errValidator, allowedValues)

				}

			}
		}

	case reflect.String:
		inputVal := input.String()
		for _, r := range rules {
			switch r.ruleType {
			case LENGTH:
				if len(inputVal) != *r.ruleValue.len {
					return fmt.Errorf("%w: length != %d", errValidator, r.ruleValue.len)

				}
			case REGEXP:
				re := r.ruleValue.regexp

				if !re.MatchString(inputVal) {
					return fmt.Errorf("%w: value not match regexp '%q'", errValidator, r.ruleValue.regexp)
				}
			case IN:
				if !slices.Contains(r.ruleValue.allowed, inputVal) {
					return fmt.Errorf("%w: value not in %v", errValidator, r.ruleValue.allowed)
				}
			}
		}

	case reflect.Slice:
		elemKind := input.Type().Elem().Kind()
		switch elemKind {
		case reflect.Int, reflect.String, reflect.Slice:
			for i := 0; i < input.Len(); i++ {
				elem := input.Index(i)
				if err := validateItem(elem, rules); err != nil {
					return err
				}
			}
		default:
		}
	}
	return nil
}

func Validate(v interface{}) error {
	var err ValidationErrors
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	if rt.Kind() == reflect.Struct {
		argCnt := rt.NumField()
		for i := 0; i < argCnt; i++ {
			rf := rt.Field(i)
			rfValue := rv.Field(i)

			if !rf.IsExported() {
				continue
			}
			tags, _ := parseTag(rf.Tag)
			elemErr := validateItem(rfValue, tags)
			if elemErr != nil {
				if errors.Is(elemErr, errParser) {
					return elemErr
				}
				err = append(err, ValidationError{Field: rf.Name, Err: elemErr})
			}
		}
	} else {
		return fmt.Errorf("v is not a struct")
	}

	if len(err) > 0 {
		return err
	}
	return nil
}
