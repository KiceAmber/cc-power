import {defineStore} from 'pinia'
import { ref } from 'vue'

export const useStatusStore = defineStore('status', () => {
    const status = ref("");

    const setStatus = (newStatus: string) => {
        status.value = newStatus;
    }

    return {
        status,
        setStatus,
    }
})