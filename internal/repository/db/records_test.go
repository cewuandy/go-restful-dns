package db

import (
	"context"
	"github.com/samber/do"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"testing"

	"github.com/cewuandy/go-restful-dns/internal/domain"
	pkgGorm "github.com/cewuandy/go-restful-dns/pkg/gorm"
)

type recordRepoTestSuite struct {
	suite.Suite

	repo domain.RecordRepo
}

func TestRecordRepo(t *testing.T) {
	suite.Run(t, &recordRepoTestSuite{})
}

func (t *recordRepoTestSuite) SetupSuite() {
	injector := do.New()
	db, err := gorm.Open(
		sqlite.Open("dns.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
	t.Nil(err)

	do.ProvideValue[*gorm.DB](injector, db)
	err = pkgGorm.AutoMigrate(db)
	t.Nil(err)

	t.repo, _ = NewRecordsRepo(injector)

	_ = t.repo.Create(
		context.Background(), &domain.Record{
			Name:   "test.com.",
			RrType: 1,
			Class:  1,
			Record: "test.com.\t1440\tIN\tA\t1.1.1.1",
		},
	)
}

func (t *recordRepoTestSuite) SetupTest() {

}

func (t *recordRepoTestSuite) TearDownSuite() {
	_ = os.Remove("dns.db")
}

func (t *recordRepoTestSuite) TestCreate() {
	t.Run(
		"success", func() {
			err := t.repo.Create(
				context.Background(), &domain.Record{
					Name:   "test.com.",
					RrType: 1,
					Class:  1,
					Record: "test.com.\t1440\tIN\tA\t1.1.1.1",
				},
			)
			t.Nil(err)
		},
	)
}

func (t *recordRepoTestSuite) TestGet() {
	t.Run(
		"success", func() {
			record, err := t.repo.Get(context.Background(), "test.com.", 1, 1)
			t.Equal("test.com.", record.Name)
			t.Equal(uint16(1), record.RrType)
			t.Equal(uint16(1), record.Class)
			t.Nil(err)
		},
	)
}

func (t *recordRepoTestSuite) TestList() {
	t.Run(
		"success", func() {
			records, err := t.repo.List(context.Background())
			t.NotNil(records)
			t.Equal("test.com.", records[0].Name)
			t.Equal(uint16(1), records[0].RrType)
			t.Equal(uint16(1), records[0].Class)
			t.Nil(err)
		},
	)
}

func (t *recordRepoTestSuite) TestUpdate() {
	t.Run(
		"success", func() {
			err := t.repo.Update(
				context.Background(), &domain.Record{
					Name:   "test.com.",
					RrType: 1,
					Class:  1,
					Record: "test.com.\t1440\tIN\tA\t2.2.2.2",
				},
			)
			t.Nil(err)
		},
	)
}

func (t *recordRepoTestSuite) TestDelete() {
	t.Run(
		"success", func() {
			err := t.repo.Delete(context.Background(), "test.com.", 1, 1)
			t.Nil(err)
		},
	)
}
