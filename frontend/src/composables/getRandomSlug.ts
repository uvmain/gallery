import { getServerUrl } from './getServerBaseUrl'

export async function getRandomSlug(limit: number) {
  try {
    const serverBaseUrl = await getServerUrl()
    const response = await fetch(`${serverBaseUrl}/api/slugs/random?limit=${limit}`)
    const jsonData = await response.json() as string[]
    return jsonData
  }
  catch (error) {
    console.error('Failed to fetch thumbnails:', error)
  }
}
