package util

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/sefikcan/address-api/pkg/logger"
)

func GetRequestId(c *fiber.Ctx) string {
	return c.Get(fiber.HeaderXRequestID)
}

func GetIPAddress(c *fiber.Ctx) string {
	return c.IP()
}

func GetRequestCtx(c *fiber.Ctx) context.Context {
	return context.WithValue(c.UserContext(), "RequestCtx", GetRequestId(c))
}

func PrepareLogging(c *fiber.Ctx, logger logger.Logger, err error) {
	logger.Errorf("Error, RequestId: %s, IPAddress: %s, Error: %s", GetRequestId(c), GetIPAddress(c), err)
}
