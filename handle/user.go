package handle

import(
	dblayer "filestore-server/db"
	"filestore-server/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const(
	pwd_salt = "*#890"
)

// 处理用户注册请求
func SignupHandler(w http.ResponseWriter,r *http.Request){
	if r.Method == http.MethodGet{
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()

	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	if len(username)<3 || len(passwd)<5 {
		w.Write([]byte("Invalid parameter"))
		return
	}

	enc_passwd:= util.Sha1([]byte(passwd+pwd_salt))
	suc:= dblayer.UserSignup(username,enc_passwd)
	if suc{
		w.Write([]byte("SUCCESS"))
	}else{
		w.Write([]byte("FAILED"))
	}
}

// 登录接口
func SignInHandler(w http.ResponseWriter,r *http.Request) {
	if r.Method == http.MethodGet{
		data, err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username:=r.Form.Get("username")
	password:=r.Form.Get("password")
	encPasswd:=util.Sha1([]byte(password+pwd_salt))

	// 1.检验用户名和密码
	pwdChecked:=dblayer.UserSignIn(username,encPasswd)

	if !pwdChecked{
		w.Write([]byte("FAILED1"))
		return
	}

	// 2.生成访问凭证(token)
	token:=GenToken(username)
	upRes := dblayer.UpdateToken(username,token)
	if!upRes{
		w.Write([]byte("FAILED2"))
		return
	}
	// 3.登录成功后重定向到首页
	w.Write([]byte("http://"+r.Host+"/static/view/home.html"))
	//resp := util.RespMsg{
	//	Code:0,
	//	Msg:"OK",
	//	Data: struct {
	//		Location string
	//		Username string
	//		Token string
	//	}{
	//		Location:"http://"+r.Host+"/static/view/home.html",
	//		Username:username,
	//		Token:token,
	//	},
	//}
	//w.Write(resp.JSONBytes())
}

func GenToken(username string) string{
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts:=fmt.Sprintf("%s",time.Now().Unix())
	tokenPrefix:=util.MD5([]byte(username+ts+"_tokensalt"))
	return tokenPrefix+ts[:8]
}