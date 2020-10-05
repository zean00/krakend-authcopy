package authcopy

import (
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/gin-gonic/gin"
)

const authHeader = "Authorization"
const Namespace = "github_com/zean00/authcopy"

//CopyConfig authcopy config
type CopyConfig struct {
	CookieKey string
	QueryKey  string
	Overwrite bool
}

//New create middleware
func New(logger logging.Logger, config config.ExtraConfig) gin.HandlerFunc {
	cfg := configGetter(logger, config)
	if cfg == nil {
		logger.Info("[authcopy] Empty config")
		return func(c *gin.Context) {
			c.Next()
		}
	}
	return func(c *gin.Context) {

		if token := c.Request.Header.Get(authHeader); token != "" && !cfg.Overwrite {
			c.Next()
			return
		}

		cookie, err := c.Request.Cookie(cfg.CookieKey)
		if err == nil {
			logger.Debug("[authcopy] Copying from cookie")
			c.Request.Header.Set(authHeader, "Bearer "+cookie.Value)
		}

		query := c.Request.URL.Query().Get(cfg.QueryKey)
		if query != "" {
			logger.Debug("[authcopy] Copying from query")
			c.Request.Header.Set(authHeader, "Bearer "+query)
			val := c.Request.URL.Query()
			val.Del(cfg.QueryKey)
			c.Request.URL.RawQuery = val.Encode()
		}

		c.Next()
	}
}

func configGetter(logger logging.Logger, config config.ExtraConfig) *CopyConfig {
	v, ok := config[Namespace]
	if !ok {
		return nil
	}
	tmp, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	cfg := new(CopyConfig)
	cfg.Overwrite = false

	ck, ok := tmp["cookie_key"].(string)
	if ok {
		cfg.CookieKey = ck
	}

	qk, ok := tmp["query_key"].(string)
	if ok {
		cfg.QueryKey = qk
	}

	ow, ok := tmp["overwrite"].(bool)
	if ok {
		cfg.Overwrite = ow
	}

	if cfg.CookieKey == "" && cfg.QueryKey == "" {
		return nil
	}

	return cfg
}
