package main

import (
	"testing"
	"path/filepath"
	"io/ioutil"
	"os"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	u := chatUser{uniqueID:"abc"}
	url, err := authAvatar.GetAvatarURL(u)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL when no value present")
	}
	testUrl := ""
	url, err = authAvatar.GetAvatarURL(u)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error when value present")
	}
	if url != testUrl {
		t.Error("AuthAvatar.GetAvatarURL should return correct URL")
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	user := &chatUser{uniqueID:"abc"}
	url, err := gravatarAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return an error")
	}
	if url != "//www.gravatar.com/avatar/abc" {
		t.Errorf("GravatarAvatar.GetAvatarURL wrongly returned %s", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	filename := filepath.Join("avatars", "abc.png")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer os.Remove(filename)
	var fileSystemAvatar FileSystemAvatar
	user := &chatUser{uniqueID:"abc"}
	url, err := fileSystemAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURL should not return an error")
	}
	if url != "/avatars/abc.png" {
		t.Errorf("FileSystemAvatar.GetAvatarURL wrongly returned %s", url)
	}
}