package convert

import "time"

const (
	DbDateFormat      = "2006-01-02"
	DisplayDateFormat = "02. 01. 2006"
)

func ChangeDateFormat(destFormat string, sourceFormat string, sourceDate string) (destDate string) {
	t, err := time.Parse(sourceFormat, sourceDate)
	if err != nil {
		return
	}
	destDate = t.Format(destFormat)
	return
}

func CurrentDateInDbFormat() string {
	return time.Now().Format(DbDateFormat)
}

func CurrentDateInDisplayFormat() string {
	return time.Now().Format(DisplayDateFormat)
}

func DbToDisplayDate(dbDate string) string {
	return ChangeDateFormat(DisplayDateFormat, DbDateFormat, dbDate)
}

func DisplayTimeToDb(displayDate string) string {
	return ChangeDateFormat(DbDateFormat, DisplayDateFormat, displayDate)
}
