package user

type User struct{
	Auth
	ID string
	Features *[]string
}

type Auth struct{
	SessionToken string
	Username string
	Password string
}

func (u *Auth) IsLoggedIn() bool{
	return u.SessionToken != ""
}

func (u *Auth) LoginUser(userId string, hashedPassword string){
	u.SessionToken = userId
}


func (u *User) HasFeature(feature string) bool{
	if u.Features == nil {
		return false
	}
	for _, f := range *u.Features {
		if f == feature {
			return true
		}
	}
	return false
}

type Features interface{
	HasFeature(feature string) bool
}