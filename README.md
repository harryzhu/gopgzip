# gopgzip

## build
`./build.sh`

## zip a file Parallely
```
./gopgzip zip --input=/home/harryzhu/test/abc.tar --output=/the/path/of/the/abc.tar.gz --thread=6 --level=9
```
### Performance:
221GB: zip to 100GB in 44 minutes

 - `--threads=6`: from 1 to maximum(your-total-cpu-core)
 - `--level=0|1|6|9`:
   - 0 - No Compression
   - 1 - Best Speed
   - 6 - Default Compression
   - 9 - Best Compression

## unzip a file
```
./gopgzip unzip --input=/home/harryzhu/test/abc.tar.gz --output=/the/path/of/the/abc.tar
```

## zstd a file
```
./gopgzip zstd --input=/home/harryzhu/test/abc.tar --output=/the/path/of/the/abc.zst --level=0|1|6|9
```

## unzstd a file
```
./gopgzip unzstd --input=/home/harryzhu/test/abc.zst --output=/the/path/of/the/abc.tar
```

## tar a folder recursively
```
./gopgzip tar --input=/the/path/of/the/folder --output=/home/harryzhu/test/abc.tar
```

## untar a file
```
./gopgzip untar --input=/home/harryzhu/test/abc.tar --output=/the/path/of/the/folder
```

## encrypt a file
```
./gopgzip encrypt --input=/home/harryzhu/test/abc.tar --output=/the/path/of/the/abc.tar.enc
```
 in your `/etc/profile`, add `export HARRYZHUENCRYPTKEY=Your-Password` to set your own PASSWORD
 or use `--password="Your-Password" --force` inline (NOT recommend)

## decrypt a file
```
./gopgzip decrypt --input=/home/harryzhu/test/abc.tar.enc --output=/the/path/of/the/abc.tar
```

## md5 sum a file: for a small file
```
./gopgzip md5 --input=/home/harryzhu/test/abc.tar.gz [--output=/the/path/of/the/abc.tar.md5]
```

## sha256 sum a file: for a small file
```
./gopgzip sha256 --input=/home/harryzhu/test/abc.tar.gz [--output=/the/path/of/the/abc.tar.md5]
```

## xxhash sum a file: for a big file
```
./gopgzip xxhash --input=/home/harryzhu/test/abc.tar.gz [--output=/the/path/of/the/abc.tar.xxhash]
```

## b3sum a file: for a very big file
```
./gopgzip b3sum --input=/home/harryzhu/test/abc.tar.gz [--output=/the/path/of/the/abc.tar.b3sum]
```
