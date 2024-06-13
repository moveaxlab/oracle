package oracle

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Model struct {
	Id   int    `gorm:"column:ID"`
	Name string `gorm:"column:NAME"`
}

func (Model) TableName() string {
	return "TEST_TABLE"
}

type OracleSuite struct {
	suite.Suite

	db *gorm.DB
}

func (s *OracleSuite) TearDownTest() {
	err := s.db.Exec("DELETE FROM TEST_TABLE WHERE 1 = 1").Error
	s.Require().Nil(err)
}

func (s *OracleSuite) TearDownSuite() {
	err := s.db.Exec("DROP TABLE TEST_TABLE").Error
	s.Require().Nil(err)
}

type Oracle19Suite struct {
	OracleSuite
}

func (s *Oracle19Suite) SetupSuite() {
	url := BuildUrl("oracle-19c", 1521, "MYATP", "ADMIN", "TVDGXvpzQat8", map[string]string{
		"CONNECTION TIMEOUT": "5",
		"LANGUAGE":           "ITALIAN",
		"TERRITORY":          "ITALY",
		"SSL":                "false",
	})

	dialector := New(Config{DSN: url})

	var err error

	s.db, err = gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            false,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	s.Require().Nil(err)

	err = s.db.Exec("CREATE TABLE TEST_TABLE (ID INT PRIMARY KEY, NAME VARCHAR2(64 CHAR) NOT NULL)").Error
	s.Require().Nil(err)
}

type Oracle23Suite struct {
	OracleSuite
}

func (s *Oracle23Suite) SetupSuite() {
	url := BuildUrl("oracle-23c", 1521, "FREE", "SYSTEM", "password", map[string]string{
		"CONNECTION TIMEOUT": "5",
		"LANGUAGE":           "ITALIAN",
		"TERRITORY":          "ITALY",
		"SSL":                "false",
	})

	dialector := New(Config{DSN: url})

	var err error

	logger.Default.LogMode(logger.Info)

	s.db, err = gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            false,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	s.Require().Nil(err)

	err = s.db.Exec("CREATE TABLE TEST_TABLE (ID INT PRIMARY KEY, NAME VARCHAR2(64 CHAR) NOT NULL)").Error
	s.Require().Nil(err)
}

func TestOracle19c(t *testing.T) {
	suite.Run(t, new(Oracle19Suite))
}

func TestOracle23c(t *testing.T) {
	suite.Run(t, new(Oracle23Suite))
}

func (s *OracleSuite) TestEntityFound() {
	var err error

	err = s.db.Create(&Model{
		Id:   1,
		Name: "test",
	}).Error
	s.Require().Nil(err)

	var m Model

	err = s.db.
		Where("id = ?", 1).
		First(&m).
		Error

	s.Require().Nil(err)
}

func (s *OracleSuite) TestEntityNotFound() {
	var err error

	var m Model

	err = s.db.
		Where("id = ?", 2).
		First(&m).
		Error

	s.Require().ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *OracleSuite) TestEntityFoundForUpdate() {
	var err error

	err = s.db.Create(&Model{
		Id:   1,
		Name: "test",
	}).Error
	s.Require().Nil(err)

	var m Model

	err = s.db.Transaction(func(tx *gorm.DB) error {
		return tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", 1).
			First(&m).
			Error
	})
	s.Require().Nil(err)
}

func (s *OracleSuite) TestEntityNotFoundForUpdate() {
	var err error

	var m Model

	err = s.db.Transaction(func(tx *gorm.DB) error {
		return tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", 2).
			First(&m).
			Error
	})
	s.Require().ErrorIs(err, gorm.ErrRecordNotFound)
}
