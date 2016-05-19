package main

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/nmcclain/ldap"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	keysURI      *string
	ldapServer   *string
	ldapPort     *int
	ldapUser     *string
	ldapPassword *string

	httpClient *http.Client
)

type client struct {
	ldapCon    *ldap.Conn
	httpClient *http.Client
}

func (c client) keys() (keys string, err error) {
	resp, err := c.httpClient.Get(*keysURI)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Requesting services.yaml file returned %v HTTP status\n", resp.Status)
	}

	byteKeys, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	keys = strings.TrimSpace(string(byteKeys))
	return
}

func buildUserKeyMap(keys string) (userKeys map[string]string) {
	userKeys = make(map[string]string)
	keysArray := strings.Split(keys, "\n")

	for _, key := range keysArray {
		splitKey := strings.Split(key, " ")
		user := splitKey[len(splitKey)-1]
		userKeys[user] = key
	}
	return
}

func (c *client) validateUsers(userKeys map[string]string) error {
	err := c.connectLDAP()
	if err != nil {
		log.Printf("Could not connect to LDAP: %v", err)
		return err
	}
	defer c.ldapCon.Close()

	for user, _ := range userKeys {
		if !c.userActive(user) {
			delete(userKeys, user)
		}
	}
	return nil
}

func (c *client) userActive(user string) bool {
	search := ldap.NewSearchRequest(
		"DC=AD,DC=FT,DC=COM",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(mail=%s)", user),
		[]string{"userAccountControl"},
		nil)

	searchResults, err := c.ldapCon.Search(search)
	if err != nil {
		log.Printf("LDAP search for user %v did not succeed, could not verify their key", user)
		return false
	}
	/* This should only return a single value: `userAccountControl`, which is burried in entries, attributes and values.
	So we verify there is only one thing returned, and it's value is `512` https://support.microsoft.com/en-gb/kb/305144
	*/
	for _, entry := range searchResults.Entries {
		if len(entry.Attributes) == 1 {
			if len(entry.Attributes[0].Values) == 1 {
				if entry.Attributes[0].Values[0] == "512" {
					return true
				}
			}
		}
	}
	return false
}

func (c *client) connectLDAP() (err error) {
	con, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", *ldapServer, *ldapPort))
	if err != nil {
		log.Printf("dial err: %v", err)
		return err
	}

	err = con.Bind(*ldapUser, *ldapPassword)
	if err != nil {
		log.Printf("bind err: %v", err)
		return err
	}
	c.ldapCon = con
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := client{httpClient: httpClient}

	keysString, err := c.keys()
	if err != nil {
		log.Printf("couldn't get keys")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	userKeyMap := buildUserKeyMap(keysString)
	err = c.validateUsers(userKeyMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	out := ""
	for _, key := range userKeyMap {
		out += key + "\n"
	}
	fmt.Fprintf(w, fmt.Sprintf("%v", out))
}

func main() {
	app := cli.App("ssh-ad-key-validator", "Exposes public ssh keys of AD active users.")
	keysURI = app.String(cli.StringOpt{
		Name:   "keysURI",
		Value:  "",
		Desc:   "URI of the public ssh keys",
		EnvVar: "KEYS_URI",
	})
	ldapServer = app.String(cli.StringOpt{
		Name:   "ldapServer",
		Value:  "",
		Desc:   "ldap server address",
		EnvVar: "LDAP_SERVER",
	})
	ldapPort = app.Int(cli.IntOpt{
		Name:   "ldapPort",
		Value:  0,
		Desc:   "ldap server port",
		EnvVar: "LDAP_PORT",
	})
	ldapUser = app.String(cli.StringOpt{
		Name:   "ldapUser",
		Value:  "",
		Desc:   "ldap login username",
		EnvVar: "LDAP_USER",
	})
	ldapPassword = app.String(cli.StringOpt{
		Name:   "ldapPassword",
		Value:  "",
		Desc:   "ldap login password",
		EnvVar: "LDAP_PASSWORD",
	})

	httpClient = &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 25}}

	http.HandleFunc("/authorized_keys", handler)
	http.ListenAndServe(":8080", nil)
}
