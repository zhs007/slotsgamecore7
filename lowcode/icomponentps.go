package lowcode

type IComponentPS interface {
	// SetPublicJson
	SetPublicJson(str string) error
	// SetPrivateJson
	SetPrivateJson(str string) error
	// GetPublicJson
	GetPublicJson() string
	// GetPrivateJson
	GetPrivateJson() string
	// Clone
	Clone() IComponentPS
}
