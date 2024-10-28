import { ref } from 'vue'
import { defineStore } from 'pinia'

export const alertStore = defineStore('alert', () => {
  const alert_state = ref<string>('');
    function $reset() {
      alert_state.value = ''
    }

    return { alert_state, $reset }
})
