package gorm_utils

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/response_builder"
	"math"
)

func GetGormInternalOrNoRecordsFoundError(
	err error,
	noRecordsErrMsg string,
	internalErrMsg string,
) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(noRecordsErrMsg)
	} else {
		return errors.New(internalErrMsg + ": " + err.Error())
	}
}

// ToNullInt64 validates a sql.NullInt64 if incoming string evaluates to an integer, invalidates if it does not
func ToNullInt64(i int) sql.NullInt64 {
	return sql.NullInt64{Int64: int64(i), Valid: i != 0}
}

// GetGormPaginatedQuery executes the paginated query, and returns a paginated object
func GetGormPaginatedQuery(
	tx *gorm.DB,
	outStruct interface{},
	page int,
	rows int,
	rowLimit int,
) (response_builder.Pagination, error) {
	var pages []int
	var countContainer int64
	var count int
	var pagination response_builder.Pagination

	// ====================================================
	// Determine count by executing the query
	err := tx.Count(&countContainer).Error
	if err != nil {
		return pagination, err
	}
	count = int(countContainer)

	// ====================================================
	// Pagination logic
	pageStart := 0
	pageEnd := rowLimit
	totalPages := int(math.Ceil(float64(count) / float64(rows)))

	if page > totalPages || totalPages == 1 {
		page = 0
	}

	var rowsPerPage int

	if rows > count {
		rowsPerPage = count
	} else {
		rowsPerPage = rows
	}

	var offset int

	if count != 0 {
		if page == 0 {
			offset = 0
		} else if totalPages == 0 {
			offset = 0
		} else {
			if page >= totalPages {
				offset = 0
			} else {
				offset = page * rows
			}
		}
	} else {
		offset = 0
	}

	for !(page >= 0 && page <= pageEnd) {
		pageStart = pageStart + rowLimit
		pageEnd = pageEnd + rowLimit
	}

	for i := pageStart; i <= pageEnd; i++ {
		if i <= totalPages - 1 {
			pages = append(pages, i)
		} else {
			if len(pages) > 0 {
				previousPage := pages[0] - 1

				if previousPage > 1 {
					pages = append([]int{previousPage}, pages...)
				}
			}
		}
	}

	if len(pages) == 0 {
		pages = append(pages, 0)
	}

	resultCount := 0
	hasNextPage := false

	for _, ele := range pages {
		if ele > page {
			hasNextPage = true
			break
		}
	}

	if hasNextPage == true {
		resultCount = rows
	} else {
		resultCount = count - offset
	}

	// ====================================================
	// Populate pagination
	// Old code (restore this when corrected below)
	pagination.Rows = rows
	pagination.Page = page
	pagination.Pages = pages
	pagination.RowsPerPage = rowsPerPage
	pagination.Offset = offset
	pagination.Page = page
	pagination.ResultCount = resultCount
	pagination.TotalCount = count

	// ====================================================
	// Execute query from calculated pagination variables
	return pagination, tx.Offset(offset).Limit(rowLimit).Find(outStruct).Error
}