package script

import (
	"errors"
	"github.com/hiboma/suexec"
	"os"
	"strings"
	"syscall"
)

type Script struct {
	path, cwd           string
	path_info, cwd_info os.FileInfo
}

func NewScript(path string, cwd string) (*Script, error) {

	cwd_info, err := os.Lstat(cwd)
	if err != nil || !cwd_info.IsDir() {
		return nil, errors.New("cannot stat directory: (%s)\n")
	}

	if !cwd_info.IsDir() {
		return nil, errors.New("cannot stat program: (%s) %s %s\n")
	}

	path_info, err := os.Lstat(path)
	if err != nil {
		return nil, errors.New("cannot stat program: (%s) %s %s\n")
	}

	if path_info.Mode()&os.ModeSymlink != 0 {
		return nil, errors.New("cannot stat program: (%s) %s %s\n")
	}

	return &Script{path: path, cwd: cwd, path_info: path_info, cwd_info: cwd_info}, nil
}

func (self *Script) VerifyToSuexec(uid int, gid int) *suexec.Error {

	if !self.HasSecurePath() {
		return suexec.NewError(104, "invalid command (%s)\n", self.path)
	}

	if self.IsDirWritableByOthers() {
		return suexec.NewError(116, "directory is writable by others: (%s)\n", self.cwd)
	}

	if self.IsWritableByOthers() {
		return suexec.NewError(118, "file is writable by others: (%s/%s)\n", self.cwd, self.path)
	}

	if self.IsSetuid() || self.IsSetgid() {
		return suexec.NewError(119, "file is either setuid or setgid: (%s/%s)\n", self.path, self.cwd)
	}

	if !self.IsExecutable() {
		return suexec.NewError(121, "file has no execute permission: (%s/%s)\n", self.cwd, self.path)
	}

	if !self.IfOwnerMatch(uid, gid) {
		return suexec.NewError(121, "target uid/gid (%d/%d) mismatch with directory (%d/%d) or program (%d/%d)\n",
			uid, gid,
			self.path_info.Sys().(*syscall.Stat_t).Uid,
			self.path_info.Sys().(*syscall.Stat_t).Gid,
			self.cwd_info.Sys().(*syscall.Stat_t).Uid,
			self.cwd_info.Sys().(*syscall.Stat_t).Gid)
	}

	return nil
}

func (self *Script) HasSecurePath() bool {
	return !self.hasAbsolutePath() && !self.hasRelativePath()
}

func (self *Script) hasAbsolutePath() bool {
	return strings.HasPrefix(self.path, "/")
}

func (self *Script) hasRelativePath() bool {
	return strings.HasPrefix(self.path, "../") || strings.Index(self.path, "/../") > 0
}

func (self *Script) IsSetuid() bool {
	return self.path_info.Mode()&os.ModeSetuid != 0
}

func (self *Script) IsSetgid() bool {
	return self.path_info.Mode()&os.ModeSetgid != 0
}

func (self *Script) IsExecutable() bool {
	return self.path_info.Sys().(*syscall.Stat_t).Mode&syscall.S_IXUSR != 0
}

func (self *Script) IsWritableByOthers() bool {
	return self.path_info.Sys().(*syscall.Stat_t).Mode&syscall.S_IWOTH != 0 ||
		self.path_info.Sys().(*syscall.Stat_t).Mode&syscall.S_IWGRP != 0
}

func (self *Script) IsDirWritableByOthers() bool {
	return self.cwd_info.Sys().(*syscall.Stat_t).Mode&syscall.S_IWOTH != 0 ||
		self.cwd_info.Sys().(*syscall.Stat_t).Mode&syscall.S_IWGRP != 0
}

func (self *Script) IfOwnerMatch(uid int, gid int) bool {

	if uint32(uid) != self.cwd_info.Sys().(*syscall.Stat_t).Uid ||
		uint32(gid) != self.cwd_info.Sys().(*syscall.Stat_t).Gid ||
		uint32(uid) != self.path_info.Sys().(*syscall.Stat_t).Uid ||
		uint32(gid) != self.path_info.Sys().(*syscall.Stat_t).Gid {
		return false
	}

	return true
}
