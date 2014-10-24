package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var client *http.Client

func StartSession(user, pass string) error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	client = &http.Client{Jar: jar}

	r, err := client.Get("https://ting.com/account/login")
	if err != nil {
		return err
	}
	if r.StatusCode != 200 {
		return fmt.Errorf("Bad status code while creating session: %s", r.Status)
	}

	vals := make(url.Values)
	vals["first_name"] = []string{""}
	vals["last_name"] = []string{""}
	vals["phone"] = []string{""}
	vals["email"] = []string{user}
	vals["password"] = []string{pass}
	vals["confirm_password"] = []string{""}
	vals["send_news"] = []string{"on"}
	vals["send_device_alerts"] = []string{"on"}
	vals["existing_user_login"] = []string{""}
	r, err = client.PostForm("https://ting.com/account/login", vals)
	if err != nil {
		return err
	}
	if r.StatusCode != 200 {
		return fmt.Errorf("Bad status code while logging in session: %s", r.Status)
	}
	return nil
}

func EndSession() error {
	r, err := client.Get("https://ting.com/account/logout")
	if err != nil {
		return err
	}
	if r.StatusCode != 200 {
		return fmt.Errorf("Bad status code while ending session: %s", r.Status)
	}
	return nil
}
