package v1handler

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/phluxx/gogert/internal/client/ldap"
	"github.com/phluxx/gogert/internal/service/config"
	"github.com/phluxx/gogert/pkg/v1model"
	"github.com/phluxx/gogert/pkg/v1request"
	"github.com/phluxx/gogert/pkg/v1view"
)

type HttpHandler struct {
	Router *httprouter.Router
	Config *config.Config
}

func NewHttpHandler(cfg *config.Config) *HttpHandler {
	return &HttpHandler{
		Router: httprouter.New(),
		Config: cfg,
	}
}

func (h *HttpHandler) RegisterHandler() {
	h.Router.GET("/health", h.healthCheck)
	h.Router.POST("/v1/password/change", h.passwordChange)
	h.Router.POST("/v1/auth/login", h.login)
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Router.ServeHTTP(w, r)
}

func (h *HttpHandler) healthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// The ldap should look for cn=USERNAME,ou=people,dc=ewnix,dc=net, password entry is userPassword
// the currentpass will read the current user's (cn=username,ou=people,dc=ewnix,dc=net) userPassword attribute, authenticate vs it, then if successful, rewrite that userPassword attribute as a CRYPT SHA-256
// the newpass will be the new password, and the newpassconfirm will be the new password confirm
func (h *HttpHandler) passwordChange(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		passChange v1request.PasswordchangeRequest
		dec        = json.NewDecoder(r.Body)
	)
	// Unmarshal the request body into the passChange struct
	err := dec.Decode(&passChange)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// Validate the passChange struct
	err = passChange.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	//Get the users account information from jwt token in authorization header
	claims, err := h.getClaims(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	c, err := ldap.New(&h.Config.LdapConfig)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Authenticate the user against the ldap server
	err = c.Authenticate(claims.Username, passChange.OldPassword)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// Change the users password
	err = c.ChangePassword(claims.Username, passChange.NewPassword)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password Changed"))
}

func (h *HttpHandler) getClaims(r *http.Request) (*v1model.Claims, error) {
	auth := r.Header.Get("Authorization")
	claims, err := v1model.GetClaims(auth, h.Config.JwtConfig)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (h *HttpHandler) login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		login v1request.LoginRequest
		dec   = json.NewDecoder(r.Body)
	)
	// Unmarshal the request body into the login struct
	err := dec.Decode(&login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	c, err := ldap.New(&h.Config.LdapConfig)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// Authenticate the user against the ldap server
	err = c.Authenticate(login.Username, login.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// Create the jwt token
	token, err := v1model.CreateToken(login.Username, h.Config.JwtConfig)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Marshal the response struct into json
	response := v1view.LoginResponse{
		Token: token,
	}
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
}
