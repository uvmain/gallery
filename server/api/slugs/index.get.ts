export default defineEventHandler(async (event) => {
  // const limit = Number.parseInt(event.context.params?.limit as string) || 10
  // const offset = Number.parseInt(event.context.params?.offset as string) || 20
  try {
    const metaData = await getSlugs()
    return metaData
  }
  catch {
    return []
  }
})

export async function getSlugs(): Promise<string[] | null> {
  return new Promise((resolve, reject) => {
    db.all('SELECT slug FROM metadata ORDER BY dateTaken DESC;', (err: Error, rows: { slug: string }[]) => {
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
          const filenames = rows.map(row => row.slug)
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
