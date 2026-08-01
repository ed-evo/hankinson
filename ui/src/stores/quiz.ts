import { getCapitoli, login, USER_REF, type Capitolo, type User } from "@/services/hankinson";
import { TrainingSettings } from "@/types/models";
import { useLocalStorage } from "@vueuse/core";
import { defineStore } from "pinia";
import { computed, ref, type Ref } from "vue";

async function fetchCapitoli(
    capitoliContainer: Map<number, Capitolo>
) {
    const body: Capitolo[] = await getCapitoli()
    for (const capitolo of body) {
        capitoliContainer.set(capitolo.id as number, capitolo)
    }
}

export const useQuizStore = defineStore('quiz', () => {

    const trainingSettings: Ref<TrainingSettings> = useLocalStorage('quiz.trainingSettings', new TrainingSettings())
    const capitoliSelezionati: Ref<number[]> = useLocalStorage('quiz.capitoliSelezionati', [])
    const downloadProgress = ref(-1)
    const capitoli: Map<number, Capitolo> = new Map()

    login().then(
        user => console.log('user', user),
        err => {
            console.error("errore login", err)
            const user = window.prompt("Inserisci la tua email.")
            console.info("user", user)
            USER_REF.value = user as User
        }
    ).then(() => fetchCapitoli(capitoli))
    .then(() => {
        if (capitoliSelezionati.value.length === 0) {
            capitoliSelezionati.value.push(1, 2, 3)
        }
        downloadProgress.value = 100
    })
    
    return {
        // state
        capitoliSelezionati,
        trainingSettings,
        // getters
        user: computed(() => USER_REF),
        isLoading: computed(() => downloadProgress.value < 100),
        capitoli: computed(() => capitoli),
        // actions
    }
})