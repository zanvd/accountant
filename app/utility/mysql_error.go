package utility

import (
	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"net/http"
)

func GetMySQLErrorMessage(err error) string {
	message := "An unknown error occurred."
	if mySQLError, ok := err.(*mysql.MySQLError); ok {
		switch mySQLError.Number {
		case mysqlerr.ER_NO_REFERENCED_ROW, mysqlerr.ER_NO_REFERENCED_ROW_2:
			message = "Selected referenced entity does not exist."
		case mysqlerr.ER_ROW_IS_REFERENCED, mysqlerr.ER_ROW_IS_REFERENCED_2:
			message = "Cannot delete entries still in use."
		}
	}
	return message
}

func MapMySQLErrorToHttpCode(err error) int {
	code := http.StatusInternalServerError
	if mySQLError, ok := err.(*mysql.MySQLError); ok {
		switch mySQLError.Number {
		case mysqlerr.ER_NO_REFERENCED_ROW,
			mysqlerr.ER_NO_REFERENCED_ROW_2,
			mysqlerr.ER_ROW_IS_REFERENCED,
			mysqlerr.ER_ROW_IS_REFERENCED_2:
			code = http.StatusBadRequest
		}
	}
	return code
}
