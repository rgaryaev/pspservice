# Testset is an utility to create  a test set of original passport data file. 
# all data are validated before copy. 
>testset Mode(d or r), Num , imput_file, out_file
Num - number of records to copy
Mode - direct (d): consecutive copy from begin of file   or random (r): records from original file are copied randomly 
input_file - path to original passport data file
output_file - path to destination file

Example 


>testset r 100000 list_of_expired_passports.csv test_random_valid_100k.csv 
This example creates  a test set with 100k records   selected randomly from original file

