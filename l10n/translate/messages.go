package translate

type Key string

const (
	KeyEmailUniqueness         Key = "EmailUniqueness"
	KeyUsernameUniqueness      Key = "UsernameUniqueness"
	KeyFarmerAlreadyFavorite   Key = "FarmerAlreadyFavorite"
	KeyFarmerAlreadyMuted      Key = "FarmerAlreadyMuted"
	KeyInvalidToken            Key = "InvalidToken"
	KeyUserBanned              Key = "UserBanned"
	KeyUserNotFound            Key = "UserNotFound"
	KeyUserSuspended           Key = "UserSuspended"
	KeyPasswordNotMatched      Key = "PasswordNotMatched"
	KeyInternalServerError     Key = "InternalServerError"
	KeyNearestPotagersNotFound Key = "NearestPotagersNotFound"
	KeyNearestAlimentsNotFound Key = "NearestAlimentsNotFound"
	KeyPotagerNotFound         Key = "PotagerNotFound"
)

var AllKey = []Key{
	KeyEmailUniqueness,
	KeyUsernameUniqueness,
	KeyFarmerAlreadyFavorite,
	KeyFarmerAlreadyMuted,
	KeyInvalidToken,
	KeyUserBanned,
	KeyUserNotFound,
	KeyUserSuspended,
	KeyPasswordNotMatched,
	KeyInternalServerError,
	KeyNearestPotagersNotFound,
	KeyNearestAlimentsNotFound,
	KeyPotagerNotFound,
}

func (m Key) String() string {
	return string(m)
}
