/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: wanghaohao <wanghaohao@kylinos.cn>
 * Date: Mon Oct 28 17:43:48 2024 +0800
 */
import { ref } from 'vue'
import { defineStore } from 'pinia'

export const alertStore = defineStore('alert', () => {
  const alert_state = ref<string>('');
    function $reset() {
      alert_state.value = ''
    }

    return { alert_state, $reset }
})
