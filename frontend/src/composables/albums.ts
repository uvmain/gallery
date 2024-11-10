import { backendFetchRequest } from './fetchFromBackend'

export interface Album {
  Slug: string
  Name: string
  DateCreated: string
  CoverSlug: string
}

export async function getAllAlbums(): Promise<Album[]> {
  try {
    const response = await backendFetchRequest('albums')
    const albums = await response.json() as Album[]
    return albums
  }
  catch (error) {
    console.error('Failed to fetch Albums:', error)
    return []
  }
}

export async function getAlbumCoverSlugThumbnailAddress(album: Album) {
  const imageSlug = album.CoverSlug
  return `/api/thumbnail/${imageSlug}`
}
