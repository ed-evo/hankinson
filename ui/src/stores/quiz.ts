import type { Capitolo, Domanda } from "@/types/hankinson";
import { defineStore } from "pinia";
import { ref, type Ref } from "vue";

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
    capitoliContainer: Ref<Map<number, Capitolo>>,
    domandeContainer: Ref<MultiMap<number, Domanda>>
) {
    const response = await fetch('/api/v1/quiz/capitoli')
    const body: Capitolo[] = await response.json()
    const singlePercentage = 100/body.length
    for (const capitolo of body) {
        capitoliContainer.value.set(capitolo.id as number, capitolo)
        const domandeResponse = await fetch(`/api/v1/quiz/capitoli/${capitolo.id}`)
        const domandeBody = await domandeResponse.json()
        domandeContainer.value.set(capitolo.id, domandeBody.edges.domande)
        progress.value += singlePercentage
    }
    progress.value = 100
}

export const useQuizStore = defineStore('quiz', () => {
    const downloadProgress = ref(0)
    const capitoli = ref<Map<number, Capitolo>>(new Map())
    const capitoliSelezionati = ref<Capitolo[]>([])
    // capitoliSelezionati.value = Array.from(capitoli.value.keys())
    const domandeByCapitoli = ref<MultiMap<number, Domanda>>(new MultiMap())
    fetchCapitoli(downloadProgress, capitoli, domandeByCapitoli)

    return {
        downloadProgress,
        capitoli,
        domandeByCapitoli,
        capitoliSelezionati
    }
}, {
    persist: {
        key: 'quiz-store',
        pick: ['capitoliSelezionati'] 
    }
})