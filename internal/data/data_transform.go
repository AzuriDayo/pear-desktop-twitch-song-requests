package data

import (
	"database/sql"

	. "github.com/azuridayo/pear-desktop-twitch-song-requests/gen/table"
	. "github.com/go-jet/jet/v2/sqlite"
)

type dataTransformType interface {
	GetKey() string
	Transform(*sql.DB) error
	IsNecessary() bool
}

func GetDataTransformTypes() []dataTransformType {
	return []dataTransformType{
		dataTransformBackportRequestersUserID,
	}
}

func IsDataTransformed(db *sql.DB, k dataTransformType) (bool, error) {
	v := struct {
		Value bool
	}{}
	stmt := SELECT(DataTransforms.Value.AS("value")).FROM(DataTransforms).WHERE(DataTransforms.Key.EQ(String(k.GetKey()))).LIMIT(1)

	err := stmt.Query(db, &v)
	if err != nil {
		return false, err
	}
	return v.Value, nil
}

func SetDataTransformed(db *sql.DB, k dataTransformType) error {
	stmt := DataTransforms.UPDATE(DataTransforms.Value).SET(Bool(true)).WHERE(DataTransforms.Key.EQ(String(k.GetKey())))

	_, err := stmt.Exec(db)
	if err != nil {
		return err
	}
	return nil
}
