package query

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Query interface {
	Validate() error
}

type SongQuery struct {
	Name        string
	Group       string	`sql_related:"name"`
	ReleaseDate string 	`valid:"date"`
	Text        string	`sql:"substring"`
	Link        string 	`valid:"link"`
	Page        int		`sql:"-"`
	Limit		int		`sql:"-"`
}

type GroupQuery struct {
	Name string
}

func (q *SongQuery) Validate() error {
	_, err := govalidator.ValidateStruct(*q)
	return err
}

func (q *SongQuery) GenerateSQL() string {
	buf := new(bytes.Buffer)
	buf.WriteString(`SELECT s."id", s."name", s."releaseDate", s."text", s."link", g."name" from songs s 
		JOIN "groups" g ON s."groupId" = g."id"`)
	v := reflect.ValueOf(*q)

	limit := q.Limit
	if q.Limit == 0 {
		limit = 10
	} 
		
	offset := limit * (q.Page - 1)
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		if field.IsZero() || fieldType.Tag.Get("sql") == "-" {
			continue
		}

		var (
			tableShort rune = 's'
			searchPattern string = fmt.Sprintf(`= '%v'`, field.Interface())
			fieldName string = fieldType.Name
		)

		if i == 0 {
			buf.WriteString(` WHERE `)
		} else {
			buf.WriteString(` AND `)
		}

		if fieldType.Tag.Get("sql") == "substring" {
			searchPattern = fmt.Sprintf(`LIKE '%%%v%%'`, field.Interface())
		}

		if related := fieldType.Tag.Get("sql_related"); related != "" {
			fieldName = related
			tableShort = 'g'
		}

		fmt.Fprintf(buf, `%c."%s%s" %s`, tableShort, strings.ToLower(fieldName[:1]), fieldName[1:], searchPattern)
	}

	if q.Page != 0 {
		fmt.Fprintf(buf, " LIMIT %d OFFSET %d", limit, offset)
	}

	return buf.String()
}

func (q *GroupQuery) Validate() error {
	return nil
}
