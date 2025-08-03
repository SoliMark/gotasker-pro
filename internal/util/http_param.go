package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseUintParam(c *gin.Context, name string, out *uint) error {
	param := c.Param(name)
	id64, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		return err
	}
	*out = uint(id64)
	return nil
}
