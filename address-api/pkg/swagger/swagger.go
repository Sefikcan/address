package swagger

import (
	"encoding/json"
	"fmt"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/pkg/errors"
	"github.com/sefikcan/address/address-api/pkg/logger"
	"os"
	"path"
)

type Swagger struct {
	FileName string
	BasePath string
	logger   logger.Logger
}

func NewSwagger(fileName, basePath string, logger logger.Logger) *Swagger {
	return &Swagger{
		FileName: fileName,
		BasePath: basePath,
		logger:   logger,
	}
}

func (s *Swagger) UseSwaggerUI(op middleware.SwaggerUIOpts) fiber.Handler {
	return adaptor.HTTPHandler(middleware.SwaggerUI(op, nil))
}

func (s *Swagger) UseSwaggerFile() (fiber.Handler, error) {
	fullPath := path.Join(s.BasePath, s.FileName)
	s.logger.Infof("Looking for Swagger file at: %s", fullPath)

	if _, err := os.Stat(s.FileName); os.IsNotExist(err) {
		s.logger.Errorf("Swagger file not found: %s", s.FileName)
		return nil, errors.New(fmt.Sprintf("%s file is not exist", s.FileName))
	}

	doc, err := loads.Spec(s.FileName)
	if err != nil {
		s.logger.Errorf("Error loading swagger spec: %v", err)
		return nil, err
	}

	b, err := json.MarshalIndent(doc.Spec(), "", " ")
	if err != nil {
		s.logger.Errorf("Error marshalling swagger spec: %v", err)
		return nil, err
	}

	return func(ctx *fiber.Ctx) error {
		ctx.Response().Header.Set("Content-Type", "application/json")
		return ctx.Send(b)
	}, nil
}

func (s *Swagger) RegisterSwagger(app *fiber.App) {
	swaggerUI := middleware.SwaggerUIOpts{
		BasePath: s.BasePath,
		Path:     "v1/docs",
		SpecURL:  path.Join(s.BasePath, "v1/swagger.json"),
	}

	swaggerUIHandler := s.UseSwaggerUI(swaggerUI)
	swaggerFileHandler, err := s.UseSwaggerFile()
	if err != nil {
		s.logger.Fatalf("Failed to initialize swagger handlers: %v", err)
		panic(err)
	}

	swaggerV2UI := middleware.SwaggerUIOpts{
		BasePath: s.BasePath,
		Path:     "v2/docs",
		SpecURL:  path.Join(s.BasePath, "v2/swagger.json"),
	}

	swaggerUIV2Handler := s.UseSwaggerUI(swaggerV2UI)
	if err != nil {
		s.logger.Fatalf("Failed to initialize swagger v2 handlers: %v", err)
		panic(err)
	}

	app.Use(path.Join(s.BasePath, "v1/docs"), swaggerUIHandler)
	app.Use(path.Join(s.BasePath, "v1/swagger.json"), swaggerFileHandler)
	app.Use(path.Join(s.BasePath, "v2/docs"), swaggerUIV2Handler)
	app.Use(path.Join(s.BasePath, "v2/swagger.json"), swaggerFileHandler)
}
