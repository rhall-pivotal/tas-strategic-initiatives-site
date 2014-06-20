echo "Building .pivotal from ${PRODUCT_DIR}"

echo "------------- metadata_parts/binaries.yml -------------"
cat ${PRODUCT_DIR}/metadata_parts/binaries.yml
echo "------------- metadata_parts/binaries.yml -------------"

rm -f ${PRODUCT_DIR}/*.pivotal
rm -f ${PRODUCT_DIR}/*.pivotal.yml
rm -f ${PRODUCT_DIR}/*.pivotal.md5

bundle install
bundle exec vara build-pivotal ${PRODUCT_DIR} ${VARA_EXTRA_FLAGS}
