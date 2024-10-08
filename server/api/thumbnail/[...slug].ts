import { promises as fs } from 'node:fs'
import path from 'node:path'
import sharp from 'sharp'

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
    try {
      const imagePath = path.resolve(imagesDirectory, slug)
      const fileBuffer = await fs.readFile(imagePath)

      const resizedImageBuffer = await sharp(fileBuffer)
        .resize({ 
          width: serverConfiguration.thumbnailMaxPixels,
          height: serverConfiguration.thumbnailMaxPixels,
          fit: 'inside'
        })
        .webp()
        .toBuffer()
      
      // Save the thumbnail as a webp file in the thumbnails directory
      console.info(`Saving thumbnail to file system: ${slug}`)
      await fs.writeFile(thumbnailPath, resizedImageBuffer)

      setHeader(event, 'Content-Type', 'image/webp')
      return resizedImageBuffer
    }
    catch(error) {
      throw createError({
        statusCode: 403,
        statusMessage: `${error}`,
      })
    }
  }
})
