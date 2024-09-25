export default defineEventHandler(async (event) => {

  const session = await getUserSession(event)

  console.log(session.user)

  try {
    // if (!loggedIn) {
    //   console.error("User not logged in")
    //   return createError({
    //     statusCode: 403,
    //     statusMessage: "User is not authenticated",
    //   })
    // }

    const body = await readBody(event)
    if (!body) {
      console.error("Request body is empty or undefined")
      return createError({
        statusCode: 400,
        statusMessage: "Request body is empty or undefined",
      })
    }

    insertImage(body)
    return {
      statusCode: 200,
      statusMessage: "success"
    }
  }
  catch (error) {
    console.error("Error handling image POST request:", error)
    return createError({
      statusCode: 500,
      statusMessage: "Failed to process request",
    })
  }
})

function insertImage(imageMetadata: ImageMetadata) {
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
    );`
  db.all(insertImageSql, [], (err) => {
    if (err) {
      console.error(`Failed to insert image to database; ${err.message}`)
      console.info(insertImageSql)
      return { outcome: err }
    }
    console.debug(`Inserted image ${imageMetadata.fileName}.`)
    return { outcome: 'success' }
  })
}