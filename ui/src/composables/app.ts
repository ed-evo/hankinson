import { createEventHook, useDocumentVisibility } from "@vueuse/core";
import { watch } from "vue";

export interface VisibilitySession {
    state: DocumentVisibilityState,
    startedAt: Date,
    finishedAt: Date,
}

class VisibilitySessionBuilder {
    constructor(
        private state: DocumentVisibilityState,
        private startedAt: Date = new Date()
    ) { }

    build (finishedAt: Date = new Date()): VisibilitySession {
        return {
            state: this.state,
            startedAt: this.startedAt,
            finishedAt
        }
    }

    sameState(state: DocumentVisibilityState): boolean {
        return state === this.state
    }
}

const visibility = useDocumentVisibility()

const appVisibilityHook = createEventHook<VisibilitySession>()

const initialState = visibility.value || "visible"

let currentBuilder: VisibilitySessionBuilder = new VisibilitySessionBuilder(initialState)

watch(visibility, (current) => {
    if (!current || currentBuilder.sameState(current)) {
        return
    }
    const session = currentBuilder.build()
    currentBuilder = new VisibilitySessionBuilder(current, session.finishedAt)
    appVisibilityHook.trigger(session)
})
export function  useAppVisibility() {
    return {
        onVisibilityChange: appVisibilityHook.on,
    }
}