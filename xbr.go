package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/ncruces/zenity"
	"github.com/topxeq/dlgs"

	"github.com/topxeq/tk"
	"github.com/topxeq/xie"

	"github.com/kbinani/screenshot"

	"github.com/jchv/go-webview2"
)

// global variables
var versionG string = "0.1.0"

func showInfoGUI(titleA string, formatA string, messageA ...interface{}) interface{} {
	rs, errT := dlgs.Info(titleA, fmt.Sprintf(formatA, messageA...))

	if errT != nil {
		return errT
	}

	return rs
}

func getConfirmGUI(titleA string, formatA string, messageA ...interface{}) interface{} {
	flagT, errT := dlgs.Question(titleA, fmt.Sprintf(formatA, messageA...), true)
	if errT != nil {
		return errT
	}

	return flagT
}

func showErrorGUI(titleA string, formatA string, messageA ...interface{}) interface{} {
	rs, errT := dlgs.Error(titleA, fmt.Sprintf(formatA, messageA...))
	if errT != nil {
		return errT
	}

	return rs
}

// mt $pln $guiG selectFileToSave -confirmOverwrite -title=savefile... -default=c:\test\test.txt `-filter=[{"Name":"Go and TextFiles", "Patterns":["*.go","*.txt"], "CaseFold":true}]`

func selectFileToSaveGUI(argsA ...string) interface{} {
	optionsT := []zenity.Option{}

	optionsT = append(optionsT, zenity.ShowHidden())

	titleT := tk.GetSwitch(argsA, "-title=", "")

	if titleT != "" {
		optionsT = append(optionsT, zenity.Title(titleT))
	}

	defaultT := tk.GetSwitch(argsA, "-default=", "")

	if defaultT != "" {
		optionsT = append(optionsT, zenity.Filename(defaultT))
	}

	if tk.IfSwitchExistsWhole(argsA, "-confirmOverwrite") {
		optionsT = append(optionsT, zenity.ConfirmOverwrite())
	}

	filterStrT := tk.GetSwitch(argsA, "-filter=", "")

	// tk.Plv(filterStrT)

	var filtersT zenity.FileFilters

	if filterStrT != "" {

		errT := jsoniter.Unmarshal([]byte(filterStrT), &filtersT)

		if errT != nil {
			return errT
		}

		optionsT = append(optionsT, filtersT)
	}

	rs, errT := zenity.SelectFileSave(optionsT...)

	if errT != nil {
		if errT == zenity.ErrCanceled {
			return nil
		}

		return errT
	}

	return rs
}

func selectFileGUI(argsA ...string) interface{} {
	optionsT := []zenity.Option{}

	optionsT = append(optionsT, zenity.ShowHidden())

	titleT := tk.GetSwitch(argsA, "-title=", "")

	if titleT != "" {
		optionsT = append(optionsT, zenity.Title(titleT))
	}

	defaultT := tk.GetSwitch(argsA, "-default=", "")

	if defaultT != "" {
		optionsT = append(optionsT, zenity.Filename(defaultT))
	}

	filterStrT := tk.GetSwitch(argsA, "-filter=", "")

	var filtersT zenity.FileFilters

	if filterStrT != "" {

		errT := jsoniter.Unmarshal([]byte(filterStrT), &filtersT)

		if errT != nil {
			return errT
		}

		optionsT = append(optionsT, filtersT)
	}

	rs, errT := zenity.SelectFile(optionsT...)

	if errT != nil {
		if errT == zenity.ErrCanceled {
			return nil
		}

		return errT
	}

	return rs
}

func getInputGUI(argsA ...string) interface{} {
	optionsT := []zenity.Option{}

	optionsT = append(optionsT, zenity.ShowHidden())

	titleT := tk.GetSwitch(argsA, "-title=", "")

	if titleT != "" {
		optionsT = append(optionsT, zenity.Title(titleT))
	}

	defaultT := tk.GetSwitch(argsA, "-default=", "")

	if defaultT != "" {
		optionsT = append(optionsT, zenity.EntryText(defaultT))
	}

	hideTextT := tk.IfSwitchExistsWhole(argsA, "-hideText")
	if hideTextT {
		optionsT = append(optionsT, zenity.HideText())
	}

	modalT := tk.IfSwitchExistsWhole(argsA, "-modal")
	if modalT {
		optionsT = append(optionsT, zenity.Modal())
	}

	textT := tk.GetSwitch(argsA, "-text=", "")

	okLabelT := tk.GetSwitch(argsA, "-okLabel=", "")

	if okLabelT != "" {
		optionsT = append(optionsT, zenity.OKLabel(okLabelT))
	}

	cancelLabelT := tk.GetSwitch(argsA, "-cancelLabel=", "")

	if cancelLabelT != "" {
		optionsT = append(optionsT, zenity.CancelLabel(cancelLabelT))
	}

	extraButtonT := tk.GetSwitch(argsA, "-extraButton=", "")

	if extraButtonT != "" {
		optionsT = append(optionsT, zenity.ExtraButton(extraButtonT))
	}

	rs, errT := zenity.Entry(textT, optionsT...)

	if errT != nil {
		if errT == zenity.ErrCanceled {
			return nil
		}

		if errT == zenity.ErrExtraButton {
			return fmt.Errorf("extraButton")
		}

		return errT
	}

	return rs
}

func newWindowWebView2(objA interface{}, paramsA interface{}) interface{} {
	paraArgsT, ok := paramsA.([]string)

	if !ok {
		nv2, ok := paramsA.([]interface{})

		if ok {
			paraArgsT = []string{}

			for i := 0; i < len(nv2); i++ {
				paraArgsT = append(paraArgsT, tk.ToStr(nv2[i]))
			}
		}
	}

	titleT := tk.GetSwitch(paraArgsT, "-title=", "dialog")
	widthT := tk.GetSwitch(paraArgsT, "-width=", "800")
	heightT := tk.GetSwitch(paraArgsT, "-height=", "600")
	iconT := tk.GetSwitch(paraArgsT, "-icon=", "2")
	debugT := tk.IfSwitchExists(paraArgsT, "-debug")
	centerT := tk.IfSwitchExists(paraArgsT, "-center")
	fixT := tk.IfSwitchExists(paraArgsT, "-fix")
	maxT := tk.IfSwitchExists(paraArgsT, "-max")
	minT := tk.IfSwitchExists(paraArgsT, "-min")

	if maxT {
		// windowStyleT = webview2.HintMax

		rectT := screenshot.GetDisplayBounds(0)

		widthT = tk.ToStr(rectT.Max.X)
		heightT = tk.ToStr(rectT.Max.Y)
	}

	if minT {
		// windowStyleT = webview2.HintMin

		widthT = "0"
		heightT = "0"
	}

	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     debugT,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  titleT,
			Width:  uint(tk.ToInt(widthT, 800)),
			Height: uint(tk.ToInt(heightT, 600)),
			IconId: uint(tk.ToInt(iconT, 2)), // icon resource id
			Center: centerT,
		},
	})

	if w == nil {
		return fmt.Errorf("failed to create window: %v", "N/A")
	}

	windowStyleT := webview2.HintNone

	if fixT {
		windowStyleT = webview2.HintFixed
	}

	w.SetSize(tk.ToInt(widthT, 800), tk.ToInt(heightT, 600), windowStyleT)

	var handlerT tk.TXDelegate

	handlerT = func(actionA string, objA interface{}, dataA interface{}, paramsA ...interface{}) interface{} {
		switch actionA {
		case "show":
			w.Run()
			return nil
		case "navigate":
			len1T := len(paramsA)
			if len1T < 1 {
				return fmt.Errorf("not enough paramters")
			}

			if len1T > 0 {
				w.Navigate(tk.ToStr(paramsA[0]))
			}

			return nil
		case "setHtml":
			len1T := len(paramsA)
			if len1T < 1 {
				return fmt.Errorf("not enough paramters")
			}

			if len1T > 0 {
				w.SetHtml(tk.ToStr(paramsA[0]))
			}

			return nil
		case "call", "eval":
			len1T := len(paramsA)
			if len1T < 1 {
				return fmt.Errorf("not enough paramters")
			}

			if len1T > 0 {
				w.Dispatch(func() {
					w.Eval(tk.ToStr(paramsA[0]))
				})
			}

			return nil
		case "close":
			w.Destroy()
			return nil
		case "setQuickDelegate":
			len1T := len(paramsA)
			if len1T < 1 {
				return fmt.Errorf("not enough parameters")
			}

			var deleT tk.QuickVarDelegate = paramsA[0].(tk.QuickVarDelegate)

			w.Bind("quickDelegateDo", func(args ...interface{}) interface{} {
				// args is the input parameter while WebView2 call delegate functions
				// strT := args[0].String()

				rsT := deleT(args...)

				if tk.IsErrX(rsT) {
					// if xie.GlobalsG.VerboseLevel > 0 {
					tk.Pl("error occurred in QuickVarDelegate: %v", rsT)
					// }
				}

				// must return a value
				return rsT
			})

			return nil
		case "setDelegate":
			len1T := len(paramsA)
			if len1T < 1 {
				return fmt.Errorf("not enough parameters")
			}

			var codeT = paramsA[0]

			dv1, ok := codeT.(tk.QuickVarDelegate)

			if ok {
				w.Bind("delegateDo", dv1)
				return nil
			}

			return fmt.Errorf("invalid type: %T(%v)", codeT, codeT)
		// return nil
		case "setGoDelegate":
			// var codeT string = tk.ToStr(paramsA[0])

			// p := objA.(*xie.XieVM)
			var codeT = paramsA[0]

			dv1, ok := codeT.(tk.QuickVarDelegate)

			if ok {
				w.Bind("delegateDo", dv1)
				return nil
			}

			return fmt.Errorf("invalid type: %T(%v)", codeT, codeT)
		default:
			return fmt.Errorf("unknown action: %v", actionA)
		}

		return nil
	}

	return handlerT

}

func guiHandler(actionA string, objA interface{}, dataA interface{}, paramsA ...interface{}) interface{} {
	switch actionA {
	case "init":
		// rs := initGUI()
		return ""
	case "lockOSThread":
		runtime.LockOSThread()
		return nil
	case "method", "mt":
		if len(paramsA) < 1 {
			return fmt.Errorf("not enough paramters")
		}

		objT := paramsA[0]

		methodNameT := tk.ToStr(paramsA[1])

		v1p := 2

		switch nv := objT.(type) {
		case zenity.ProgressDialog:
			switch methodNameT {
			case "close":
				rs := nv.Close()
				return rs
			case "complete":
				rs := nv.Complete()
				return rs
			case "text":
				if len(paramsA) < v1p+1 {
					return fmt.Errorf("not enough paramters")
				}

				v1 := tk.ToStr(paramsA[v1p])

				rs := nv.Text(v1)
				return rs
			case "value":
				if len(paramsA) < v1p+1 {
					return fmt.Errorf("not enough paramters")
				}

				v1 := tk.ToInt(paramsA[v1p])

				rs := nv.Value(v1)
				return rs
			case "maxValue":
				rs := nv.MaxValue()
				return rs
			case "done":
				return nv.Done()
			}
		}

		rvr := tk.ReflectCallMethod(objT, methodNameT, paramsA[2:]...)

		return rvr

	case "new":
		if len(paramsA) < 1 {
			return fmt.Errorf("not enough paramters")
		}

		vs1 := tk.ToStr(paramsA[0])

		p, ok := objA.(*xie.XieVM)

		if !ok {
			p = nil
		}

		switch vs1 {
		case "window", "webView2":
			return newWindowWebView2(p, paramsA[1:])
		}

		return fmt.Errorf("unsupported type: %v", vs1)

	case "close":
		if len(paramsA) < 1 {
			return fmt.Errorf("not enough paramters")
		}

		switch nv := paramsA[0].(type) {
		case zenity.ProgressDialog:
			nv.Close()
		}

		return ""

	case "showInfo":
		if len(paramsA) < 2 {
			return fmt.Errorf("not enough paramters")
		}
		return showInfoGUI(tk.ToStr(paramsA[0]), tk.ToStr(paramsA[1]), paramsA[2:]...)

	case "showError":
		if len(paramsA) < 2 {
			return fmt.Errorf("not enough paramters")
		}
		return showErrorGUI(tk.ToStr(paramsA[0]), tk.ToStr(paramsA[1]), paramsA[2:]...)

	case "getConfirm":
		if len(paramsA) < 2 {
			return fmt.Errorf("not enough paramters")
		}
		return getConfirmGUI(tk.ToStr(paramsA[0]), tk.ToStr(paramsA[1]), paramsA[2:]...)
	case "getInput":
		// if len(paramsA) < 2 {
		// 	return fmt.Errorf("not enough paramters")
		// }
		return getInputGUI(tk.InterfaceToStringArray(paramsA)...)
	case "selectFile":
		// if len(paramsA) < 2 {
		// 	return fmt.Errorf("not enough paramters")
		// }
		return selectFileGUI(tk.InterfaceToStringArray(paramsA)...)
	case "selectFileToSave":
		// if len(paramsA) < 2 {
		// 	return fmt.Errorf("not enough paramters")
		// }
		return selectFileToSaveGUI(tk.InterfaceToStringArray(paramsA)...)
	case "getActiveDisplayCount":
		return screenshot.NumActiveDisplays()
	case "getScreenResolution":
		var paraArgsT []string = []string{}

		for i := 0; i < len(paramsA); i++ {
			paraArgsT = append(paraArgsT, tk.ToStr(paramsA[i]))
		}

		pT := objA.(*xie.XieVM)

		formatT := pT.GetSwitchVarValue(pT.Running, paraArgsT, "-format=", "")

		idxStrT := pT.GetSwitchVarValue(pT.Running, paraArgsT, "-index=", "0")

		idxT := tk.StrToInt(idxStrT, 0)

		rectT := screenshot.GetDisplayBounds(idxT)

		if formatT == "" {
			return []interface{}{rectT.Max.X, rectT.Max.Y}
		} else if formatT == "raw" || formatT == "rect" {
			return rectT
		} else if formatT == "json" {
			return tk.ToJSONX(rectT, "-sort")
		}

		return []interface{}{rectT.Max.X, rectT.Max.Y}
	case "showProcess":
		var paraArgsT []string = []string{}

		for i := 0; i < len(paramsA); i++ {
			paraArgsT = append(paraArgsT, tk.ToStr(paramsA[i]))
		}

		optionsT := []zenity.Option{}

		titleT := tk.GetSwitch(paraArgsT, "-title=", "")

		if titleT != "" {
			optionsT = append(optionsT, zenity.Title(titleT))
		}

		okButtonT := tk.GetSwitch(paraArgsT, "-ok=", "")

		if titleT != "" {
			optionsT = append(optionsT, zenity.OKLabel(okButtonT))
		}

		cancelButtonT := tk.GetSwitch(paraArgsT, "-cancel=", "")

		if titleT != "" {
			optionsT = append(optionsT, zenity.CancelLabel(cancelButtonT))
		}

		if tk.IfSwitchExistsWhole(paraArgsT, "-noCancel") {
			optionsT = append(optionsT, zenity.NoCancel())
		}

		if tk.IfSwitchExistsWhole(paraArgsT, "-modal") {
			optionsT = append(optionsT, zenity.Modal())
		}

		if tk.IfSwitchExistsWhole(paraArgsT, "-pulsate") {
			optionsT = append(optionsT, zenity.Pulsate())
		}

		maxT := tk.GetSwitch(paraArgsT, "-max=", "")

		if maxT != "" {
			optionsT = append(optionsT, zenity.MaxValue(tk.ToInt(maxT, 100)))
		}

		dlg, errT := zenity.Progress(optionsT...)
		if errT != nil {
			return fmt.Errorf("创建进度框失败（failed to create progress dialog）：%v", errT)
		}

		return dlg

	case "newWindow":
		return newWindowWebView2(objA, paramsA)

	default:
		return fmt.Errorf("unknown method")
	}

	return ""
}

var scriptPathG = ""
var guiHandlerG tk.TXDelegate

func main() {
	guiHandlerG = guiHandler

	argsT := []interface{}{"window"}

	for _, v := range os.Args {
		argsT = append(argsT, v)
	}

	rsT := guiHandlerG("new", nil, nil, argsT...)

	if tk.IsErrX(rsT) {
		tk.Pl("failed to init window: %v", rsT)
		return
	}

	windowT := rsT.(tk.TXDelegate)

	urlT := tk.GetSwitch(os.Args, "-url=", "")

	if urlT != "" {
		windowT("navigate", nil, nil, urlT)
	}

	htmlT := tk.GetSwitch(os.Args, "-html=", "")

	if htmlT != "" {
		if htmlT == "clip" {
			htmlT = tk.GetClipText()
		} else if strings.HasPrefix(htmlT, "file:") {
			htmlT = tk.LoadStringFromFile(htmlT[5:])
		} else if strings.HasPrefix(htmlT, "http") {
			htmlT = tk.ToStr(tk.GetWeb(htmlT))
		}

		windowT("setHtml", nil, nil, htmlT)
	}

	xieFilePathT := tk.GetSwitch(os.Args, "-xie=", "./default.xie")
	var scriptT string

	if strings.HasPrefix(xieFilePathT, "http") {
		rsT := tk.DownloadWebPageX(xieFilePathT)
		scriptPathG = xieFilePathT

		if tk.IsErrStr(rsT) {
			scriptT = ""
		} else {
			scriptT = rsT
		}

		if tk.IsErrX(scriptT) {
			tk.Pl("failed to load script: %v", tk.GetErrStrX(scriptT))
			tk.Exit()
		}
	} else if tk.IfFileExists(xieFilePathT) {
		scriptT = tk.LoadStringFromFile(xieFilePathT)
		scriptPathG = xieFilePathT

		if tk.IsErrX(scriptT) {
			tk.Pl("failed to load script: %v", tk.GetErrStrX(scriptT))
			tk.Exit()
		}

		if strings.HasPrefix(scriptT, "//TXDEF#") {
			scriptT = tk.TKX.DecryptStringByTXDEF(scriptT)

			if tk.IsErrStrX(scriptT) {
				tk.Fatalf("invalid code")
			}
		}

		rs := xie.RunCode(scriptT, nil, map[string]interface{}{"guiG": guiHandlerG, "windowG": windowT, "versionG": versionG, "scriptPathG": scriptPathG}, os.Args...)
		if !tk.IsUndefined(rs) {
			tk.Pl("%v", rs)
		}

		return
	} else {
		windowT("show", nil, nil)
	}
}
