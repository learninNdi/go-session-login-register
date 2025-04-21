package libraries

import (
	"database/sql"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/learninNdi/go-session-login-register/config"

	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validation struct {
	conn *sql.DB
}

func NewValidation() *Validation {
	conn, err := config.DBConn()

	if err != nil {
		panic(err)
	}

	return &Validation{
		conn: conn,
	}
}

func (v *Validation) Init() (*validator.Validate, ut.Translator) {
	// VALIDASI USING GO PLAYGROUND
	// memanggil translator package
	translator := en.New()
	uni := ut.New(translator, translator)
	trans, _ := uni.GetTranslator("en")

	validate := validator.New()

	// register default translation (en)
	en_translations.RegisterDefaultTranslations(validate, trans)

	// change default label
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		labelName := field.Tag.Get("label")

		return labelName
	})

	// translate error to Bahasa Indonesia
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} tidak boleh kosong", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())

		return t
	})

	// custom validation
	validate.RegisterValidation("isUnique", func(fl validator.FieldLevel) bool {
		params := fl.Param()
		splittedParams := strings.Split(params, "-")

		tableName := splittedParams[0]
		columnName := splittedParams[1]
		fieldValue := fl.Field().String()

		return v.checkIsUnique(tableName, columnName, fieldValue)
	})

	validate.RegisterTranslation("isUnique", trans, func(ut ut.Translator) error {
		return ut.Add("isUnique", "{0} sudah digunakan", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("isUnique", fe.Field())

		return t
	})

	return validate, trans
}

func (v *Validation) Struct(s interface{}) interface{} {
	validate, trans := v.Init()

	vErrors := make(map[string]interface{})

	err := validate.Struct(s)

	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			vErrors[e.StructField()] = e.Translate(trans)
		}
	}

	if len(vErrors) > 0 {
		return vErrors
	}

	return nil
}

func (v *Validation) checkIsUnique(tableName, columnName, fieldValue string) bool {

	row, _ := v.conn.Query("select "+columnName+" from "+tableName+" where "+columnName+" = ?", fieldValue)

	defer row.Close()

	var result string

	for row.Next() {
		row.Scan(&result)
	}

	return result != fieldValue
}
