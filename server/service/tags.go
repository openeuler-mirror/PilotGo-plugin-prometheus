package service

import (
	"gitee.com/openeuler/PilotGo/sdk/common"
	"openeuler.org/PilotGo/prometheus-plugin/server/dao"
	"openeuler.org/PilotGo/prometheus-plugin/server/plugin"
)

func GetTags() {
	tag_cb := func(uuids []string) []common.Tag {
		var tags []common.Tag
		for _, uuid := range uuids {
			ok, _ := dao.IsExistTargetUUID(uuid)
			if ok {
				tag := common.Tag{
					UUID: uuid,
					Type: common.TypeOk,
					Data: "exporter",
				}
				tags = append(tags, tag)
			} else {
				tag := common.Tag{
					UUID: uuid,
					Type: common.TypeError,
					Data: "",
				}
				tags = append(tags, tag)
			}
		}
		return tags
	}
	plugin.Client.OnGetTags(tag_cb)
}
