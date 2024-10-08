import { promises as fs } from 'node:fs'
import path from 'node:path'

const imagesPath = path.resolve(serverConfiguration.imagePath)

export default defineEventHandler(async (event) => {
  const slug = event.context.params?.slug

  if (!slug) {
    throw createError({
      statusCode: 400,
      statusMessage: 'Missing slug parameter',
    })
  }

  const imagePath = path.resolve(path.join(imagesPath, slug))

  try {
    const fileBuffer = await fs.readFile(imagePath)
    const mimeType = getMimeType(slug)
    setHeader(event, 'Content-Type', mimeType)
    return fileBuffer
  }
  catch (error) {
    throw createError({
      statusCode: 404,
      statusMessage: `${error}`,
    })
  }
})
