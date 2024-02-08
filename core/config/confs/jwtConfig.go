package confs

// JwtConfig
type (
	JwtConfig struct {
		Jwt jwt `yaml:"auth"`
	}
	jwt struct {
		JwtSecret string `yaml:"jwtSecret"` //JWT密钥
		Expire    int    `yaml:"expire"`    //JWT过期时间，单位为秒
	}
)

func (c *JwtConfig) IsLoad() bool {
	return true
}
