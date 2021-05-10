package catalog_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"

	"shopping/config"
	"shopping/daos"
	"shopping/database"
	"shopping/model/catalog"
	routeCatalog "shopping/rest/catalog"
	"shopping/service"
	"shopping/test"
)

type CatalogTestCaseSuite struct {
	db database.IDatabase
	c  *gin.Engine
}

func setupCatalogTestCaseSuite(t *testing.T) (CatalogTestCaseSuite, func(t *testing.T)) {
	s := CatalogTestCaseSuite{
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
	catalogRepo := daos.NewCatalogRepo(s.db.GetDB())
	commodityRepo := daos.NewCommodityRepo(s.db.GetDB())
	service.CatalogService = service.NewCatalogService(catalogRepo, commodityRepo)

	routeCatalog.MakeHandler(s.c.Group("/v1"))

	return s, func(t *testing.T) {
		s.db.Close()
		os.Remove(df.Name())
	}
}

func GetCatalog1() *catalog.Catalog {
	return &catalog.Catalog{
		Id:   1,
		Name: "food",
		Hide: false,
	}
}

func GetCatalog2() *catalog.Catalog {
	return &catalog.Catalog{
		Id:   2,
		Name: "water",
		Hide: true,
	}
}

func TestCatalogFindHandler(t *testing.T) {
	s, teardownTestCase := setupCatalogTestCaseSuite(t)
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
			route:       "/v1/catalog",
			method:      "GET",
			params: map[string]string{
				"page":   "1",
				"number": "2",
			},
			expected:     `{"code":2000000,"data":{"datas":[{"id":1,"name":"food","hide":false},{"id":2,"name":"water","hide":true}],"page":{"page":1,"number":2}},"msg":"Success"}`,
			expectedCode: http.StatusOK,
			setupSubTest: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&catalog.Catalog{})
				s.db.GetDB().AutoMigrate(&catalog.Catalog{})
				s.db.GetDB().Create(GetCatalog1())
				s.db.GetDB().Create(GetCatalog2())

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

func TestCatalogRegisterHandler(t *testing.T) {
	s, teardownTestCase := setupCatalogTestCaseSuite(t)
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
			route:        "/v1/catalog",
			method:       "POST",
			body:         `{"name":"food","hide":false}`,
			expectedCode: http.StatusOK,
			setupSubTest: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&catalog.Catalog{})
				s.db.GetDB().AutoMigrate(&catalog.Catalog{})

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

func TestCatalogDeleteHandler(t *testing.T) {
	s, teardownTestCase := setupCatalogTestCaseSuite(t)
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
			route:        fmt.Sprintf("/v1/catalog/%d", GetCatalog1().Id),
			method:       "DELETE",
			expectedCode: http.StatusOK,
			setupSubTest: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&catalog.Catalog{})
				s.db.GetDB().AutoMigrate(&catalog.Catalog{})
				s.db.GetDB().Create(GetCatalog1())

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

func TestCatalogUpdateHandler(t *testing.T) {
	s, teardownTestCase := setupCatalogTestCaseSuite(t)
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
			route:        "/v1/catalog",
			method:       "PUT",
			body:         `{"id":1,"name":"tool"}`,
			expectedCode: http.StatusOK,
			setupSubTest: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&catalog.Catalog{})
				s.db.GetDB().AutoMigrate(&catalog.Catalog{})
				s.db.GetDB().Create(GetCatalog1())

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
