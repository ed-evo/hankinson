import {
  addMilliseconds,
  formatDuration,
  interval,
  intervalToDuration,
} from 'date-fns'

export function formatDurationMs (millis: number) {
  const strat = new Date()
  const end = addMilliseconds(strat, millis)
  const i = interval(strat, end)
  const duration = intervalToDuration(i)
  return formatDuration(duration)
}

export function formatDurationMin (minutes: number) {
  return formatDurationMs(minutes * 60 * 1000)
}
