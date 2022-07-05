package wgmanager

import "golang.zx2c4.com/wireguard/wgctrl/wgtypes"

func ParseKey(key string) (parsedKey wgtypes.Key, err error) { //parses string into key
	parsedKey, err = wgtypes.ParseKey(key)
	return
}
