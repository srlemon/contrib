package session

import (
	"crypto/ecdsa"
	error2 "github.com/olongfen/contrib/utils"
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"

	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
)

// Key 密钥实例
type Key struct {
	encryptMethod string // 默认加密方法
	// hook
	hookSessionCheck func(session *Session) error  // 二次检测session合法性
	hookTokenCheck   func(token interface{}) error // 二次检测token合法性
	// encrypt rsa
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey
	// encrypt hmac
	hmacSecret []byte // HMAC密钥
	// encrypt ecdsa
	ecdsaPrivateKey *ecdsa.PrivateKey
	ecdsaPublicKey  *ecdsa.PublicKey
}

// SetHookSessionCheck 二次检测session合法性
func (k *Key) SetHookSessionCheck(f func(session *Session) error) {
	k.hookSessionCheck = f
}

// SetHookTokenCheck 二次检测token合法性
func (k *Key) SetHookTokenCheck(f func(token interface{}) error) {
	k.hookTokenCheck = f
}

// SessionDecode 解析出session
func (k *Key) SessionDecode(inf interface{}) (ret *Session, err error) {
	if k == nil {
		err = error2.ErrSessionKeyUndefined
		return
	}
	var (
		val map[string]interface{}
	)
	//if val, err = SessionDecode(inf); err != nil {
	//	return
	//}
	if val, err = k.TokenDecode(inf); err != nil {
		return
	}
	// 手动解析
	ret = new(Session)
	// 逻辑属性
	if v, ok := val[TokenTagExp]; ok {
		switch v_ := v.(type) {
		case int32:
			ret.ExpireTime = int64(v_)
			break
		case int64:
			ret.ExpireTime = v_
			break
		case float64:
			ret.ExpireTime = int64(v_)
			break
		default:
			err = fmt.Errorf("unknown %s value: %v", TokenTagExp, v)
			break
		}
	}
	if v, ok := val[TokenTagUid]; ok {
		if s, ok := v.(string); ok {
			ret.UID = s
		}
	}
	if v, ok := val[TokenTagContent]; ok {
		if s, ok := v.(map[string]interface{}); ok {
			ret.Content = s
		}
	}

	// 合法性
	if err = ret.Valid(); err != nil {
		return
	}
	if k.hookSessionCheck != nil {
		if err = k.hookSessionCheck(ret); err != nil {
			return
		}
	}
	return
}

// SessionEncode 将session编码为token
func (k *Key) SessionEncode(s *Session) (token string, err error) {

	if s == nil {
		err = error2.ErrSessionUndefined
		return
	} else if err = s.Valid(); err != nil {
		return
	}

	// 手动填充map
	var m = make(map[string]interface{})

	// 必填:逻辑属性
	m[TokenTagExp] = s.ExpireTime // ExpireTime

	if len(s.UID) > 0 {
		// uid
		m[TokenTagUid] = s.UID
	}
	if s.Content != nil {
		m[TokenTagContent] = s.Content
	}
	token, err = k.TokenEncode(m) // 默认加密
	return
}

// TokenDecode 解析出需要的值
func (k *Key) TokenDecode(inf interface{}) (ret map[string]interface{}, err error) {
	if k == nil {
		err = error2.ErrSessionKeyUndefined
		return
	}
	switch v := inf.(type) {
	case string:
		ret, err = k.tokenParse(v)
		break
	case []byte:
		ret, err = k.tokenParse(string(v))
		break
	case json.RawMessage:
		ret, err = k.tokenParse(string(v))
		break
	case *http.Request:
		var token string
		if token, err = PubGetTokenFromReq(v); err != nil {
			return
		}
		ret, err = k.tokenParse(token)
		break
	default:
		err = error2.ErrTokenParseTypeNotSupport
		break
	}

	// 执行hook函数
	if err == nil && k.hookTokenCheck != nil {
		if err = k.hookTokenCheck(inf); err != nil {
			return
		}
	}
	return
}

// SetRSA 设置RSA密钥对
func (k *Key) SetRSA(priPath, pubPath string) (err error) {
	var (
		priPem []byte
		pubPem []byte
	)
	if priPem, err = ioutil.ReadFile(priPath); err != nil {
		return
	}
	if pubPem, err = ioutil.ReadFile(pubPath); err != nil {
		return
	}
	if k.rsaPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(priPem); err != nil {
		return
	}

	if k.rsaPublicKey, err = jwt.ParseRSAPublicKeyFromPEM(pubPem); err != nil {
		return
	}
	return
}

// SetHmac 设置HMAC密钥
func (k *Key) SetHmac(hmacKeyPath string) (err error) {
	var (
		hmacKey []byte
	)
	if hmacKey, err = ioutil.ReadFile(hmacKeyPath); err != nil {
		return
	}
	k.hmacSecret = hmacKey

	return
}

// SetECDSA 设这ECDSA
func (k *Key) SetECDSA(priPath, pubPath string) (err error) {
	var (
		priPem []byte
		pubPem []byte
	)
	if priPem, err = ioutil.ReadFile(priPath); err != nil {
		return
	}
	if pubPem, err = ioutil.ReadFile(pubPath); err != nil {
		return
	}
	if k.ecdsaPrivateKey, err = jwt.ParseECPrivateKeyFromPEM(priPem); err != nil {
		return
	}

	if k.ecdsaPublicKey, err = jwt.ParseECPublicKeyFromPEM(pubPem); err != nil {
		return
	}
	return
}

// tokenParse 解析出token中的值
func (k *Key) tokenParse(tokenStr string) (ret map[string]interface{}, err error) {
	var (
		token *jwt.Token
	)

	// 解析token
	if token, err = jwt.Parse(tokenStr, k.parseKey); err != nil {
		return
	} else if token.Valid == false {
		err = error2.ErrTokenInvalid
		return
	}

	// 解析map
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		//// 解析合法性
		if err = claims.Valid(); err != nil {
			return
		}
		ret = claims
	} else {
		err = error2.ErrTokenClaimsInvalid
		return
	}

	return
}

// TokenEncode 将map转为加密的token
func (k *Key) TokenEncode(val map[string]interface{}) (tokenStr string, err error) {
	switch k.encryptMethod {
	case "RS256":
		tokenStr, err = jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(val)).SignedString(k.rsaPrivateKey)
	case "RS384":
		tokenStr, err = jwt.NewWithClaims(jwt.SigningMethodRS384, jwt.MapClaims(val)).SignedString(k.rsaPrivateKey)
	case "RS512":
		tokenStr, err = jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims(val)).SignedString(k.rsaPrivateKey)
	case "HS256":
		tokenStr, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(val)).SignedString(k.hmacSecret)
	case "HS384":
		tokenStr, err = jwt.NewWithClaims(jwt.SigningMethodES384, jwt.MapClaims(val)).SignedString(k.hmacSecret)
	case "HS512":
		tokenStr, err = jwt.NewWithClaims(jwt.SigningMethodES512, jwt.MapClaims(val)).SignedString(k.hmacSecret)
	case "ES256":
		tokenStr, err = jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims(val)).SignedString(k.ecdsaPrivateKey)
	case "ES384":
		tokenStr, err = jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims(val)).SignedString(k.ecdsaPrivateKey)
	case "ES512":
		tokenStr, err = jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims(val)).SignedString(k.ecdsaPrivateKey)
	default:
		err = fmt.Errorf("encrypt method %s not exist", k.encryptMethod)
		return
	}

	return
}

// parseKey 按类型取出密钥
func (k *Key) parseKey(token *jwt.Token) (ret interface{}, err error) {
	switch token.Method.(type) {
	case *jwt.SigningMethodRSA:
		// RSA
		ret = k.rsaPublicKey
	case *jwt.SigningMethodHMAC:
		// HMAC
		ret = k.hmacSecret
	case *jwt.SigningMethodECDSA:
		ret = k.ecdsaPublicKey
	default:
		err = fmt.Errorf("%s '%v'", error2.ErrTokenParseSignMethodNotSupport.Error(), token.Header["alg"])
		break
	}
	return
}

// NewKey new Key 直接使用key可以一个密钥对多个用户一起使用
func NewKey(encryptMethod string) *Key {
	d := new(Key)
	d.encryptMethod = encryptMethod
	return d
}
