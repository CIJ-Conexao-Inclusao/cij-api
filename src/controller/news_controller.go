package controller

import (
	"cij_api/src/model"
	"cij_api/src/service"
	"mime/multipart"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type NewsController struct {
	newsService service.NewsService
}

func NewNewsController(newsService service.NewsService) *NewsController {
	return &NewsController{
		newsService: newsService,
	}
}

// ListNews
// @Summary List all registered news.
// @Description list all registered news.
// @Tags News
// @Accept application/json
// @Produce json
// @Success 200 {array} model.NewsResponse
// @Failure 404 {object} string "not found"
// @Failure 500 {object} string "internal server error"
// @Router /news [get]
func (n *NewsController) ListNews(ctx *fiber.Ctx) error {
	var response model.Response

	news, err := n.newsService.ListNews()
	if err.Code != "" {
		response = model.Response{
			Message: err.Error(),
			Code:    err.Code,
		}

		return ctx.Status(http.StatusInternalServerError).JSON(response)
	}

	responseData := model.ResponseData[[]model.NewsResponse]{
		Message: "success",
		Data:    news,
	}

	return ctx.Status(http.StatusOK).JSON(responseData)
}

// CreateNews
// @Summary Create a new news.
// @Description create a new news.
// @Tags News
// @Accept json
// @Produce json
// @Param news formData model.NewsRequest true "news"
// @Param banner formData file true "banner"
// @Param authorImage formData file true "author_image"
// @Success 201 {object} model.Response
// @Failure 400 {object} string "bad request"
// @Failure 500 {object} string "internal server error"
// @Router /news [post]
func (n *NewsController) CreateNews(ctx *fiber.Ctx) error {
	var response model.Response

	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(model.Response{
			Message: err.Error(),
		})
	}

	request := model.NewsRequest{
		Title:       form.Value["title"][0],
		Description: form.Value["description"][0],
		Author:      form.Value["author"][0],
		Date:        form.Value["date"][0],
	}

	files := make(map[string]multipart.FileHeader)

	for filename, file := range form.File {
		files[filename] = *file[0]
	}

	newsError := n.newsService.CreateNews(request, files)
	if newsError.Code != "" {
		return ctx.Status(http.StatusInternalServerError).JSON(model.Response{
			Message: newsError.Error(),
			Code:    newsError.Code,
		})
	}

	response = model.Response{
		Message: "success",
	}

	return ctx.Status(http.StatusCreated).JSON(response)
}
