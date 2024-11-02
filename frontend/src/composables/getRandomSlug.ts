import { backendFetchRequest } from './fetchFromBackend'

export async function getRandomSlug(limit: number) {
  try {
    const response = await backendFetchRequest(`slugs/random?limit=${limit}`)
    const jsonData = await response.json() as string[]
    console.log(jsonData)
    return jsonData
  }
  catch (error) {
    console.error('Failed to fetch thumbnails:', error)
  }
}
