/*
Copyright 2021 The AtomCI Group Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package migrations

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/go-atomci/atomci/internal/core/pipelinemgr"
	"github.com/go-atomci/atomci/internal/dao"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/internal/models"
)

type Migration20220726 struct {
}

func (m Migration20220726) GetCreateAt() time.Time {
	return time.Date(2022, 7, 26, 0, 0, 0, 0, time.Local)
}

func (m Migration20220726) Upgrade(ormer orm.Ormer) error {
	// update component
	_ = updateComponent20220726()

	// update task tmpls
	_ = updateTaskTemplates20220726()
	return nil
}

// initComponent ..
func updateComponent20220726() error {
	components := []component{
		{
			Name: "人工卡点",
			Type: "manual",
		},
		{
			Name: "构建",
			Type: "build",
		},
		{
			Name: "部署",
			Type: "deploy",
		},
		{
			Name: "单测",
			Type: "test",
		},
	}
	for _, comp := range components {
		pipelineModel := dao.NewPipelineStageModel()
		_, err := pipelineModel.GetFlowComponentByType(comp.Type)
		if err != nil {
			if err == orm.ErrNoRows {
				component := &models.FlowComponent{
					Addons: models.NewAddons(),
					Name:   comp.Name,
					Type:   comp.Type,
					Params: comp.Params,
				}
				if err := pipelineModel.CreateFlowComponent(component); err != nil {
					log.Log.Warn("when init component, occur error: %s", err.Error())
				}
			} else {
				log.Log.Warn("when init component, occur error: %s", err.Error())
			}
		} else {
			log.Log.Debug("component type `%s` already exists, skip", comp.Type)
		}
	}
	return nil
}

func updateTaskTemplates20220726() error {

	taskTmpls := []pipelinemgr.TaskTmplReq{
		{
			Name:        "单元测试",
			Type:        "test",
			Description: "运行单元测试",
			SubTask: []pipelinemgr.SubTask{
				{
					Index: 1,
					Type:  "checkout",
					Name:  "检出代码",
				},
				{
					Index: 2,
					Type:  "run-test",
					Name:  "执行",
				},
			},
		},
	}

	pipeline := pipelinemgr.NewPipelineManager()

	for _, item := range taskTmpls {
		_, err := pipeline.GetTaskTmplByName(item.Name)
		if err != nil {
			if err == orm.ErrNoRows {
				if err := pipeline.CreateTaskTmpl(&item, "admin"); err != nil {
					log.Log.Error("when init task template, occur error: %s", err.Error())
				}
			} else {
				logs.Warn("init task template occur error: %s", err.Error())
			}
		} else {
			log.Log.Debug("init task template `%s` already exists, skip", item.Name)
		}
	}
	return nil
}
