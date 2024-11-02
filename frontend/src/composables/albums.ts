import { backendFetchRequest } from './fetchFromBackend'
import { getRandomSlug } from './getRandomSlug'

export interface Album {
  Slug: string
  Name: string
  DateCreated: string
  CoverSlug: string
}

export async function getAlbums(): Promise<Album[]> {
  try {
    const response = await backendFetchRequest('albums')
    let albums = await response.json() as Album[]
    albums = albums.slice(0, Math.floor(Math.random() * 8) + 1)
    return albums
  }
  catch (error) {
    console.error('Failed to fetch Albums:', error)
    return []
  }
}

export async function getAlbumCoverSlugThumbnailAddress(album: Album) {
  // const imageSlug = album.CoverSlug
  const result = await getRandomSlug(1)
  const imageSlug = result ? result[0] : ''
  // console.log(imageSlug)
  return `/api/thumbnail/${imageSlug}`
}
