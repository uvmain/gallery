import { useSessionStorage } from '@vueuse/core'
import { backendFetchRequest } from './fetchFromBackend'

interface AllSlugs {
  timeLastChecked: Date
  slugs: string[]
}

const allSlugs = useSessionStorage<AllSlugs>('all-slugs', {} as AllSlugs)

export async function getAllSlugs() {
  if (allSlugs.value.slugs?.length > 0) {
    const now = new Date()
    const timeLastChecked = new Date(allSlugs.value.timeLastChecked)
    const timeAgo = (now.getTime() - timeLastChecked.getTime()) / 1000
    if (timeAgo < 300) { // 5 minutes
      return allSlugs.value.slugs
    }
  }

  try {
    const response = await backendFetchRequest('slugs')
    const jsonData = await response.json() as string[]
    const newSlugs: AllSlugs = {
      timeLastChecked: new Date(),
      slugs: jsonData,
    }
    allSlugs.value = newSlugs
    return jsonData
  }
  catch (error) {
    console.error('Failed to fetch thumbnails:', error)
    return []
  }
}
