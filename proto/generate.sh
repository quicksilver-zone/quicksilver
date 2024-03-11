PROTO_BUILDER=$1
cd proto
buf generate --template buf.gen.yaml

cd ..

cp -r github.com/quicksilver-zone/quicksilver/* ./
rm -rf github.com

swagger-combine ./docs/config.json -o ./docs/swagger.yml
rm -rf tmp-swagger-gen
