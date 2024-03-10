package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func Default(c echo.Context) error {
	return c.JSON(http.StatusOK, "webhook router 2.0")
}

func Bot(c echo.Context) error {
	requestBody := make(map[string]interface{})
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	params := strings.Split(c.QueryParams()["params"][0], ",")
	paramsMap := make(map[string]interface{})

	for p := range params {
		kv := strings.Split(params[p], "=")
		paramsMap[kv[0]] = kv[1]
	}

	responseLink := ""

	if _, ok := paramsMap["route_ip"]; ok {
		if _, ok := paramsMap["route_port"]; ok {
			responseLink += fmt.Sprintf("http://%s:%s", paramsMap["route_ip"], paramsMap["route_port"])
			delete(paramsMap, "route_ip")
			delete(paramsMap, "route_port")
		}
	}

	if len(responseLink) == 0 {
		return c.JSON(http.StatusBadRequest, "Missing route_ip or route_port params")
	}

	if len(paramsMap) > 0 {
		i := 0
		for key, value := range paramsMap {
			if i == 0 {
				responseLink += fmt.Sprintf("/?%s=%s", key, value)
			} else {
				responseLink += fmt.Sprintf("&%s=%s", key, value)
			}
			i++
		}
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response, err := http.Post(responseLink, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusServiceUnavailable, err.Error())
	}
	defer response.Body.Close()

	return c.JSON(response.StatusCode, response.Body)
}

func SetWebhook(c echo.Context) error {
	routerUrl := os.Getenv("THIS_SERVER_HTTPS_ADDRESS")
	var requestBody struct {
		TelegramToken      string                 `json:"telegram_token" validate:"required"`
		RouteIP            string                 `json:"route_ip" validate:"required,ipv4"`
		RoutePort          int                    `json:"route_port" validate:"required"`
		MaxConnections     int                    `json:"max_connections" validate:"required"`
		DropPendingUpdates bool                   `json:"drop_pending_updates" validate:"required"`
		ExtraParams        map[string]interface{} `json:"extra_params"`
	}
	var responseBody json.RawMessage

	err := json.NewDecoder(c.Request().Body).Decode(&requestBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	link := "https://api.telegram.org/bot" + requestBody.TelegramToken + "/setWebhook" +
		"?max_connections=" + strconv.Itoa(requestBody.MaxConnections) +
		"&drop_pending_updates=" + strconv.FormatBool(requestBody.DropPendingUpdates) +
		"&url=" + routerUrl + "/bot?params=route_ip=" + requestBody.RouteIP + ",route_port=" + strconv.Itoa(requestBody.RoutePort)

	for key, value := range requestBody.ExtraParams {
		link += fmt.Sprintf(",%s=%v", key, value)
	}

	res, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	res.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(res)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, err.Error())
	}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, responseBody)
}

func DeleteWebhook(c echo.Context) error {
	var requestBody struct {
		TelegramToken string `json:"telegram_token" validate:"required"`
	}
	var responseBody json.RawMessage

	err := json.NewDecoder(c.Request().Body).Decode(&requestBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := http.NewRequest(
		"GET",
		"https://api.telegram.org/bot"+requestBody.TelegramToken+"/deleteWebhook",
		nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	res.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(res)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, err.Error())
	}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.StatusCode, responseBody)
}
