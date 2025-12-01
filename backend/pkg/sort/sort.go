package sort

import (
	"fmt"
	"strings"
)

type Sort struct {
	Name      string
	Direction string
}

type Sorts struct {
	sorts []Sort
	set   map[string]struct{}
}

func NewSorts() Sorts {
	return Sorts{
		set: make(map[string]struct{}),
	}
}

func (sorts *Sorts) Add(s Sort) {
	if _, ok := sorts.set[s.Name]; ok {
		return
	}

	sorts.sorts = append(sorts.sorts, s)
	sorts.set[s.Name] = struct{}{}
}

func (sorts Sorts) Data() []Sort {
	return sorts.sorts
}

func (sorts Sorts) GetFields() []string {
	fields := make([]string, len(sorts.sorts))
	for i, sort := range sorts.sorts {
		fields[i] = sort.Name
	}
	return fields
}

func (sorts Sorts) Validate(allowedColumns []string) error {
	if len(allowedColumns) == 0 {
		return nil
	}

	allowedSet := make(map[string]bool)
	for _, col := range allowedColumns {
		allowedSet[col] = true
	}

	var invalidFields []string
	for _, sortField := range sorts.GetFields() {
		if !allowedSet[sortField] {
			invalidFields = append(invalidFields, sortField)
		}
	}

	if len(invalidFields) > 0 {
		return fmt.Errorf("invalid sort fields: %v. Allowed fields: %v", invalidFields, allowedColumns)
	}

	return nil
}

func ParseSort(value string) Sort {
	if strings.TrimSpace(value) == "" {
		return Sort{}
	}
	direction := "asc"
	if string(value[0]) == "-" {
		direction = "desc"
	}
	return Sort{
		Name:      value[1:],
		Direction: direction,
	}
}

func ParseMultipleSorts(values []string) Sorts {
	result := NewSorts()

	for _, val := range values {
		result.Add(ParseSort(val))
	}

	return result
}
