package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

type validationRule struct {
	ruleType  string
	ruleValue string
}

type Sample struct {
	Code int    `validate:"in:200,404,500"`
	Body string `json:"omitempty"`
}

func parseTag(stag reflect.StructTag) ([]validationRule, error) {
	validateTag, ok := stag.Lookup("validate")
	if !ok {
		return nil, nil
	}
	strRules := strings.Split(validateTag, "|")
	fmt.Println(strRules)

	var rules []validationRule

	for _, r := range strRules {
		ruleType, ruleValue, found := strings.Cut(r, ":")

		if !found {
			// TODO: program tag err
			continue
		}
		rule := validationRule{ruleType: ruleType, ruleValue: ruleValue}
		fmt.Println("Parsed rule:", rule)

		rules = append(rules, rule)
	}
	return rules, nil

}

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	if rt.Kind() == reflect.Struct {
		argCnt := rt.NumField()
		for i := 0; i < argCnt; i++ {
			rf := rt.Field(i)
			rfValue := rv.Field(i)
			fmt.Printf("Field %s : %s : %s\t\tValue %s\n", rf.Name, rf.Type, rf.Tag, rfValue)
			// fmt.Println(rf)

			if !rf.IsExported() {
				continue
			}
			_, _ = parseTag(rf.Tag)
		}
	}

	return nil
}
