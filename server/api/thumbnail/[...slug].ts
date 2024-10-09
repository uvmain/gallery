import { promises as fs } from 'node:fs'
import path from 'node:path'

export default defineEventHandler(async (event) => {
  const slug = event.context.params?.slug

  if (!slug) {
    throw createError({
      statusCode: 400,
      statusMessage: 'Missing slug parameter',
    })
  }

  const thumbnailPath = path.resolve(thumbnailsPath, slug)

  try {
    const fileBuffer = await fs.readFile(thumbnailPath)
    console.info(`Returning thumbnail from file system: ${slug}`)
    setHeader(event, 'Content-Type', 'image/webp')
    return fileBuffer
  }
  catch {
      const resizedImageBuffer = await createThumbnail(slug)
      setHeader(event, 'Content-Type', 'image/webp')
      return resizedImageBuffer
    }
})
