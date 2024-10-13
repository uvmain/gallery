export default defineNitroPlugin(async() => {
  await createDatabaseDirectory()
  await createImagesDirectory()
  await createThumbnailsDirectory()
  await getImageDirectoryListing()
  await createMetadataTable()
  await createMetaDataForAllImages()
  await createThumbnailsForAllImages()
})
