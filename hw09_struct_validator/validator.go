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
	ruleValue string
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
		ruleType, ruleValue, found := strings.Cut(r, ":")

		if !found || ruleValue == "" {
			// TODO: program tag err
			return nil, fmt.Errorf("%w: empty tag '%s'", errParser, validateTag)
		}
		rule := validationRule{ruleType: ruleType, ruleValue: ruleValue}
		fmt.Println("Parsed rule:", rule)

		if !isAllowedRule(ruleType) {
			// TODO: program err
			return nil, fmt.Errorf("%w: unknown tag value '%s'", errParser, ruleType)
		}

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
				ruleMin, err := strconv.Atoi(r.ruleValue)
				if err != nil {
					return fmt.Errorf("%w: '%s' is not a number", errParser, r.ruleValue)
				}
				if inputVal < int64(ruleMin) {
					return fmt.Errorf("%w: %d < %d", errValidator, inputVal, ruleMin)
				}
			case MAX:
				ruleMax, err := strconv.Atoi(r.ruleValue)
				if err != nil {
					return fmt.Errorf("%w: '%s' is not a number", errParser, r.ruleValue)
				}
				if inputVal > int64(ruleMax) {
					return fmt.Errorf("%w: %d > %d", errValidator, inputVal, ruleMax)
				}
			case IN:
				var allowedValues []int64
				for _, part := range strings.Split(r.ruleValue, ",") {
					part = strings.TrimSpace(part)
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
				requiredLen, _ := strconv.Atoi(r.ruleValue)
				if len(inputVal) != requiredLen {
					return fmt.Errorf("%w: length != %d", errValidator, requiredLen)

				}
			case REGEXP:
				re, err := regexp.Compile(r.ruleValue)
				if err != nil {
					// TODO: program err
					return fmt.Errorf("%w: value not match regexp '%s'", errParser, r.ruleValue)
				}
				if !re.MatchString(inputVal) {
					return fmt.Errorf("%w: invalid regexp '%q'", errValidator, r.ruleValue)
				}
			case IN:
				var allowedValues []string
				for _, part := range strings.Split(r.ruleValue, ",") {
					part = strings.TrimSpace(part)
					allowedValues = append(allowedValues, part)
				}
				if !slices.Contains(allowedValues, inputVal) {
					return fmt.Errorf("%w: value not in %v", errValidator, allowedValues)
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
			fmt.Printf("Field %s : %s : %s\t\tValue %s\n", rf.Name, rf.Type, rf.Tag, rfValue)

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
