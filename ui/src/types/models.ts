
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