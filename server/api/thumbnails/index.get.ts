export default defineEventHandler(async () => {
  try {
    const metaData = await getThumbnailPaths()
    return metaData
  }
  catch {
    return []
  }
})

export async function getThumbnailPaths(): Promise<string[] | null> {
  return new Promise((resolve, reject) => {
    db.all('SELECT filename FROM metadata ORDER BY dateTaken DESC;', (err: Error, rows: { fileName: string }[]) => {
      if (err) {
        console.error(`Failed to retrieve all filenames; ${err.message}`)
        reject({
          statusCode: 500,
          statusMessage: `Error retrieving filenames: ${err.message}`,
        })
      }
      else {
        if (rows && rows.length > 0) {
          console.info('Retrieved all filenames from database')
          const filenames = rows.map(row => `/api/thumbnail/${row.fileName}`)
          resolve(filenames)
        }
      else {
          console.warn('No metadata found')
          resolve(null)
        }
      }
    })
  })
}
