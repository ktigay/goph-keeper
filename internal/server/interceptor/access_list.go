package interceptor

import (
	"github.com/ktigay/goph-keeper/internal/contracts/v1/auth"
	"github.com/ktigay/goph-keeper/internal/contracts/v1/data"
)

// AccessList список для проверки доступа методов.
func AccessList() map[string]bool {
	return map[string]bool{
		auth.AuthService_Register_FullMethodName:                false,
		auth.AuthService_Login_FullMethodName:                   false,
		data.UserDataService_CreateUserDataItem_FullMethodName:  true,
		data.UserDataService_UpdateUserDataItem_FullMethodName:  true,
		data.UserDataService_GetUserDataItem_FullMethodName:     true,
		data.UserDataService_GetUserDataItems_FullMethodName:    true,
		data.UserDataService_DeleteUserDataItems_FullMethodName: true,
	}
}
