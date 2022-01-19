package adapters

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rafaelcalleja/go-kit/internal/store/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductRepository_Save_Err(t *testing.T) {
	productId := "37a0f027-15e6-47cc-a5d2-64183281087e"

	product, err := domain.NewProduct(productId)
	require.NoError(t, err)

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlMock.ExpectExec(
		"INSERT INTO products (id) VALUES (?)").
		WithArgs(productId).
		WillReturnError(errors.New("something-failed"))

	repo := NewMysqlProductRepository(db, 1*time.Millisecond)

	err = repo.Save(context.Background(), product)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Error(t, err)
}

func TestProductRepository_Save_Success(t *testing.T) {
	productId := "37a0f027-15e6-47cc-a5d2-64183281087e"

	product, err := domain.NewProduct(productId)
	require.NoError(t, err)

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlMock.ExpectExec(
		"INSERT INTO products (id) VALUES (?)").
		WithArgs(productId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := NewMysqlProductRepository(db, 1*time.Millisecond)

	err = repo.Save(context.Background(), product)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
}

func TestProductRepository_Of_Success(t *testing.T) {
	productId := "37a0f027-15e6-47cc-a5d2-64183281087e"

	productIdVO, err := domain.NewProductId(productId)
	require.NoError(t, err)

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow("37a0f027-15e6-47cc-a5d2-64183281087e")

	sqlMock.ExpectQuery(
		"SELECT products.id FROM products WHERE id = ? LIMIT 1").
		WithArgs(productId).
		WillReturnRows(rows)

	repo := NewMysqlProductRepository(db, 60*time.Second)

	_, err = repo.Of(context.Background(), productIdVO)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
}

func TestProductRepository_Of_Empty(t *testing.T) {
	productId := "37a0f027-15e6-47cc-a5d2-64183281087e"

	productIdVO, err := domain.NewProductId(productId)
	require.NoError(t, err)

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlMock.ExpectQuery(
		"SELECT products.id FROM products WHERE id = ? LIMIT 1").
		WithArgs(productId).
		WillReturnRows(sqlmock.NewRows([]string{}))

	repo := NewMysqlProductRepository(db, 60*time.Second)

	_, err = repo.Of(context.Background(), productIdVO)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Error(t, err)
}
