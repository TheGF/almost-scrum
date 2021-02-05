package core

import (
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"github.com/jtblin/go-ldap-client"
)


var ldapClient ldap.LDAPClient
func InitLdap(ldapConfig *LdapConfig) {
	ldapClient :=  &ldap.LDAPClient{
		Base:         "dc=example,dc=com",
		Host:         "ldap.example.com",
		Port:         389,
		UseSSL:       false,
		BindDN:       "uid=readonlysuer,ou=People,dc=example,dc=com",
		BindPassword: "readonlypassword",
		UserFilter:   "(uid=%s)",
		GroupFilter: "(memberUid=%s)",
		Attributes:   []string{"givenName", "sn", "mail", "uid"},
	}
	defer ldapClient.Close()
}

func AuthenticateWithLdap(user, password string) bool {
	ok, _, err := ldapClient.Authenticate(user, password)
	return ok && err != nil
}

// SetPassword add a user with password to the global configuration.
func SetPassword(user, password string) error {
	config := LoadConfig()

	if password == "" {
		delete(config.Passwords, user)
		SaveConfig(config)
		return nil
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		logrus.Errorf("SetUser - Cannot save user %s and password: %v", user, err)
		return err
	}
	config.Passwords[user] = hex.EncodeToString(bytes)
	SaveConfig(config)
	logrus.Debugf("SetPassword - set password for user %s", user)
	return nil
}

//CheckUser checks if a user has expected password
func CheckUser(user, password string) bool {
	config := LoadConfig()
	hash, _ := hex.DecodeString(config.Passwords[user])

	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}
