package klash

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var HelpError = errors.New("Help error")

type Help struct {
	Params *Params
	Name   string
	Desc   string
}

func NewHelpFromParams(name, desc string, params *Params) *Help {
	return &Help{
		Params: params,
		Name:   name,
		Desc:   desc,
	}
}

func NewHelp(name, desc string, rawParams interface{}) (*Help, error) {
	params, err := NewParams(rawParams)
	if err != nil {
		return nil, err
	}
	return NewHelpFromParams(name, desc, params), nil
}

func (h *Help) Usage() string {
	usages := make([]string, 0, len(h.Params.Listing)+2)
	usages = append(usages, "[-h]")

	for _, parameter := range h.Params.Listing {
		var paramUsage string
		shortName := DecomposeName(parameter.Alias, true)

		if shortName == "" {
			shortName = DecomposeName(parameter.Name, true)
		}
		shortName = Dashed(shortName)

		if parameter.Value.Type().Kind() != reflect.Bool {
			paramUsage = fmt.Sprintf("[%s %s]",
				shortName,
				DecomposeName(parameter.Name, false),
			)
		} else {
			paramUsage = fmt.Sprintf("[%s]", shortName)
		}

		usages = append(usages, paramUsage)
	}
	usages = append(usages, "ARGS...")
	return strings.Join(usages, " ")
}

func (h *Help) Details() string {
	detail := ""
	length := len(h.Params.Listing)

	if length > 0 {
		helpArgs := "-h, --help=false"
		maxLength := len(helpArgs)

		details := make([][2]string, 0, length+2)

		detail = fmt.Sprintf("argument details:\n")
		details = append(details, [2]string{helpArgs, "Show this help"})

		for _, parameter := range h.Params.Listing {
			var paramName string

			name := fmt.Sprintf("%s", Dashed(DecomposeName(parameter.Name, true)))
			if parameter.Alias != "" {
				name = fmt.Sprintf("%s, %s", Dashed(DecomposeName(parameter.Alias, true)), name)
			}

			switch parameter.Value.Type().Kind() {
			case reflect.Bool:
				paramName = fmt.Sprintf("%s=%t", name, parameter.Value.Interface())
			case reflect.Slice:
				paramName = fmt.Sprintf("%s=[]", name)
			case reflect.String:
				paramName = fmt.Sprintf("%s=\"%s\"", name, parameter.Value.Interface())
			case reflect.Int, reflect.Uint:
				paramName = fmt.Sprintf("%s=%d", name, parameter.Value.Interface())
			case reflect.Float32, reflect.Float64:
				paramName = fmt.Sprintf("%s=%.3f", name, parameter.Value.Interface())
			default:
				paramName = fmt.Sprintf("%s=%s", name, parameter.Value.Interface())
			}

			paramLength := len(paramName)
			if paramLength > maxLength {
				maxLength = paramLength
			}

			details = append(details, [2]string{paramName, parameter.Help})
		}

		separator := "\t"
		for _, detailArg := range details {
			spaces := ""
			for i := 0; i < maxLength-len(detailArg[0])+8; i++ {
				spaces = spaces + " "
			}

			detail = fmt.Sprintf("%s%s%s%s%s", detail, separator, detailArg[0], spaces, detailArg[1])
			separator = "\n\t"
		}
	}
	return detail
}

func (h *Help) Generate() string {
	usageLine := fmt.Sprintf("Usage: %s %s", h.Name, h.Usage())

	usage := fmt.Sprintf("%s\n\n",
		usageLine,
	)

	if h.Desc != "" {
		usage = fmt.Sprintf("%s%s\n\n", usage, h.Desc)
	}

	details := h.Details()
	if details != "" {
		usage = fmt.Sprintf("%s%s\n", usage, details)
	}

	return usage
}
