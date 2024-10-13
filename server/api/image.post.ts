export default defineEventHandler(async (event) => {

  const session = await getUserSession(event)

  console.log(session.user)

  try {
    if (!session.user) {
      console.error("User not logged in")
      return createError({
        statusCode: 401,
        statusMessage: "User is not authenticated",
      })
    }

    const body = await readBody(event)
    if (!body) {
      console.error("Request body is empty or undefined")
      return createError({
        statusCode: 400,
        statusMessage: "Request body is empty or undefined",
      })
    }

    return await insertImage(body)
  }
  catch (error) {
    console.error("Error handling image POST request:", error)
    return createError({
      statusCode: 500,
      statusMessage: "Failed to process request",
    })
  }
})

async function insertImage(imageMetadata: ImageMetadata): Promise<{ statusCode: number, statusMessage: string }> {
  const insertImageSql = `INSERT INTO metadata
    (fileName, title, dateTaken, dateUploaded, cameraModel, lensModel, aperture, shutterSpeed, flashStatus, focusLength, iso, exposureMode, whiteBalance)
    VALUES (
    "${imageMetadata.fileName}",
    "${imageMetadata.title}",
    "${imageMetadata.dateTaken}",
    "${imageMetadata.dateUploaded}",
    "${imageMetadata.cameraModel}",
    "${imageMetadata.lensModel}",
    "${imageMetadata.aperture}",
    "${imageMetadata.shutterSpeed}",
    "${imageMetadata.flashStatus}",
    "${imageMetadata.focusLength}",
    "${imageMetadata.iso}",
    "${imageMetadata.exposureMode}",
    "${imageMetadata.whiteBalance}"
    )`

  try {
    await new Promise((resolve, reject) => {
      db.run(insertImageSql, [], (err) => {
        if (err) {
          console.error(`Failed to insert image to database: ${err.message}`)
          console.info(insertImageSql)
          reject(err)
        }
        else {
          console.debug(`Inserted image ${imageMetadata.fileName}.`)
          resolve(true)
        }
      })
    })
    return {
      statusCode: 200,
      statusMessage: 'success!',
    }
  }
  catch {
    return {
      statusCode: 400,
      statusMessage: 'Failed to insert image to database',
    }
  }
}
