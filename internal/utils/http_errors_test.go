package utils

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestSendErrorResponse(t *testing.T) {
	testCases := []struct {
		name           string
		status         int
		message        interface{}
		expectedStatus int
		expectedBody   string
		typeError      string
	}{
		{
			name:           "Error Response 500",
			status:         500,
			message:        "Internal Server Error",
			expectedStatus: 500,
			expectedBody:   `{"error":"Internal Server Error"}`,
			typeError:      "error",
		},
		{
			name:           "Error Response 404",
			status:         404,
			message:        "Not Found",
			expectedStatus: 404,
			expectedBody:   `{"error":"Not Found"}`,
			typeError:      "error",
		},
		{
			name:           "Custom JSON Error Response 500",
			status:         500,
			message:        fiber.Map{"error": "Internal Server Error", "message": "Service User Down"},
			expectedStatus: 500,
			expectedBody:   `{"error":"Internal Server Error", "message": "Service User Down"}`,
			typeError:      "customError",
		},
		{
			name:           "Custom JSON Error Response 404",
			status:         500,
			message:        fiber.Map{"error": "Not Found", "message": "Book Id Not Found"},
			expectedStatus: 500,
			expectedBody:   `{"error": "Not Found", "message": "Book Id Not Found"}`,
			typeError:      "customError",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			app := fiber.New()

			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			defer app.ReleaseCtx(ctx)

			var err error
			if testCase.typeError == "error" {
				err = sendErrorResponse(ctx, testCase.status, testCase.message.(string))
			} else if testCase.typeError == "customError" {
				err = CustomErrorJSON(ctx, testCase.status, testCase.message)
			}

			require.NoError(t, err)

			assert.Equal(t, testCase.expectedStatus, ctx.Response().StatusCode())
			assert.JSONEq(t, testCase.expectedBody, string(ctx.Response().Body()))
		})
	}
}
