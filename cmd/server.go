package main

import (
	"amazonApi/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Amazon struct {
	Name    string `json:"name" validate:"gt=4,lt=100"`
	Maker   string `json:"maker" validate:"gt=0,lt=100"`
	Price   int    `json:"price" validate:"min=1,max=9999999999"`
	Comment string `json:"comment" validate:"lte=256"`
	Url     string `json:"url" validate:"url"`
	Asin    string `json:"asin" validate:"len=10,alphanum"`
}

type AmazonPatch struct {
	Name    *string `json:"name" validate:"omitempty,gt=4,lt=100"`
	Maker   *string `json:"maker" validate:"omitempty,gt=0,lt=100"`
	Price   *int    `json:"price" validate:"omitempty,min=1,max=9999999999"`
	Comment *string `json:"comment" validate:"omitempty,lte=256"`
	Url     *string `json:"url" validate:"omitempty,url"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Validator = NewValidator()

	//mysql connection
	dsn := "docker:docker@tcp(127.0.0.1:3306)/amazonApi?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	// Migrate the schema
	if err := db.AutoMigrate(&repository.AmazonItems{}); err != nil {
		panic(err.Error())
	}

	e.POST("/amazon", func(context echo.Context) error {
		// リクエストを取得する
		req := new(Amazon)
		_ = context.Bind(req)
		// バリデーション
		if err := context.Validate(req); err != nil {
			return context.JSON(http.StatusBadRequest, err)
		}
		// Create
		now := time.Now()
		tx := db.Create(&repository.AmazonItems{
			Asin:      req.Asin,
			Name:      req.Name,
			Maker:     req.Maker,
			Price:     req.Price,
			Url:       req.Url,
			Comment:   req.Comment,
			CreatedAt: now,
			UpdatedAt: now,
		})
		if tx.Error != nil {
			return context.JSON(http.StatusBadRequest, tx.Error)
		}
		return context.JSON(http.StatusCreated, req)
	})
	e.PATCH("/amazon/inactive/:asin", func(context echo.Context) error {
		// リクエストを取得する
		asin := context.Param("asin")
		m := new(repository.AmazonItems)

		m.IsDelete = repository.DELETE
		if tx := db.Model(m).Where("asin = ?", asin).Updates(m); tx.Error != nil {
			return context.JSON(http.StatusBadRequest, tx.Error)
		}

		return context.JSON(http.StatusNoContent, nil)
	})
	e.PATCH("/amazon/active/:asin", func(context echo.Context) error {
		// リクエストを取得する
		asin := context.Param("asin")
		m := new(repository.AmazonItems)
		if tx := db.Unscoped().Model(m).Where("asin = ?", asin).Update("is_delete", repository.NOT_DELETE); tx.Error != nil {
			return context.JSON(http.StatusBadRequest, tx.Error)
		}

		return context.JSON(http.StatusNoContent, nil)
	})
	e.GET("/amazon/:asin", func(context echo.Context) error {
		asin := context.Param("asin")
		m := new(repository.AmazonItems)
		if tx := db.First(m, "asin = ?", asin); tx.Error != nil {
			return context.JSON(http.StatusNotFound, tx.Error)
		}
		res := &Amazon{
			Name:    m.Name,
			Maker:   m.Maker,
			Price:   m.Price,
			Comment: m.Comment,
			Url:     m.Url,
			Asin:    m.Asin,
		}
		return context.JSON(http.StatusOK, res)
	})
	e.PUT("/amazon/:asin", func(context echo.Context) error {
		asin := context.Param("asin")
		req := new(Amazon)
		_ = context.Bind(req)
		req.Asin = asin
		if err := context.Validate(req); err != nil {
			return context.JSON(http.StatusBadRequest, err)
		}

		now := time.Now()
		m := &repository.AmazonItems{
			Name:      req.Name,
			Maker:     req.Maker,
			Price:     req.Price,
			Comment:   req.Comment,
			Url:       req.Url,
			UpdatedAt: now,
		}
		tx := db.Model(m).
			Where("asin = ?", asin).
			Updates(m)
		if tx.Error != nil {
			return context.JSON(http.StatusBadRequest, err)
		}
		return context.JSON(http.StatusOK, req)
	})
	e.PATCH("/amazon/:asin", func(context echo.Context) error {
		asin := context.Param("asin")
		req := new(AmazonPatch)
		_ = context.Bind(req)
		if err := context.Validate(req); err != nil {
			return context.JSON(http.StatusBadRequest, err)
		}

		m := new(repository.AmazonItems)
		if tx := db.Model(m).Find(m, "asin = ?", asin); tx.Error != nil {
			return context.JSON(http.StatusBadRequest, tx.Error)
		}

		if req.Name != nil {
			m.Name = *req.Name
		}
		if req.Maker != nil {
			m.Maker = *req.Maker
		}
		if req.Price != nil {
			m.Price = *req.Price
		}
		if req.Comment != nil {
			m.Comment = *req.Comment
		}
		if req.Url != nil {
			m.Url = *req.Url
		}

		tx := db.Model(m).Where("asin = ?", asin).Updates(*m)
		if tx.Error != nil {
			return context.JSON(http.StatusBadRequest, err)
		}

		return context.JSON(http.StatusOK, &Amazon{
			Name:    m.Name,
			Maker:   m.Maker,
			Price:   m.Price,
			Comment: m.Comment,
			Url:     m.Url,
			Asin:    m.Asin,
		})
	})
	e.DELETE("/amazon/:asin", func(context echo.Context) error {
		asin := context.Param("asin")
		m := new(repository.AmazonItems)
		if tx := db.Unscoped().First(m, "asin = ?", asin); tx.Error != nil {
			return context.JSON(http.StatusNotFound, tx.Error)
		}
		if tx := db.Unscoped().Delete(m, "asin = ?", asin); tx.Error != nil {
			return context.JSON(http.StatusBadRequest, tx.Error)
		}
		return context.JSON(http.StatusNoContent, nil)
	})

	e.Logger.Fatal(e.Start(":1232"))
}
