import type { Domanda, RispostaEnum } from '@/types/hankinson'

export class TrainingSettings {
  public static readonly minQuesiti = 10
  public static readonly maxQuesiti = 50
  public static readonly minErrori = 0
  public static readonly maxErrori = 10
  public static readonly minSecondi = 20
  public static readonly maxSecondi = 90
  constructor (
    public numeroQuesiti = 15,
    public erroriAmmessi = 3,
    public secondiPerDomanda = 30,
  ) {}
}

export class QuizItem {
  private _answer: RispostaEnum | null = null
  private _isCorrect = false
  constructor (
    public readonly quesitoId: number,
    public readonly domanda: Domanda,
  ) {}

  get answer (): RispostaEnum | null {
    return this._answer
  }

  get isAnswered (): boolean {
    return !!this._answer
  }

  get isCorrect (): boolean {
    return this._isCorrect
  }

  set answer (answer: RispostaEnum | null) {
    this._answer = answer
    if (!answer) {
      this._isCorrect = false
      return
    }
    this._isCorrect = this.domanda.rispostaCorretta === answer
  }
}
