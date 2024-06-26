package common

var (
	modLevelStrings = [...]string{"", "janitors", "moderators", "owners",
		"admin"}
)

// ModerationAction is an action performable by moderation staff
type ModerationAction uint8

// All supported moderation actions
const (
	BanPost ModerationAction = iota
	UnbanPost
	DeletePost
	DeleteImage
	SpoilerImage
	LockThread
	DeleteBoard
	MeidoVision
	PurgePost
	ShadowBinPost
	RedirectIP
	RedirectThread
	ToggleSticky
	ConfigureBoard
	AssignStaff
	BoardAssets
	ConfigureServer
	AdminNotification
	PlaylistLock
)

// Contains fields of a post moderation log entry
type ModerationEntry struct {
	Type   ModerationAction `json:"type"`
	Length uint64           `json:"length"`
	By     string           `json:"by"`
	Data   string           `json:"data"`
}

// ModerationLevel defines the level required to perform an action or the
// permission level held by a user
type ModerationLevel int8

// Returns string representation of moderation level
func (l ModerationLevel) String() string {
	if l < Janitor {
		return ""
	}
	return modLevelStrings[int(l)]
}

// All available moderation levels
const (
	NotLoggedIn ModerationLevel = iota - 1
	NotStaff
	Janitor
	Moderator
	BoardOwner
	Admin
)

// Pseudo-dictionary of moderation privilege level required to perform action
var ActionPrivilege = [...]ModerationLevel{
	BanPost:           Moderator,
	UnbanPost:         Moderator,
	DeletePost:        Janitor,
	DeleteImage:       Janitor,
	SpoilerImage:      Janitor,
	LockThread:        Moderator,
	DeleteBoard:       BoardOwner,
	MeidoVision:       Janitor,
	PurgePost:         Admin,
	ShadowBinPost:     Moderator,
	RedirectIP:        Admin,
	RedirectThread:    Admin,
	ToggleSticky:      Moderator,
	ConfigureBoard:    BoardOwner,
	AssignStaff:       BoardOwner,
	ConfigureServer:   Admin,
	AdminNotification: Admin,
	PlaylistLock:      Moderator,
}
