package validation

import (
	"myapp/models"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

/*
 ログインパラメータのバリデーション
*/
func SignInValidate(signInRequest models.SingInRequest) error {
	return validation.ValidateStruct(&signInRequest,
		validation.Field(
			&signInRequest.Email,
			validation.Required.Error("メールアドレスは必須入力です。"),
			validation.RuneLength(5, 40).Error("メールアドレスは 5～40 文字です"),
			is.Email.Error("メールアドレスの形式が間違っています。"),
		),
		validation.Field(
			&signInRequest.Password,
			validation.Required.Error("パスワードは必須入力です。"),
			validation.Length(6, 20).Error("パスワードは6文字以上、20字以内で入力してください。"),
			is.Alphanumeric.Error("パスワードは英数字で入力してください。"),
		),
	)
}


/*
 会員登録パラメータのバリデーション
*/
func SignUpValidate(signUpRequest models.SignUpRequest) error {
	return validation.ValidateStruct(&signUpRequest,
		validation.Field(
			&signUpRequest.Name,
			validation.Required.Error("お名前は必須入力です。"),
			validation.RuneLength(5, 10).Error("お名前は 5～10 文字です"),
		),
		validation.Field(
			&signUpRequest.Email,
			validation.Required.Error("メールアドレスは必須入力です。"),
			validation.RuneLength(5, 40).Error("メールアドレスは 5～40 文字です"),
			is.Email.Error("メールアドレスの形式が間違っています。"),
		),
		validation.Field(
			&signUpRequest.Password,
			validation.Required.Error("パスワードは必須入力です。"),
			validation.RuneLength(6, 20).Error("パスワードは 6~20 文字です。"),
			is.Alphanumeric.Error("パスワードは英数字で入力してください。"),
		),
	)
}