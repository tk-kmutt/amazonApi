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
	Code *string `json:"code"`
	Name *string `json:"name"`
	Age  *int    `json:"age" validate:"max=150"`
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
			return context.String(http.StatusBadRequest, err.Error())
		}
		// Create
		now := time.Now()
		db.Create(&repository.AmazonItems{
			Asin:      req.Asin,
			Name:      req.Name,
			Maker:     req.Maker,
			Price:     req.Price,
			Url:       req.Url,
			Comment:   req.Comment,
			CreatedAt: now,
			UpdatedAt: now,
			IsDelete:  false,
		})
		return context.JSON(http.StatusCreated, req)
	})
	e.PATCH("/amazon/inactive/:asin", func(context echo.Context) error {
		// リクエストを取得する
		asin := context.Param("asin")
		m := new(repository.AmazonItems)
		// First
		if tx := db.First(m, "asin = ?", asin); tx.Error != nil {
			return context.String(http.StatusNotFound, tx.Error.Error())
		}

		m.IsDelete = true
		if tx := db.Model(m).Where("asin = ?", asin).Updates(m); tx.Error != nil {
			return context.String(http.StatusBadRequest, tx.Error.Error())
		}

		return context.JSON(http.StatusAccepted, nil)
	})
	e.PATCH("/amazon/active/:asin", func(context echo.Context) error {
		// リクエストを取得する
		asin := context.Param("asin")
		m := new(repository.AmazonItems)
		// First
		if tx := db.First(m, "asin = ?", asin); tx.Error != nil {
			return context.String(http.StatusNotFound, tx.Error.Error())
		}

		if tx := db.Model(m).Where("asin = ?", asin).Update("is_delete", false); tx.Error != nil {
			return context.String(http.StatusBadRequest, tx.Error.Error())
		}

		return context.JSON(http.StatusAccepted, nil)
	})
	e.GET("/simple/:code", func(context echo.Context) error {
		// リクエストを取得する
		code := context.Param("code")
		m := new(repository.AmazonItems)
		if tx := db.First(m, "code = ?", code); tx.Error != nil {
			return context.String(http.StatusNotFound, tx.Error.Error())
		}

		user := &Amazon{
			//Code: m.Code,
			//Name: m.Name,
			//Age:  m.Age,
		}
		return context.JSON(http.StatusOK, user)
	})
	e.PUT("/simple/:code", func(context echo.Context) error {
		// リクエストを取得する
		code := context.Param("code")
		user := new(Amazon)
		// バリデーション
		_ = context.Bind(user)
		if err := context.Validate(user); err != nil {
			return context.JSON(http.StatusBadRequest, err)
		}

		m := new(repository.AmazonItems)
		// First
		if tx := db.First(m, "code = ?", code); tx.Error != nil {
			return context.String(http.StatusNotFound, tx.Error.Error())
		}
		// Update
		now := time.Now()
		db.Model(m).
			Where("code = ?", code).
			Updates(repository.AmazonItems{
				Name: user.Name,
				//Age:       user.Age,
				UpdatedAt: now,
			})
		return context.JSON(http.StatusOK, user)
	})
	e.PATCH("/simple/:code", func(context echo.Context) error {
		// リクエストを取得する
		code := context.Param("code")
		user := new(AmazonPatch)
		_ = context.Bind(user)
		// バリデーション
		if err := context.Validate(user); err != nil {
			return context.JSON(http.StatusBadRequest, err)
		}

		m := new(repository.AmazonItems)
		// First
		if tx := db.First(m, "code = ?", code); tx.Error != nil {
			return context.String(http.StatusNotFound, tx.Error.Error())
		}

		tx := db.Model(m).Where("code = ?", code)
		if user.Age != nil {
			//m.Age = *user.Age
		}
		if user.Name != nil {
			m.Name = *user.Name
		}
		tx.Updates(*m)
		return context.JSON(http.StatusOK, &Amazon{
			//Code: m.Code,
			//Name: m.Name,
			//Age:  m.Age,
		})
	})
	e.DELETE("/simple/:code", func(context echo.Context) error {
		// リクエストを取得する
		code := context.Param("code")

		m := new(repository.AmazonItems)
		// First
		if tx := db.First(m, "code = ?", code); tx.Error != nil {
			return context.String(http.StatusNotFound, tx.Error.Error())
		}
		db.Delete(m, "code = ?", code)

		return context.String(http.StatusNoContent, "")
	})
	e.Logger.Fatal(e.Start(":1232"))
}
