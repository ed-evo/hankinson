import type { Choice, Domanda, Quesito } from '@/services/hankinson'
import { validateAnswer } from '@/utils/quesiti'
export class TrainingSettings {
  public static readonly minQuesiti = 10
  public static readonly maxQuesiti = 50
  public static readonly minErrori = 0
  public static readonly maxErrori = 10
  public static readonly minSecondi = 20
  public static readonly maxSecondi = 90
  constructor(
    public numeroQuesiti: number = 15,
    public erroriAmmessi: number = 3,
    public secondiPerDomanda: number = 30
  ) {}

  get tempoTotaleAmmesso() {
    return this.numeroQuesiti * this.secondiPerDomanda
  }
}

export class QuizItem {
  private _answer: Choice | null = null
  private _isCorrect: boolean = false
  constructor(
    public readonly quesito: Quesito,
    public readonly domanda: Domanda
  ) {}

  set answer(answer: Choice | null) {
    this._answer = answer
    if (!answer) {
      this._isCorrect = false
      return
    }
    this._isCorrect = validateAnswer(this.domanda, answer)
  }

  get answer(): Choice | null {
    return this._answer
  }

  get isAnswered(): boolean {
    return !!this._answer
  }

  get isCorrect(): boolean {
    return this._isCorrect
  }
}
