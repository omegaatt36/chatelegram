package chatelegram

type config interface {
	use(*Service)
}

// UseAllowedUsers implements config. telebot.Bot will be add whitelist.
type UseAllowedUsers struct {
	AllowedUsers []int64
}

var _ config = &UseAllowedUsers{}

func (o UseAllowedUsers) use(service *Service) {
	service.allowedUsers = o.AllowedUsers
}
