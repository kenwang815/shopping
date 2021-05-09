package daos_test

import (
	"fmt"
	"os"
	"path/filepath"
	"shopping/config"
	"shopping/daos"
	"shopping/database"
	"shopping/model"
	"shopping/model/commodity"
	"shopping/test"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

type CommodityTestCaseSuite struct {
	db            database.IDatabase
	commodityRepo commodity.Repository
}

func setupCommodityTestCaseSuite(t *testing.T) (CommodityTestCaseSuite, func(t *testing.T)) {
	s := CommodityTestCaseSuite{}

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
	s.commodityRepo = daos.NewCommodityRepo(s.db.GetDB())

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

func GetCommodity3() *commodity.Commodity {
	return &commodity.Commodity{
		Id:          1,
		Name:        "apple",
		Cost:        30,
		Price:       50,
		Description: "test",
		Sell:        true,
		StartTime:   time.Date(2021, time.Month(5), 9, 10, 10, 30, 0, time.UTC),
		EndTime:     time.Date(2021, time.Month(5), 10, 10, 10, 30, 0, time.UTC),
	}
}

func UpdateCommodity1() *commodity.Commodity {
	return &commodity.Commodity{
		Id:   1,
		Sell: true,
	}
}

func UpdateCommodity2() *commodity.Commodity {
	return &commodity.Commodity{
		Id: 2,
	}
}

func TestCommodityDaos_Create(t *testing.T) {
	s, teardownTestCase := setupCommodityTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		name          string
		testData      *commodity.Commodity
		wantResult    *commodity.Commodity
		err           error
		setupTestCase test.SetupSubTest
	}{
		{
			name:       "success",
			testData:   GetCommodity1(),
			wantResult: GetCommodity1(),
			err:        nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&commodity.Commodity{})
				s.db.GetDB().AutoMigrate(&commodity.Commodity{})

				return func(t *testing.T) {
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)

			p, err := s.commodityRepo.Create(tc.testData)
			if err != nil {
				assert.EqualError(t, err, tc.err.Error(), "An error was expected")
			} else {
				assert.Equal(t, p, tc.wantResult)
			}
		})
	}
}

func TestCommodityDaos_Update(t *testing.T) {
	s, teardownTestCase := setupCommodityTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		name          string
		testData      *commodity.Commodity
		wantResult    *commodity.Commodity
		rowAffected   int64
		err           error
		setupTestCase test.SetupSubTest
	}{
		{
			name:        "success",
			testData:    UpdateCommodity1(),
			wantResult:  GetCommodity3(),
			rowAffected: 1,
			err:         nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&commodity.Commodity{})
				s.db.GetDB().AutoMigrate(&commodity.Commodity{})
				s.db.GetDB().Create(GetCommodity1())

				return func(t *testing.T) {
				}
			},
		},
		{
			name:        "ignore_id_update",
			testData:    UpdateCommodity2(),
			wantResult:  nil,
			rowAffected: 0,
			err:         nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&commodity.Commodity{})
				s.db.GetDB().AutoMigrate(&commodity.Commodity{})
				s.db.GetDB().Create(GetCommodity1())

				return func(t *testing.T) {
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)

			tmp, affected, err := s.commodityRepo.Update(tc.testData)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.rowAffected, affected)
			if err == nil && tmp != nil {
				d, _ := s.commodityRepo.Find(&commodity.Commodity{Id: tc.testData.Id}, &model.Page{Limit: 1, Offset: 0})
				assert.Equal(t, tc.wantResult, d[0])
			}
		})
	}
}

func TestCommodityDaos_Delete(t *testing.T) {
	s, teardownTestCase := setupCommodityTestCaseSuite(t)
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
			id:          GetCommodity1().Id,
			rowAffected: 1,
			err:         nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&commodity.Commodity{})
				s.db.GetDB().AutoMigrate(&commodity.Commodity{})
				s.db.GetDB().Create(GetCommodity1())

				return func(t *testing.T) {
				}
			},
		},
		{
			name:        "not exist id",
			id:          GetCommodity2().Id,
			rowAffected: 0,
			err:         nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&commodity.Commodity{})
				s.db.GetDB().AutoMigrate(&commodity.Commodity{})
				s.db.GetDB().Create(GetCommodity1())

				return func(t *testing.T) {
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)

			affected, err := s.commodityRepo.Delete(tc.id)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.rowAffected, affected)
		})
	}
}

func TestCommodityDaos_Find(t *testing.T) {
	s, teardownTestCase := setupCommodityTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		name          string
		testData      *commodity.Commodity
		testPage      *model.Page
		wantResult    []*commodity.Commodity
		err           error
		setupTestCase test.SetupSubTest
	}{
		{
			name:       "no data",
			testData:   &commodity.Commodity{},
			testPage:   &model.Page{Limit: 0, Offset: 0},
			wantResult: []*commodity.Commodity{},
			err:        nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&commodity.Commodity{})
				s.db.GetDB().AutoMigrate(&commodity.Commodity{})
				s.db.GetDB().Create(GetCommodity1())
				s.db.GetDB().Create(GetCommodity2())

				return func(t *testing.T) {
				}
			},
		},
		{
			name:       "input limit > count",
			testData:   &commodity.Commodity{},
			testPage:   &model.Page{Limit: 3, Offset: 0},
			wantResult: []*commodity.Commodity{GetCommodity1(), GetCommodity2()},
			err:        nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&commodity.Commodity{})
				s.db.GetDB().AutoMigrate(&commodity.Commodity{})
				s.db.GetDB().Create(GetCommodity1())
				s.db.GetDB().Create(GetCommodity2())

				return func(t *testing.T) {
				}
			},
		},
		{
			name:       "find id",
			testData:   &commodity.Commodity{Id: GetCommodity2().Id},
			testPage:   &model.Page{Limit: 2, Offset: 0},
			wantResult: []*commodity.Commodity{GetCommodity2()},
			err:        nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)

			commoditys, err := s.commodityRepo.Find(tc.testData, tc.testPage)
			if err != nil {
				assert.EqualError(t, err, tc.err.Error(), "An error was expected")
			} else {
				assert.Equal(t, commoditys, tc.wantResult)
			}
		})
	}
}
