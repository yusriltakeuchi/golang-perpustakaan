package util

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func Validate[T any](data T) map[string]string {
	err := validator.New().Struct(data)
	res := map[string]string{}
	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			res[v.StructField()] = TranslateTag(v)
		}
	}
	return res
}

// TranslateTag translates validation tags to Indonesian messages.
func TranslateTag(fd validator.FieldError) string {
	switch fd.ActualTag() {
	case "required":
		return fmt.Sprintf("field %s wajib diisi", fd.StructField())
	case "min":
		return fmt.Sprintf("field %s size minimal %s", fd.StructField(), fd.Param())
	case "unique":
		return fmt.Sprintf("field %s tidak boleh ada yang sama", fd.StructField())
	}
	return "validasi gagal"
}
