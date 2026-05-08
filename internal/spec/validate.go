package spec

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var pluginFormat = regexp.MustCompile(`^[A-Za-z0-9_.-]+@[A-Za-z0-9_.-]+$`)

func Validate(p *Profile) error {
	v := validator.New()
	var errs []string

	if err := v.Struct(p); err != nil {
		ve, _ := err.(validator.ValidationErrors)
		for _, fe := range ve {
			errs = append(errs, formatFieldError(fe))
		}
	}

	for i, plugin := range p.Plugins {
		if !pluginFormat.MatchString(plugin) {
			errs = append(errs, fmt.Sprintf(
				"plugins[%d]: invalid format, expected \"name@marketplace\", got %q",
				i, plugin))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("profile validation failed:\n  - %s", strings.Join(errs, "\n  - "))
	}
	return nil
}

func formatFieldError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "oneof":
		return fmt.Sprintf("%s: must be one of [%s], got %q",
			toYAMLPath(fe.Namespace()), fe.Param(), fe.Value())
	case "required":
		return fmt.Sprintf("%s: required", toYAMLPath(fe.Namespace()))
	case "eq":
		return fmt.Sprintf("%s: must equal %q, got %q",
			toYAMLPath(fe.Namespace()), fe.Param(), fe.Value())
	default:
		return fmt.Sprintf("%s: failed %s validation", toYAMLPath(fe.Namespace()), fe.Tag())
	}
}

func toYAMLPath(ns string) string {
	parts := strings.Split(ns, ".")
	if len(parts) > 1 {
		parts = parts[1:]
	}
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToLower(part[:1]) + part[1:]
	}
	return strings.Join(parts, ".")
}
