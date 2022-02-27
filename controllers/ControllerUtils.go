package controllers

import (
	"errors"
	"fmt"
	"net/url"
)

func parseStrings(urlValues url.Values, paramNames ...string) (map[string]string, error) {
	var parsedParams map[string]string
	for _, paramName := range paramNames {
		if paramValues, ok := urlValues[paramName]; ok {
			if len(paramValues) == 1 {
				parsedParams[paramName] = paramValues[0]
				continue
			}
		}
		msg := fmt.Sprintf("cannot parse request parameter: %s", paramName)
		return nil, errors.New(msg)
	}
	return parsedParams, nil
}
