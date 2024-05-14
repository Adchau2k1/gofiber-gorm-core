package utils

import (
	"github.com/dlclark/regexp2"
)

type checkFunc func(value string) string

type validator struct {
	Username checkFunc
	Password checkFunc
	Fullname checkFunc
}

func checkUsername(value string) string {
	re := regexp2.MustCompile(`^[a-zA-Z][a-zA-Z0-9]{5,}$`, 0)
	match, _ := re.MatchString(value)

	if match {
		return ""
	}

	return "Tài khoản: tối thiểu 6 kí tự, bắt đầu bằng chữ, không chứa kí tự đặc biệt!"
}

func checkPassword(value string) string {
	re := regexp2.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[_@$!%*?&])[A-Za-z\d_@$!%*?&]{6,}$`, 0)
	match, _ := re.MatchString(value)

	if match {
		return ""
	}

	return "Mật khẩu: tối thiểu 6 kí tự, phải có ít nhất 1 kí tự đặc biệt, in hoa, in thường và số!"
}

func checkEmail(value string) string {
	re := regexp2.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, 0)
	match, _ := re.MatchString(value)

	if match {
		return ""
	}

	return "Email không hợp lệ!"
}

func checkFullname(value string) string {
	pattern := `^[a-zA-Z\sÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚĂĐĨŨƠàáâãèéêìíòóôõùúăđĩũơƯĂẠẢẤẦẨẪẬẮẰẲẴẶẸẺẼỀỀỂưăạảấầẩẫậắằẳẵặẹẻẽềềểỄỆỈỊỌỎỐỒỔỖỘỚỜỞỠỢỤỦỨỪễệỉịọỏốồổỗộớờởỡợụủứừỬỮỰỲỴÝỶỸửữựỳỵỷỹ]{6,}$`
	re := regexp2.MustCompile(pattern, 0)
	match, _ := re.MatchString(value)

	if match {
		return ""
	}

	return "Tên hiển thị: tối thiểu 6 kí tự, chỉ bao gồm chữ và khoảng trắng!"
}

var Validator = validator{
	Username: checkUsername,
	Password: checkPassword,
	Fullname: checkFullname,
}
