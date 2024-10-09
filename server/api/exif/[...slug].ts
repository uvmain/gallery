export default defineEventHandler(async (event) => {
  const slug = event.context.params?.slug

  if (!slug) {
    throw createError({
      statusCode: 400,
      statusMessage: 'Missing slug parameter',
    })
  }

  try {
    const exifData = await getExifForImage(slug)
    return exifData
  }
  catch (error) {
    throw createError({
      statusCode: 404,
      statusMessage: `${error}`,
    })
  }
})