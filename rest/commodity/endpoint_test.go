package commodity_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"

	"shopping/config"
	"shopping/daos"
	"shopping/database"
	"shopping/model/commodity"
	routeCommodity "shopping/rest/commodity"
	"shopping/service"
	"shopping/test"
)

type CommodityTestCaseSuite struct {
	db database.IDatabase
	c  *gin.Engine
}

func setupCommodityTestCaseSuite(t *testing.T) (CommodityTestCaseSuite, func(t *testing.T)) {
	s := CommodityTestCaseSuite{
		c: gin.New(),
	}
	s.c.Use(gin.Recovery())

	name := filepath.Join(os.TempDir(), "gorm"+uuid.Must(uuid.NewV4()).String()+".db")
	df, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if df == nil || err != nil {
		panic(fmt.Sprintf("No error should happen when creating db file, but got %+v", err))
	}

	c := &config.Database{
		Dialect: "sqlite",
		Host:    df.Name(),
	}

	s.db, err = database.NewDatabase(c)
	commodityRepo := daos.NewCommodityRepo(s.db.GetDB())
	service.CommodityService = service.NewCommodityService(commodityRepo)

	routeCommodity.MakeHandler(s.c.Group("/v1"))

	return s, func(t *testing.T) {
		s.db.Close()
		os.Remove(df.Name())
	}
}

func GetCommodity1() *commodity.Commodity {
	return &commodity.Commodity{
		Id:          1,
		Name:        "apple",
		Cost:        30,
		Price:       50,
		Description: "test",
		Sell:        false,
		StartTime:   time.Date(2021, time.Month(5), 9, 10, 10, 30, 0, time.UTC),
		EndTime:     time.Date(2021, time.Month(5), 10, 10, 10, 30, 0, time.UTC),
	}
}

func GetCommodity2() *commodity.Commodity {
	return &commodity.Commodity{
		Id:          2,
		Name:        "banana",
		Cost:        10,
		Price:       20,
		Description: "test",
		Sell:        true,
		StartTime:   time.Date(2021, time.Month(5), 8, 10, 10, 30, 0, time.UTC),
		EndTime:     time.Date(2021, time.Month(5), 11, 10, 10, 30, 0, time.UTC),
	}
}

func TestCommodityFindHandler(t *testing.T) {
	s, teardownTestCase := setupCommodityTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		description  string
		route        string
		method       string
		params       map[string]string
		expected     string
		expectedCode int
		setupSubTest test.SetupSubTest
	}{
		{
			description: "input page=1, number=2",
			route:       "/v1/commodity",
			method:      "GET",
			params: map[string]string{
				"page":   "1",
				"number": "2",
			},
			expected:     `{"code":2000000,"data":{"datas":[{"id":1,"catalog_id":0,"name":"apple","cost":30,"price":50,"description":"test","sell":false,"start_time":"2021-05-09T10:10:30Z","end_time":"2021-05-10T10:10:30Z"},{"id":2,"catalog_id":0,"name":"banana","cost":10,"price":20,"description":"test","sell":true,"start_time":"2021-05-08T10:10:30Z","end_time":"2021-05-11T10:10:30Z"}],"page":{"page":1,"number":2}},"msg":"Success"}`,
			expectedCode: http.StatusOK,
			setupSubTest: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&commodity.Commodity{})
				s.db.GetDB().AutoMigrate(&commodity.Commodity{})
				s.db.GetDB().Create(GetCommodity1())
				s.db.GetDB().Create(GetCommodity2())

				return func(t *testing.T) {
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			teardownSubTest := tc.setupSubTest(t)
			defer teardownSubTest(t)

			req := httptest.NewRequest(tc.method, tc.route, nil)
			q := req.URL.Query()
			for k, v := range tc.params {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()

			req.Header.Set("Content-Type", gin.MIMEJSON)
			actul := httptest.NewRecorder()
			s.c.ServeHTTP(actul, req)
			assert.Equal(t, tc.expectedCode, actul.Code)
			assert.Equal(t, tc.expected, strings.Replace(actul.Body.String(), "\n", "", -1))
		})
	}
}

// func TestCommodityRegisterHandler(t *testing.T) {
// 	s, teardownTestCase := setupCommodityTestCaseSuite(t)
// 	defer teardownTestCase(t)

// 	tt := []struct {
// 		description  string
// 		route        string
// 		method       string
// 		body         string
// 		expectedCode int
// 		setupSubTest test.SetupSubTest
// 	}{
// 		{
// 			description:  "success",
// 			route:        "/v1/commodity",
// 			method:       "POST",
// 			body:         `{"name":"apple","cost":30,"price":50,"description":"test","sell":false,"start_time":"2021-05-11 10:10:30","end_time":"2021-05-13 10:10:30"}`,
// 			expectedCode: http.StatusOK,
// 			setupSubTest: func(t *testing.T) func(t *testing.T) {
// 				s.db.GetDB().DropTable(&commodity.Commodity{})
// 				s.db.GetDB().AutoMigrate(&commodity.Commodity{})

// 				return func(t *testing.T) {
// 				}
// 			},
// 		},
// 	}

// 	for _, tc := range tt {
// 		t.Run(tc.description, func(t *testing.T) {
// 			teardownSubTest := tc.setupSubTest(t)
// 			defer teardownSubTest(t)

// 			fmt.Printf("%+v\n", strings.NewReader(tc.body))
// 			req := httptest.NewRequest(tc.method, tc.route, strings.NewReader(tc.body))
// 			req.Header.Set("Content-Type", gin.MIMEJSON)
// 			actul := httptest.NewRecorder()
// 			s.c.ServeHTTP(actul, req)

// 			assert.Equal(t, tc.expectedCode, actul.Code)
// 		})
// 	}
// }

func TestCommodityDeleteHandler(t *testing.T) {
	s, teardownTestCase := setupCommodityTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		description  string
		route        string
		method       string
		expectedCode int
		setupSubTest test.SetupSubTest
	}{
		{
			description:  "success",
			route:        fmt.Sprintf("/v1/commodity/%d", GetCommodity1().Id),
			method:       "DELETE",
			expectedCode: http.StatusOK,
			setupSubTest: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&commodity.Commodity{})
				s.db.GetDB().AutoMigrate(&commodity.Commodity{})
				s.db.GetDB().Create(GetCommodity1())

				return func(t *testing.T) {
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			teardownSubTest := tc.setupSubTest(t)
			defer teardownSubTest(t)

			req := httptest.NewRequest(tc.method, tc.route, nil)
			req.Header.Set("Content-Type", gin.MIMEJSON)
			actul := httptest.NewRecorder()
			s.c.ServeHTTP(actul, req)

			assert.Equal(t, tc.expectedCode, actul.Code)
		})
	}
}

func TestCommodityUpdateHandler(t *testing.T) {
	s, teardownTestCase := setupCommodityTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		description  string
		route        string
		method       string
		body         string
		expectedCode int
		setupSubTest test.SetupSubTest
	}{
		{
			description:  "success",
			route:        "/v1/commodity",
			method:       "PUT",
			body:         `{"id":1,"name":"orange"}`,
			expectedCode: http.StatusOK,
			setupSubTest: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&commodity.Commodity{})
				s.db.GetDB().AutoMigrate(&commodity.Commodity{})
				s.db.GetDB().Create(GetCommodity1())

				return func(t *testing.T) {
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			teardownSubTest := tc.setupSubTest(t)
			defer teardownSubTest(t)

			req := httptest.NewRequest(tc.method, tc.route, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", gin.MIMEJSON)
			actul := httptest.NewRecorder()
			s.c.ServeHTTP(actul, req)
			assert.Equal(t, tc.expectedCode, actul.Code)
		})
	}
}
