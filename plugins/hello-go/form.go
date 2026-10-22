package hello

import (
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/plugins/internal/acceptance"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
)

// Config documents the finite supported settings of this example. Runtime keeps
// unknown saved keys intact; they are not executable settings or form controls.
type Config struct {
	Greeting          string `json:"greeting"`
	Scheduled         bool   `json:"scheduled"`
	ChainDelay        bool   `json:"chainDelay"`
	ChainError        bool   `json:"chainError"`
	DiscoveryDelay    bool   `json:"discoveryDelay"`
	Intercept         bool   `json:"intercept"`
	ModuleTitle       string `json:"moduleTitle"`
	OrganizeDelay     bool   `json:"organizeDelay"`
	OrganizeFolder    string `json:"organizeFolder"`
	OrganizeReject    bool   `json:"organizeReject"`
	RecognizeTitle    string `json:"recognizeTitle"`
	ResourceDelay     bool   `json:"resourceDelay"`
	ResourceNoCleanup bool   `json:"resourceNoCleanup"`
	ResourceReject    bool   `json:"resourceReject"`
	ResourceReverse   bool   `json:"resourceReverse"`
}

func configurationForm() sdk.Form {
	fields := []contracts.Component{acceptance.Field("VTextField", "greeting", "问候语"), acceptance.Field("VSwitch", "scheduled", "启用定时服务")}
	diagnostic := []contracts.Component{{Component: "VAlert", Props: map[string]json.RawMessage{"title": acceptance.Prop("以下为示例插件测试选项，可能延迟或拒绝处理；生产使用请保持关闭。"), "type": acceptance.Prop("warning")}}}
	for _, field := range [][2]string{{"chainDelay", "模拟识别延迟"}, {"chainError", "模拟识别错误"}, {"discoveryDelay", "模拟发现延迟"}, {"intercept", "拦截识别结果"}, {"organizeDelay", "模拟整理延迟"}, {"organizeReject", "拒绝整理"}, {"resourceDelay", "模拟资源延迟"}, {"resourceNoCleanup", "保留下载源文件"}, {"resourceReject", "拒绝资源下载"}, {"resourceReverse", "反转资源顺序"}} {
		diagnostic = append(diagnostic, acceptance.Field("VSwitch", field[0], field[1]))
	}
	for _, field := range [][2]string{{"moduleTitle", "示例处理标题"}, {"organizeFolder", "示例整理目录"}, {"recognizeTitle", "示例识别标题"}} {
		diagnostic = append(diagnostic, acceptance.Field("VTextField", field[0], field[1]))
	}
	// Progressive disclosure is presentation-only; all controls remain expressible.
	fields = append(fields, contracts.Component{Component: "VDetails", Props: map[string]json.RawMessage{"label": acceptance.Prop("示例测试选项")}, Children: diagnostic})
	return sdk.Form{Content: fields, Defaults: json.RawMessage(`{"greeting":"Hello","scheduled":true}`)}
}
