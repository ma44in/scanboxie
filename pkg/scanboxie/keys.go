package scanboxie

import (
	"scanboxie/pkg/device"
	"strings"
)

// Map device keys to actual chars
//
// Howto build this map automatically:
// 	- for char in {a..z}; do; echo "device.KEY_$(echo ${char} | tr '[:lower:]' '[:upper:]'): '${char}',"; done
// 	- for char in {0..9}; do; echo "device.KEY_$(echo ${char} | tr '[:lower:]' '[:upper:]'): '${char}',"; done
var keyToCharMap = map[int]byte{
	device.KEY_A:     'a',
	device.KEY_B:     'b',
	device.KEY_C:     'c',
	device.KEY_D:     'd',
	device.KEY_E:     'e',
	device.KEY_F:     'f',
	device.KEY_G:     'g',
	device.KEY_H:     'h',
	device.KEY_I:     'i',
	device.KEY_J:     'j',
	device.KEY_K:     'k',
	device.KEY_L:     'l',
	device.KEY_M:     'm',
	device.KEY_N:     'n',
	device.KEY_O:     'o',
	device.KEY_P:     'p',
	device.KEY_Q:     'q',
	device.KEY_R:     'r',
	device.KEY_S:     's',
	device.KEY_T:     't',
	device.KEY_U:     'u',
	device.KEY_V:     'v',
	device.KEY_W:     'w',
	device.KEY_X:     'x',
	device.KEY_Y:     'y',
	device.KEY_Z:     'z',
	device.KEY_0:     '0',
	device.KEY_1:     '1',
	device.KEY_2:     '2',
	device.KEY_3:     '3',
	device.KEY_4:     '4',
	device.KEY_5:     '5',
	device.KEY_6:     '6',
	device.KEY_7:     '7',
	device.KEY_8:     '8',
	device.KEY_9:     '9',
	device.KEY_ENTER: '\n',
}

func getCharFromKeyEvent(ev *device.InputEvent, shiftPressed bool) string {
	if ev.Type != device.EV_KEY {
		return ""
	}

	val, haskey := keyToCharMap[int(ev.Code)]
	if haskey {
		if shiftPressed {
			return strings.ToUpper(string(val))
		}

		return string(val)
	}

	return "?"
}
