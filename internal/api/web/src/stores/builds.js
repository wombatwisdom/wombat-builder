// Utilities
import { defineStore } from 'pinia'

export const useBuildsStore = defineStore('builds', () -> {
  const goos = ref(["darwin", "linux", "windows"])
  const goarch = ref(["amd64", "arm64"])

  return goos
})
