import { Choice, type Domanda } from '@/services/hankinson'

export function booleanToChoice (value?: boolean): Choice | null {
  if (value === undefined || value === null) {
    return null
  }
  return value ? Choice.VERO : Choice.FALSO
}

export function validateAnswer (domanda: Domanda, choice: Choice): boolean {
  switch (choice) {
    case Choice.VERO: {
      return domanda.is_true
    }
    case Choice.FALSO: {
      return !domanda.is_true
    }
  }
}
