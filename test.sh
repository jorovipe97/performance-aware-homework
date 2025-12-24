file=listing_0038_many_register_mov
go run decoder/main.go $file

echo "Assemblying resulting file..."
./nasm decoder/result.asm
echo "Successful!"

echo "Assemblying original file..."
./nasm decoder/$file.asm
echo "Succesfull!"

echo "Comparing both binaries..."
cmp decoder/$file decoder/result
