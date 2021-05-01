package pomodoroclock

import (
	"syscall"
	"unsafe"
)

const (
	SND_SYNC uint = 0x0000 /* play synchronously (default) */
)

var (
	mmsystem      = syscall.MustLoadDLL("winmm.dll")
	sndPlaySoundA = mmsystem.MustFindProc("sndPlaySoundA")
)

// SndPlaySoundA play sound file in Windows

func SndPlaySoundA(sound string, flags uint) {
	b := append([]byte(sound), 0)
	sndPlaySoundA.Call(uintptr(unsafe.Pointer(&b[0])), uintptr(flags))
}

//Raise beep using Windows dll
func beep() {
	file := "Alarm01.wav"
	SndPlaySoundA(file, SND_SYNC)
}
