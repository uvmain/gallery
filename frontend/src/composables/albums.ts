import { backendFetchRequest, getCachedServerUrl } from './fetchFromBackend'

export interface Album {
  Slug: string
  Name: string
  DateCreated: string
  CoverSlug: string
}

export async function getAlbums(): Promise<Album[]> {
  try {
    const response = await backendFetchRequest('albums')
    const albums = await response.json() as Album[]
    console.log(albums)
    return albums
  }
  catch (error) {
    console.error('Failed to fetch Albums:', error)
    return []
  }
}

export function getAlbumCoverSlugThumbnailAddress(album: Album) {
  const imageSlug = album.CoverSlug
  const serverBaseUrl = getCachedServerUrl()
  return `${serverBaseUrl}/api/thumbnail/${imageSlug}`
}
