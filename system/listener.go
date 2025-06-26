package system

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	WM_POWERBROADCAST        = 0x0218
	PBT_APMSUSPEND           = 0x0004
	PBT_APMRESUMEAUTOMATIC   = 0x0012
	PBT_APMRESUMESUSPEND     = 0x0007
	PBT_APMPOWERSTATUSCHANGE = 0x000A

	// セッション変更イベント
	WM_WTSSESSIONCHANGE        = 0x02B1
	WTS_SESSION_LOCK           = 0x7
	WTS_SESSION_UNLOCK         = 0x8
	WTS_SESSION_LOGON          = 0x5
	WTS_SESSION_LOGOFF         = 0x6
	WTS_SESSION_REMOTE_CONTROL = 0x9

	// WTSRegisterSessionNotificationのフラグ
	NOTIFY_FOR_THIS_SESSION = 0
	NOTIFY_FOR_ALL_SESSIONS = 1
)

var (
	moduser32   = windows.NewLazySystemDLL("user32.dll")
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
	modwtsapi32 = windows.NewLazySystemDLL("wtsapi32.dll")

	procCreateWindowExW  = moduser32.NewProc("CreateWindowExW")
	procDefWindowProcW   = moduser32.NewProc("DefWindowProcW")
	procDispatchMessageW = moduser32.NewProc("DispatchMessageW")
	procGetMessageW      = moduser32.NewProc("GetMessageW")
	procRegisterClassExW = moduser32.NewProc("RegisterClassExW")
	procTranslateMessage = moduser32.NewProc("TranslateMessage")
	procGetModuleHandleW = modkernel32.NewProc("GetModuleHandleW")

	// セッション通知用
	procWTSRegisterSessionNotification   = modwtsapi32.NewProc("WTSRegisterSessionNotification")
	procWTSUnRegisterSessionNotification = modwtsapi32.NewProc("WTSUnRegisterSessionNotification")
)

type MSG struct {
	Hwnd    windows.HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 }
}

type WNDCLASSEXW struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   windows.Handle
	Icon       windows.Handle
	Cursor     windows.Handle
	Background windows.Handle
	MenuName   *uint16
	ClassName  *uint16
	IconSm     windows.Handle
}

// ListenSystemEvents システムイベントをリッスンする
func ListenSystemEvents() error {
	fmt.Println("電源イベントとセッションイベントをリッスンしています...")

	// ウィンドウクラスを登録
	className, _ := syscall.UTF16PtrFromString("PowerEventListener")

	hInstance, _, _ := procGetModuleHandleW.Call(0)

	wndClass := WNDCLASSEXW{
		Size:      uint32(unsafe.Sizeof(WNDCLASSEXW{})),
		WndProc:   syscall.NewCallback(wndProc),
		Instance:  windows.Handle(hInstance),
		ClassName: className,
	}

	atom, _, err := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wndClass)))
	if atom == 0 {
		return fmt.Errorf("RegisterClassEx failed: %v", err)
	}

	// 非表示ウィンドウを作成
	windowName, _ := syscall.UTF16PtrFromString("Power Event Listener")
	hwnd, _, err := procCreateWindowExW.Call(
		0,
		uintptr(atom),
		uintptr(unsafe.Pointer(windowName)),
		0,
		0, 0, 0, 0,
		0,
		0,
		hInstance,
		0,
	)

	if hwnd == 0 {
		return fmt.Errorf("CreateWindowEx failed: %v", err)
	}

	// セッション変更通知を登録
	ret, _, err := procWTSRegisterSessionNotification.Call(
		hwnd,
		NOTIFY_FOR_THIS_SESSION,
	)
	if ret == 0 {
		log.Printf("WTSRegisterSessionNotification failed: %v", err)
	} else {
		fmt.Println("セッション変更通知を登録しました")
		defer procWTSUnRegisterSessionNotification.Call(hwnd)
	}

	// メッセージループ
	var msg MSG
	for {
		ret, _, _ := procGetMessageW.Call(
			uintptr(unsafe.Pointer(&msg)),
			0, 0, 0,
		)

		if ret == 0 {
			break
		}

		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}

	return nil
}

func wndProc(hwnd windows.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_POWERBROADCAST:
		switch wParam {
		case PBT_APMSUSPEND:
			fmt.Println("[システム] スリープモードに移行します")
		case PBT_APMRESUMEAUTOMATIC:
			fmt.Println("[システム] スリープから復帰しました")
		case PBT_APMRESUMESUSPEND:
			fmt.Println("[システム] ユーザー操作による復帰")
		case PBT_APMPOWERSTATUSCHANGE:
			fmt.Println("[システム] 電源状態が変更されました")
		default:
			fmt.Printf("[システム] 電源イベント: 0x%X\n", wParam)
		}

	case WM_WTSSESSIONCHANGE:
		sessionID := uint32(lParam)
		switch wParam {
		case WTS_SESSION_LOCK:
			fmt.Printf("[セッション] 画面がロックされました (セッションID: %d)\n", sessionID)
		case WTS_SESSION_UNLOCK:
			fmt.Printf("[セッション] 画面のロックが解除されました (セッションID: %d)\n", sessionID)
		case WTS_SESSION_LOGON:
			fmt.Printf("[セッション] ユーザーがログオンしました (セッションID: %d)\n", sessionID)
		case WTS_SESSION_LOGOFF:
			fmt.Printf("[セッション] ユーザーがログオフしました (セッションID: %d)\n", sessionID)
		case WTS_SESSION_REMOTE_CONTROL:
			fmt.Printf("[セッション] リモートコントロールステータスが変更されました (セッションID: %d)\n", sessionID)
		default:
			fmt.Printf("[セッション] セッションイベント: 0x%X (セッションID: %d)\n", wParam, sessionID)
		}
	}

	ret, _, _ := procDefWindowProcW.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam,
	)
	return ret
}
