package notifieer

import (
	"NoticeServices/app/dao"
	"NoticeServices/app/model"
	"NoticeServices/library/tools"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"plugin"
)

type instance struct {
}

var Instance = new(instance)

//GateWaySend 通过发送通道进行发送
func (n *instance) GateWaySend(message *model.InfoData) {

	//获取指定通知的配置信息
	config, _ := n.getInfoConfig(message.ConfigId)
	glog.Info("信息配置：", config)
	sendGatewayList := gstr.Explode("|", config.Config.SendGateway)
	if sendGatewayList == nil {
		return
	}

	messageBaseBody := message.MsgBody
	for _, gatewayName := range sendGatewayList {
		message.MsgBody = messageBaseBody

		//获取发送通道的通知模板
		where := g.Map{
			"config_id":    message.ConfigId,
			"send_gateway": gatewayName,
		}
		template, err := dao.Template.FindOne(where)
		if template != nil {
			paramDataMap := gconv.Map(message)
			message.MsgBody = tools.StringLiteralTemplate(template.Content, paramDataMap)

		}

		// 加载插件
		pluginPath := g.Config().GetString("system.PluginPath")
		filename := pluginPath + "/" + gatewayName + "/" + gatewayName + ".so"
		p, err := plugin.Open(filename)
		if err != nil {
			glog.Error(err)
			return
		}

		// 查找插件里的指定函数
		symbol, err := p.Lookup("Send")
		if err != nil {
			panic(err)
		}
		sendFunc, ok := symbol.(func(map[string]interface{}, *model.InfoData))

		if !ok {
			glog.Error(gerror.New("Plugin has no Send function"))
			return
		}

		sendParam := make(map[string]interface{})

		// 调用插件函数
		if template != nil {
			sendParam["code"] = template.Code
		}

		sendFunc(sendParam, message)
	}
}

//getInfoConfig 读取通知信息的配置文件
func (n *instance) getInfoConfig(configId string) (*model.EntityConfig, error) {

	var entityConfig = new(model.EntityConfig)
	err := dao.Config.Fields("*").Where(dao.Config.Columns.Id, configId).
		Scan(&entityConfig.Config)
	if err != nil {
		return nil, err
	}

	err = dao.Template.Fields("*").Where(dao.Template.Columns.ConfigId, configId).
		Scan(&entityConfig.Template)
	if err != nil {
		return nil, err
	}

	return entityConfig, nil
}
