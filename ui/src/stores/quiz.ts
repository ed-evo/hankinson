import { getCapitoli, getDomandeByCapitolo, login, USER_REF, type Capitolo, type Domanda, type User } from "@/services/hankinson";
import { useLocalStorage } from "@vueuse/core";
import { defineStore } from "pinia";
import { computed, ref, type Ref } from "vue";

class MultiMap<K, V> extends Map<K, V[]> {
    // Overrides standard set to append to an array
    add(key: K, value: V) {
        let values: V[] | undefined = this.get(key)
        if (!values) {
            values = []
            super.set(key, values);
        }
        values.push(value);
        return this; // Allows chaining
    }
}

async function fetchCapitoli(
    progress: Ref<number>,
    capitoliContainer: Map<number, Capitolo>,
    domandeContainer: MultiMap<number, Domanda>
) {
    const body: Capitolo[] = await getCapitoli()
    const singlePercentage = 100/body.length
    for (const capitolo of body) {
        capitoliContainer.set(capitolo.id as number, capitolo)
        const domande = await getDomandeByCapitolo(capitolo.id)
        domandeContainer.set(capitolo.id, domande)
        progress.value += singlePercentage
    }
}

export const useQuizStore = defineStore('quiz', () => {

    const capitoliSelezionati: Ref<number[]> = useLocalStorage('quiz.capitoliSelezionati', [])
    const downloadProgress = ref(-1)
    const capitoli: Map<number, Capitolo> = new Map()
    const domandeByCapitoli: MultiMap<number, Domanda> = new MultiMap()

    login().then(
        user => console.log('user', user),
        err => {
            console.error("errore login", err)
            const user = window.prompt("Inserisci la tua email.")
            console.info("user", user)
            USER_REF.value = user as User
        }
    ).then(() => fetchCapitoli(downloadProgress, capitoli, domandeByCapitoli))
    .then(() => {
        if (capitoliSelezionati.value.length === 0) {
            capitoliSelezionati.value.push(1, 2, 3)
        }
        downloadProgress.value = 100
    })
    
    return {
        // state
        capitoliSelezionati,
        // getters
        user: computed(() => USER_REF),
        downloadProgress: computed(() => downloadProgress.value),
        capitoli: computed(() => capitoli),
        domandeByCapitoli: computed(() => domandeByCapitoli),
        // actions
    }
})