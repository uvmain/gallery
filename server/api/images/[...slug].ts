import { promises as fs } from 'node:fs'
import path from 'node:path'
import sharp from 'sharp'

export default defineEventHandler(async (event) => {
  let slug = event.context.params?.slug

  if (!slug) {
    throw createError({
      statusCode: 400,
      statusMessage: 'Missing slug parameter',
    })
  }

  slug = decodeURIComponent(slug)

  const filename = await new Promise<string | null>((resolve, reject) => {
    db.get(
      `SELECT fileName FROM metadata WHERE slug = ?`,
      [slug],
      (err: Error, row: { fileName: string } | null) => {
        if (err) {
          reject(err)
        }
        else {
          resolve(row ? row.fileName : null)
        }
      }
    )
  })

  if (!filename) {
    throw createError({
      statusCode: 404,
      statusMessage: 'File not found',
    })
  }

  const imagePath = path.resolve(path.join(imagesDirectory, filename))

  try {
    const fileBuffer = await createOptimisedImage(imagePath)
    return fileBuffer
  }
  catch (error) {
    throw createError({
      statusCode: 404,
      statusMessage: `File not found: ${error}`,
    })
  }
})

export async function createOptimisedImage(imagePath: string) {
  try {
    const fileBuffer = await fs.readFile(imagePath)
    const resizedImageBuffer = await sharp(fileBuffer)
      .rotate()
      .resize({
        width: serverConfiguration.imageMaxPixels,
        height: serverConfiguration.imageMaxPixels,
        fit: 'outside'
      })
      .webp()
      .toBuffer()

    return resizedImageBuffer
  }
  catch (error) {
    throw createError({
      statusCode: 403,
      statusMessage: `${error}`,
    })
  }
}