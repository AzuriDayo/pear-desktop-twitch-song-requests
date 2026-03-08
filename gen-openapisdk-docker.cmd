rem Script to create the proper openapi go client-side sdk
curl -fsSL -o pear3.11.0-oapi3.1.yaml http://127.0.0.1:26538/doc
bunx @apiture/openapi-down-convert@0.14.2 -i ./pear3.11.0-oapi3.1.yaml -o ./pear3.11.0-oapi3.0.yaml
docker run --rm -v "%cd%:/local" openapitools/openapi-generator-cli:v7.20.0 generate -i /local/pear3.11.0-oapi3.0.yaml -g go -o /local/gen/pearsdk --package-name pearsdk --additional-properties=basePath=http://127.0.0.1:26538 --skip-validate-spec
