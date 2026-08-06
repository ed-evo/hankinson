import { Choice, type Domanda } from '@/services/hankinson'

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
