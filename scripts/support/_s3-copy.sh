#!/bin/bash

echo "access_key = $S3_KEY" >> /s3cfg.s3
echo "secret_key = $S3_SECRET" >> /s3cfg.s3

s3cmd -c /s3cfg.s3 cp $*
