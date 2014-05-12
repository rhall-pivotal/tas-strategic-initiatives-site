echo "Building .pivotal from ${PRODUCT_DIR}"

echo "------------- metadata_parts/binaries.yml -------------"
cat ${PRODUCT_DIR}/metadata_parts/binaries.yml
echo "------------- metadata_parts/binaries.yml -------------"

rm -f ${PRODUCT_DIR}/*.pivotal
rm -f ${PRODUCT_DIR}/*.pivotal.yml
rm -f ${PRODUCT_DIR}/*.pivotal.md5

METADATA_FILE=${PRODUCT_DIR}/metadata/cf.yml

bundle install
bundle exec vara-build-metadata     --product-dir="${PRODUCT_DIR}"
bundle exec vara-download-artifacts --product-metadata="${METADATA_FILE}"
bundle exec vara-build-pivotal      --product-metadata="${METADATA_FILE}" --rc="-build${BUILD_NUMBER:--local}"
