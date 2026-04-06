package configure

import (
	"fmt"

	"github.com/gwuhaolin/livego/utils/uid"
	"github.com/puzpuzpuz/xsync/v3"

	log "github.com/sirupsen/logrus"
)

type RoomKeysType struct {
	channelToKey *xsync.MapOf[string, string]
	keyToChannel *xsync.MapOf[string, string]
}

var RoomKeys = RoomKeysType{
	channelToKey: xsync.NewMapOf[string, string](),
	keyToChannel: xsync.NewMapOf[string, string](),
}

// set/reset a random key for channel
func (r *RoomKeysType) SetKey(channel string) (key string, err error) {
	for {
		key = uid.RandStringRunes(48)
		if _, found := r.keyToChannel.Load(key); !found {
			r.channelToKey.Store(channel, key)
			r.keyToChannel.Store(key, channel)
			break
		}
	}
	return
}

func (r *RoomKeysType) GetKey(channel string) (newKey string, err error) {
	var key interface{}
	var found bool
	if key, found = r.channelToKey.Load(channel); found {
		return key.(string), nil
	}
	newKey, err = r.SetKey(channel)
	log.Debugf("[KEY] new channel [%s]: %s", channel, newKey)
	return
}

func (r *RoomKeysType) GetChannel(key string) (channel string, err error) {
	chann, found := r.keyToChannel.Load(key)
	if found {
		return chann, nil
	} else {
		return "", fmt.Errorf("%s does not exists", key)
	}
}

func (r *RoomKeysType) DeleteChannel(channel string) bool {
	key, ok := r.channelToKey.Load(channel)
	if ok {
		r.channelToKey.Delete(channel)
		r.keyToChannel.Delete(key)
		return true
	}
	return false
}

func (r *RoomKeysType) DeleteKey(key string) bool {
	channel, ok := r.keyToChannel.Load(key)
	if ok {
		r.channelToKey.Delete(channel)
		r.keyToChannel.Delete(key)
		return true
	}
	return false
}
