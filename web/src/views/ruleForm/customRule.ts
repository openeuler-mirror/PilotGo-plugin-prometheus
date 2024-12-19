/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: wanghaohao <wanghaohao@kylinos.cn>
 * Date: Mon Oct 28 17:58:01 2024 +0800
 */
// 自定义校验

// 校验阈值是否在0-100之间
export const checkThreshold = (_rule: any, value: any, callback: any) => {
  setTimeout(() => {
    if (!Number.isInteger(parseInt(value))) {
      callback(new Error('请输入数字'))
    } else {
      if (value < 0 || value > 100) {
        callback(new Error('请输入0-100之间的数字'))
      } else {
        callback()
      }
    }
  }, 100)
}

export const checkDuration = (_rule: any, value: any, callback: any) => {
  setTimeout(() => {
    if (isNaN(value)) {
      callback(new Error('请输入数字'))
    } else {
      if (value < 0) {
        callback(new Error('请输入大于等于0的数字'))
      } else {
        callback()
      }
    }
  }, 100)
}