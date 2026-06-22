package api

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/req"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/idcodec"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/pkg/validate"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/common/strs"

	"bbs-go/internal/cache"
	"bbs-go/internal/handlers/render"
	"bbs-go/internal/services"
)

func UserCurrent(ctx *gin.Context) {
	if !config.Instance.Installed {
		ginx.WriteJSON(ctx, nil)
		return
	}
	user := common.GetCurrentUser(ctx)
	if user != nil {
		ginx.WriteJSON(ctx, render.BuildUserProfile(user))
		return
	}
	ginx.WriteJSON(ctx, nil)
}

func UserDetail(ctx *gin.Context) {
	userIdStr := ctx.Param("id")

	userId := idcodec.Decode(userIdStr)
	user := cache.UserCache.Get(userId)
	if user != nil && user.Status != constants.StatusDeleted {
		ginx.WriteJSON(ctx, render.BuildUserDetail(user))
		return
	}
	ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("user.not_found")))

}

func UserUpdate(ctx *gin.Context) {
	userIdStr := ctx.Param("id")

	userId := idcodec.Decode(userIdStr)
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}
	if user.Id != userId {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("user.no_permission")))
		return
	}
	var req req.UserUpdateReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	nickname := strings.TrimSpace(req.Nickname)
	homePage := req.HomePage
	description := req.Description
	gender := strings.TrimSpace(req.Gender)

	var (
		minLength = constants.NicknameMinLengthEnUS
		maxLength = constants.NicknameMaxLengthEnUS
	)
	if strings.EqualFold(string(config.Instance.Language), string(config.LanguageZhCN)) {
		minLength = constants.NicknameMinLengthZhCN
		maxLength = constants.NicknameMaxLengthZhCN
	}
	if nicknameLength := utf8.RuneCountInString(nickname); nicknameLength < minLength || nicknameLength > maxLength {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Getf("user.nickname_length_invalid", minLength, maxLength)))
		return
	}

	if strs.IsNotBlank(gender) {
		if gender != string(constants.GenderMale) && gender != string(constants.GenderFemale) {
			ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("user.gender_error")))
			return
		}
	}

	if len(homePage) > 0 && validate.IsURL(homePage) != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("user.homepage_error")))
		return
	}

	err := services.UserService.Updates(user.Id, map[string]any{
		"nickname":    nickname,
		"home_page":   homePage,
		"description": description,
		"gender":      gender,
	})
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func UserUpdateAvatar(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}
	avatar := strings.TrimSpace(params.FormValue(ctx, "avatar"))
	if len(avatar) == 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("user.avatar_empty")))
		return
	}
	err := services.UserService.UpdateAvatar(user.Id, avatar)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func UserSetUsername(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}
	username := strings.TrimSpace(params.FormValue(ctx, "username"))
	err := services.UserService.SetUsername(user.Id, username)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)
}

func UserSetEmail(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}
	email := strings.TrimSpace(params.FormValue(ctx, "email"))
	err := services.UserService.SetEmail(user.Id, email)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func UserSetPassword(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}
	password := params.FormValue(ctx, "password")
	rePassword := params.FormValue(ctx, "rePassword")
	err := services.UserService.SetPassword(user.Id, password, rePassword)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func UserUpdatePassword(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}
	var req req.PasswordUpdateReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	if err := services.UserService.UpdatePassword(user.Id, req.OldPassword, req.Password, req.RePassword); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func UserSetBackgroundImage(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}
	backgroundImage := params.FormValue(ctx, "backgroundImage")
	if strs.IsBlank(backgroundImage) {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("user.upload_image_required")))
		return
	}
	if err := services.UserService.UpdateBackgroundImage(user.Id, backgroundImage); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func UserForbidden(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}
	var req req.UserForbiddenReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	userId := idcodec.Decode(req.UserId)
	if !services.PermissionService.CanForbiddenUser(user, req.Days) {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("user.no_permission")))
		return
	}
	if userId < 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("param: userId required"))
		return
	}
	if req.Days == 0 {
		services.UserService.RemoveForbidden(user.Id, userId, ctx.Request)
	} else {
		if err := services.UserService.Forbidden(user.Id, userId, req.Days, req.Reason, ctx.Request); err != nil {
			ginx.WriteJSON(ctx, err)
			return
		}
	}
	ginx.WriteJSON(ctx, nil)

}

func UserSendVerifyEmail(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}
	if err := services.UserService.SendEmailVerifyEmail(user.Id); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func UserVerifyEmail(ctx *gin.Context) {
	token := params.FormValue(ctx, "token")
	if strs.IsBlank(token) {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("Illegal request"))
		return
	}
	var (
		email string
		err   error
	)
	if email, err = services.UserService.VerifyEmail(token); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, map[string]any{"email": email})

}

func UserWxBindInfo(ctx *gin.Context) {
	user, err := common.CheckLogin(ctx)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	thirdUser := services.ThirdUserService.GetByUserId(user.Id, constants.ThirdTypeWeixin)
	if thirdUser != nil {
		ginx.WriteJSON(ctx, map[string]any{
			"bind":     true,
			"nickname": thirdUser.Nickname,
			"avatar":   thirdUser.Avatar,
		})
		return
	}
	ginx.WriteJSON(ctx, map[string]any{
		"bind": false,
	})

}

func UserGoogleBindInfo(ctx *gin.Context) {
	user, err := common.CheckLogin(ctx)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	thirdUser := services.ThirdUserService.GetByUserId(user.Id, constants.ThirdTypeGoogle)
	if thirdUser != nil {
		ginx.WriteJSON(ctx, map[string]any{
			"bind":     true,
			"nickname": thirdUser.Nickname,
			"avatar":   thirdUser.Avatar,
		})
		return
	}
	ginx.WriteJSON(ctx, map[string]any{
		"bind": false,
	})

}

func UserGithubBindInfo(ctx *gin.Context) {
	user, err := common.CheckLogin(ctx)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	thirdUser := services.ThirdUserService.GetByUserId(user.Id, constants.ThirdTypeGithub)
	if thirdUser != nil {
		ginx.WriteJSON(ctx, map[string]any{
			"bind":     true,
			"nickname": thirdUser.Nickname,
			"avatar":   thirdUser.Avatar,
		})
		return
	}
	ginx.WriteJSON(ctx, map[string]any{
		"bind": false,
	})

}
