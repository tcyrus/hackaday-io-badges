// +build openwhisk

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
)

// Main function for the action
func Main(obj map[string]interface{}) map[string]interface{} {
	var badgeId int

	_, err := fmt.Sscanf(obj["__ow_path"].(string), "/%d", &badgeId)
	if err != nil {
		return map[string]interface{}{
			"statusCode": http.StatusInternalServerError,
			"body": err.Error(),
		}
	}

	data, err := getProject(strconv.Itoa(badgeId))
	if err != nil {
		return map[string]interface{}{
			"statusCode": http.StatusInternalServerError,
			"body": err.Error(),
		}
	}

	tmp_data := &BadgeData{Skulls: int(data["skulls"].(float64)),Name: data["name"].(string)}

	var tpl bytes.Buffer
	if err := Badge.Execute(&tpl, tmp_data); err != nil {
		return map[string]interface{}{
			"statusCode": http.StatusInternalServerError,
			"body": err.Error(),
		}
	}

	return map[string]interface{}{
		"headers": map[string]interface{}{
			"Content-Type": "image/svg+xml",
		},
		"statusCode": http.StatusOK,
		"body": tpl.String(),
	}
}
