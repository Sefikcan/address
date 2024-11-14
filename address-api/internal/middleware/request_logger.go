package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sefikcan/address-api/pkg/util"
	"time"
)

func (mw Manager) RequestLogger(c *fiber.Ctx) error {
	start := time.Now()

	// proceed to the next middleware or handler
	err := c.Next()

	// log request details after handler execution
	req := c.Request()
	res := c.Response()
	status := res.StatusCode()
	size := res.Header.ContentLength()
	elapsed := time.Since(start).String()
	requestId := util.GetRequestId(c)

	mw.logger.Infof("RequestId: %s, Method: %s, Url: %s, Status: %v, Size: %v, Time: %s",
		requestId, req.Header.Method(), req.URI().String(), status, size, elapsed)

	return err
}
