/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Jul 26 16:42:38 2023 +0800
 */
import { defineStore } from 'pinia';

export const useMacStore = defineStore('mac', {
  state: () => {
    return {
      macIp: '',
    };
  },
  getters: {
    newIp(state) {
      return state.macIp.length > 0 ? state.macIp.split(':')[0] : '';
    },
  },
  actions: {
    setMacIp(ip: string) {
      this.macIp = ip;
      this.newIp;
    }
  }
});
