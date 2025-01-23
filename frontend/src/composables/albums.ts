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

export async function addImageToAlbum(albumSlug: string, imageSlug: string) {
  const newAlbum = {
    AlbumSlug: albumSlug,
    ImageSlug: imageSlug,
  }
  const options = {
    body: JSON.stringify(newAlbum),
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
  }
  const response = await backendFetchRequest('link', options)
  return response
}

export async function removeImageFromAlbum(albumSlug: string, imageSlug: string) {
  const newAlbum = {
    AlbumSlug: albumSlug,
    ImageSlug: imageSlug,
  }
  const options = {
    body: JSON.stringify(newAlbum),
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
  }
  const response = await backendFetchRequest('link', options)
  return response
}
