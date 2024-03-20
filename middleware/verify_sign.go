package middleware

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/binbinly/gin-pkg/app"
	"github.com/binbinly/pkg/errno"
	"github.com/binbinly/pkg/signature"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const (
	// HeaderSignToken 签名验证 Authorization，Header 中传递的参数
	HeaderSignToken = "Auth"

	// HeaderSignTokenDate 签名验证 Date，Header 中传递的参数
	HeaderSignTokenDate = "Auth-Date"

	// HeaderSignTokenTimeout 签名有效期为 1 分钟
	HeaderSignTokenTimeout = time.Minute
)

type VerifyOpts struct {
	timeout time.Duration
	except  []string
}

type VerifyOpt func(opts *VerifyOpts)

func WithTimeout(timeout time.Duration) VerifyOpt {
	return func(opts *VerifyOpts) {
		opts.timeout = timeout
	}
}

func WithExcept(uris []string) VerifyOpt {
	return func(opts *VerifyOpts) {
		opts.except = uris
	}
}

// VerifySign 验证签名
func VerifySign(secretKey string, opts ...VerifyOpt) gin.HandlerFunc {
	os := &VerifyOpts{
		timeout: HeaderSignTokenTimeout,
		except:  make([]string, 0),
	}
	for _, opt := range opts {
		opt(os)
	}

	return func(c *gin.Context) {
		if len(os.except) > 0 {
			for _, uri := range os.except {
				if c.Request.URL.Path == uri {
					c.Next()
					return
				}
			}
		}

		// 签名信息
		authorization := c.GetHeader(HeaderSignToken)
		if authorization == "" {
			app.Error(c, errno.ErrSignParam)
			return
		}

		// 时间信息
		timestamp, _ := strconv.ParseInt(c.GetHeader(HeaderSignTokenDate), 10, 64)
		if timestamp == 0 {
			app.Error(c, errno.ErrSignParam)
			return
		}

		// 通过签名信息获取 key
		authorizationSplit := strings.Split(authorization, " ")
		if len(authorizationSplit) < 2 {
			app.Error(c, errno.ErrSignParam)
			return
		}
		key := authorizationSplit[0]

		data, _ := c.GetRawData()
		// 这里防止body只能读一次
		c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
		params := getParams(c)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(data))

		st := signature.New(key, secretKey, HeaderSignTokenTimeout)
		ok, err := st.Verify(authorization, timestamp, params)
		if err != nil || !ok {
			app.Error(c, errno.ErrSignParam)
			return
		}

		c.Next()
	}
}

// getParams 获取请求参数
func getParams(c *gin.Context) map[string]any {
	params := map[string]any{}
	if c.Request.Method == "POST" {
		contextType := c.Request.Header.Get("Content-Type")
		if strings.Index(contextType, "json") >= 0 {
			if err := c.ShouldBindBodyWith(&params, binding.JSON); err != nil {
				return nil
			}
		} else {
			_ = c.Request.ParseMultipartForm(32 << 20)
			if len(c.Request.PostForm) > 0 {
				for k, v := range c.Request.PostForm {
					params[k] = v[0]
				}
			}
		}
	} else {
		var tmpParams = make(map[string]string)
		if err := c.ShouldBind(&tmpParams); err != nil {
			return nil
		}
		for k, v := range tmpParams {
			params[k] = v
		}
	}
	params["method"] = c.Request.Method
	params["path"] = c.Request.URL.Path
	return params
}
