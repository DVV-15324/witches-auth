package utils

import (
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"strconv"
)

var Bits int64

type UID struct {
	LocalID  int64
	ObjectID int64
}

func NewUID(localID int64, objectID int64) *UID {
	return &UID{
		LocalID:  localID,
		ObjectID: objectID,
	}
}

func (u *UID) Validate() error {
	if u.LocalID < 0 {
		return fmt.Errorf("localID cannot be negative")
	}
	if u.ObjectID < 0 {
		return fmt.Errorf("objectID cannot be negative")
	}
	return nil
}

func (u *UID) ToBase58() string {
	uid := uint64(u.LocalID)<<Bits | uint64(u.ObjectID)
	buf := make([]byte, 0, 20)
	buf = fmt.Appendf(buf, "%d", uid)
	return base58.Encode(buf)
}

func DecodeFromBase58(fakeUID string) *UID {
	mask := uint64((1 << Bits) - 1)
	uidByte := base58.Decode(fakeUID)
	uid, _ := strconv.ParseUint(string(uidByte), 10, 64)
	return &UID{
		LocalID:  int64(uid >> Bits),
		ObjectID: int64(uid & mask),
	}
}
