package daos_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"

	"shopping/config"
	"shopping/daos"
	"shopping/database"
	"shopping/model"
	"shopping/model/catalog"
	"shopping/test"
)

type CatalogTestCaseSuite struct {
	db          database.IDatabase
	catalogRepo catalog.Repository
}

func setupCatalogTestCaseSuite(t *testing.T) (CatalogTestCaseSuite, func(t *testing.T)) {
	s := CatalogTestCaseSuite{}

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
	s.catalogRepo = daos.NewCatalogRepo(s.db.GetDB())
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
		Hide: false,
	}
}

func UpdateCatalog1() *catalog.Catalog {
	return &catalog.Catalog{
		Id:   1,
		Name: "food",
		Hide: true,
	}
}

func UpdateCatalog2() *catalog.Catalog {
	return &catalog.Catalog{
		Id: 1,
	}
}

func TestCatalogDaos_Create(t *testing.T) {
	s, teardownTestCase := setupCatalogTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		name          string
		testData      *catalog.Catalog
		wantResult    *catalog.Catalog
		err           error
		setupTestCase test.SetupSubTest
	}{
		{
			name:       "success",
			testData:   GetCatalog1(),
			wantResult: GetCatalog1(),
			err:        nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&catalog.Catalog{})
				s.db.GetDB().AutoMigrate(&catalog.Catalog{})

				return func(t *testing.T) {
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)

			p, err := s.catalogRepo.Create(tc.testData)
			if err != nil {
				assert.EqualError(t, err, tc.err.Error(), "An error was expected")
			} else {
				assert.Equal(t, p, tc.wantResult)
			}
		})
	}
}

func TestCatalogDaos_Update(t *testing.T) {
	s, teardownTestCase := setupCatalogTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		name          string
		testData      *catalog.Catalog
		wantResult    *catalog.Catalog
		rowAffected   int64
		err           error
		setupTestCase test.SetupSubTest
	}{
		{
			name:        "success",
			testData:    UpdateCatalog1(),
			wantResult:  UpdateCatalog1(),
			rowAffected: 1,
			err:         nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&catalog.Catalog{})
				s.db.GetDB().AutoMigrate(&catalog.Catalog{})
				s.db.GetDB().Create(GetCatalog1())

				return func(t *testing.T) {
				}
			},
		},
		{
			name:        "ignore_id_update",
			testData:    UpdateCatalog2(),
			wantResult:  nil,
			rowAffected: 0,
			err:         nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&catalog.Catalog{})
				s.db.GetDB().AutoMigrate(&catalog.Catalog{})
				s.db.GetDB().Create(GetCatalog1())

				return func(t *testing.T) {
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)

			tmp, affected, err := s.catalogRepo.Update(tc.testData)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.rowAffected, affected)
			if err == nil && tmp != nil {
				d, _ := s.catalogRepo.Find(&catalog.Catalog{Id: tc.testData.Id}, &model.Page{Limit: 1, Offset: 0})
				assert.Equal(t, tc.testData, d[0])
			}
		})
	}
}

func TestCatalogDaos_Delete(t *testing.T) {
	s, teardownTestCase := setupCatalogTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		name          string
		id            int
		rowAffected   int64
		err           error
		setupTestCase test.SetupSubTest
	}{
		{
			name:        "success",
			id:          GetCatalog1().Id,
			rowAffected: 1,
			err:         nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&catalog.Catalog{})
				s.db.GetDB().AutoMigrate(&catalog.Catalog{})
				s.db.GetDB().Create(GetCatalog1())

				return func(t *testing.T) {
				}
			},
		},
		{
			name:        "not exist id",
			id:          GetCatalog2().Id,
			rowAffected: 0,
			err:         nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&catalog.Catalog{})
				s.db.GetDB().AutoMigrate(&catalog.Catalog{})
				s.db.GetDB().Create(GetCatalog1())

				return func(t *testing.T) {
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)

			affected, err := s.catalogRepo.Delete(tc.id)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.rowAffected, affected)
		})
	}
}

func TestCatalogDaos_Find(t *testing.T) {
	s, teardownTestCase := setupCatalogTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		name          string
		testData      *catalog.Catalog
		testPage      *model.Page
		wantResult    []*catalog.Catalog
		err           error
		setupTestCase test.SetupSubTest
	}{
		{
			name:       "no data",
			testData:   &catalog.Catalog{},
			testPage:   &model.Page{Limit: 0, Offset: 0},
			wantResult: []*catalog.Catalog{},
			err:        nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&catalog.Catalog{})
				s.db.GetDB().AutoMigrate(&catalog.Catalog{})
				s.db.GetDB().Create(GetCatalog1())
				s.db.GetDB().Create(GetCatalog2())

				return func(t *testing.T) {
				}
			},
		},
		{
			name:       "input limit > count",
			testData:   &catalog.Catalog{},
			testPage:   &model.Page{Limit: 3, Offset: 0},
			wantResult: []*catalog.Catalog{GetCatalog1(), GetCatalog2()},
			err:        nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&catalog.Catalog{})
				s.db.GetDB().AutoMigrate(&catalog.Catalog{})
				s.db.GetDB().Create(GetCatalog1())
				s.db.GetDB().Create(GetCatalog2())

				return func(t *testing.T) {
				}
			},
		},
		{
			name:       "find id",
			testData:   &catalog.Catalog{Id: GetCatalog2().Id},
			testPage:   &model.Page{Limit: 2, Offset: 0},
			wantResult: []*catalog.Catalog{GetCatalog2()},
			err:        nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)

			catalogs, err := s.catalogRepo.Find(tc.testData, tc.testPage)
			if err != nil {
				assert.EqualError(t, err, tc.err.Error(), "An error was expected")
			} else {
				assert.Equal(t, catalogs, tc.wantResult)
			}
		})
	}
}
