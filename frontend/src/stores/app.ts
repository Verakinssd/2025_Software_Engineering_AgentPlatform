import { reactive } from "vue";

export const AppStore = reactive<{
    loading: boolean;
}>({
    loading: true,
})