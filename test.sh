file=listing_0038_many_register_mov
go run decoder/main.go $file

echo "Assemblying resulting file..."
./nasm decoder/result.asm
echo "Successful!"

echo "Assemblying original file..."
./nasm listings/$file.asm
echo "Succesfull!"

echo "Comparing both binaries..."
cmp listings/$file decoder/result
