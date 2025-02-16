package assistants

import (
	cl "avitomaxwin/curloger"

	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Возвращает response клиенту
func ResponseReturner(code int, resp string, err error, c *gin.Context) {
	if err != nil {
		c.JSON(code, gin.H{"error": err.Error()})
	} else {
		c.String(code, resp)
	}
}

// Чтение параметров из тела запроса
func ReadBodyData(c *gin.Context) (int, map[string]string, error) {
	var bodyData map[string]interface{}

	body, err := c.GetRawData()
	if err != nil {
		cl.Log(logrus.ErrorLevel, "error while reading request body", map[string]interface{}{
			"error": err,
		})
		return http.StatusInternalServerError, nil, errors.New("error while processing client data")
	}

	if err := json.Unmarshal(body, &bodyData); err != nil {
		cl.Log(logrus.ErrorLevel, "error while unmarshaling json", map[string]interface{}{
			"error": err,
		})
		return http.StatusInternalServerError, nil, errors.New("error while processing client data")
	}

	result := make(map[string]string)

	for key, value := range bodyData {
		switch v := value.(type) {
		case string:
			result[key] = v
		case float64:
			result[key] = strconv.FormatFloat(v, 'f', -1, 64)
		default:
			cl.Log(logrus.ErrorLevel, "unsupported value type - "+key, map[string]interface{}{
				"error": err,
			})
			return http.StatusBadRequest, nil, errors.New("error while processing client data")
		}
	}

	return http.StatusOK, result, nil
}
