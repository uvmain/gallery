export default defineNitroPlugin(() => {
  createDatabaseDirectory()
  createMetadataTable()
  createImagesDirectory()
  createThumbnailsDirectory()
})