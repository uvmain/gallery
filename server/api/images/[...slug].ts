import { promises as fs } from 'node:fs'
import path from 'node:path'

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
    const fileBuffer = await fs.readFile(imagePath)
    const mimeType = getMimeType(slug)
    setHeader(event, 'Content-Type', mimeType)
    return fileBuffer
  }
  catch (error) {
    throw createError({
      statusCode: 404,
      statusMessage: `File not found: ${error}`,
    })
  }
})

