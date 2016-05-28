#!/bin/bash
/usr/bin/aws s3 cp s3://ft-ce-repository/amazon-ftbase/bootstrap.sh . 
bash ./bootstrap.sh
